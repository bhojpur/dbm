package tags

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/bhojpur/dbm/pkg/orm/cache"
	"github.com/bhojpur/dbm/pkg/orm/convert"
	"github.com/bhojpur/dbm/pkg/orm/dialect"
	"github.com/bhojpur/dbm/pkg/orm/name"
	"github.com/bhojpur/dbm/pkg/orm/schema"
)

var (
	// ErrUnsupportedType represents an unsupported type error
	ErrUnsupportedType = errors.New("unsupported type")
)

// Parser represents a parser for ORM tag
type Parser struct {
	identifier   string
	dialect      dialect.Dialect
	columnMapper name.Mapper
	tableMapper  name.Mapper
	handlers     map[string]Handler
	cacherMgr    *cache.Manager
	tableCache   sync.Map // map[reflect.Type]*schemas.Table
}

// NewParser creates a tag parser
func NewParser(identifier string, dialect dialect.Dialect, tableMapper, columnMapper name.Mapper, cacherMgr *cache.Manager) *Parser {
	return &Parser{
		identifier:   identifier,
		dialect:      dialect,
		tableMapper:  tableMapper,
		columnMapper: columnMapper,
		handlers:     defaultTagHandlers,
		cacherMgr:    cacherMgr,
	}
}

// GetTableMapper returns table mapper
func (parser *Parser) GetTableMapper() name.Mapper {
	return parser.tableMapper
}

// SetTableMapper sets table mapper
func (parser *Parser) SetTableMapper(mapper name.Mapper) {
	parser.ClearCaches()
	parser.tableMapper = mapper
}

// GetColumnMapper returns column mapper
func (parser *Parser) GetColumnMapper() name.Mapper {
	return parser.columnMapper
}

// SetColumnMapper sets column mapper
func (parser *Parser) SetColumnMapper(mapper name.Mapper) {
	parser.ClearCaches()
	parser.columnMapper = mapper
}

// SetIdentifier sets tag identifier
func (parser *Parser) SetIdentifier(identifier string) {
	parser.ClearCaches()
	parser.identifier = identifier
}

// ParseWithCache parse a struct with cache
func (parser *Parser) ParseWithCache(v reflect.Value) (*schema.Table, error) {
	t := v.Type()
	tableI, ok := parser.tableCache.Load(t)
	if ok {
		return tableI.(*schema.Table), nil
	}
	table, err := parser.Parse(v)
	if err != nil {
		return nil, err
	}
	parser.tableCache.Store(t, table)
	if parser.cacherMgr.GetDefaultCacher() != nil {
		if v.CanAddr() {
			gob.Register(v.Addr().Interface())
		} else {
			gob.Register(v.Interface())
		}
	}
	return table, nil
}

// ClearCacheTable removes the database mapper of a type from the cache
func (parser *Parser) ClearCacheTable(t reflect.Type) {
	parser.tableCache.Delete(t)
}

// ClearCaches removes all the cached table information parsed by structs
func (parser *Parser) ClearCaches() {
	parser.tableCache = sync.Map{}
}
func addIndex(indexName string, table *schema.Table, col *schema.Column, indexType int) {
	if index, ok := table.Indexes[indexName]; ok {
		index.AddColumn(col.Name)
		col.Indexes[index.Name] = indexType
	} else {
		index := schema.NewIndex(indexName, indexType)
		index.AddColumn(col.Name)
		table.AddIndex(index)
		col.Indexes[index.Name] = indexType
	}
}

// ErrIgnoreField represents an error to ignore field
var ErrIgnoreField = errors.New("field will be ignored")

func (parser *Parser) getSQLTypeByType(t reflect.Type) (schema.SQLType, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Struct {
		v, ok := parser.tableCache.Load(t)
		if ok {
			pkCols := v.(*schema.Table).PKColumns()
			if len(pkCols) == 1 {
				return pkCols[0].SQLType, nil
			}
			if len(pkCols) > 1 {
				return schema.SQLType{}, fmt.Errorf("unsupported mulitiple primary key on cascade")
			}
		}
	}
	return schema.Type2SQLType(t), nil
}
func (parser *Parser) parseFieldWithNoTag(fieldIndex int, field reflect.StructField, fieldValue reflect.Value) (*schema.Column, error) {
	var sqlType schema.SQLType
	if fieldValue.CanAddr() {
		if _, ok := fieldValue.Addr().Interface().(convert.Conversion); ok {
			sqlType = schema.SQLType{Name: schema.Text}
		}
	}
	if _, ok := fieldValue.Interface().(convert.Conversion); ok {
		sqlType = schema.SQLType{Name: schema.Text}
	} else {
		var err error
		sqlType, err = parser.getSQLTypeByType(field.Type)
		if err != nil {
			return nil, err
		}
	}
	col := schema.NewColumn(parser.columnMapper.Obj2Table(field.Name),
		field.Name, sqlType, sqlType.DefaultLength,
		sqlType.DefaultLength2, true)
	col.FieldIndex = []int{fieldIndex}
	if field.Type.Kind() == reflect.Int64 && (strings.ToUpper(col.FieldName) == "ID" || strings.HasSuffix(strings.ToUpper(col.FieldName), ".ID")) {
		col.IsAutoIncrement = true
		col.IsPrimaryKey = true
		col.Nullable = false
	}
	return col, nil
}
func (parser *Parser) parseFieldWithTags(table *schema.Table, fieldIndex int, field reflect.StructField, fieldValue reflect.Value, tags []tag) (*schema.Column, error) {
	var col = &schema.Column{
		FieldName:       field.Name,
		FieldIndex:      []int{fieldIndex},
		Nullable:        true,
		IsPrimaryKey:    false,
		IsAutoIncrement: false,
		MapType:         schema.TWOSIDES,
		Indexes:         make(map[string]int),
		DefaultIsEmpty:  true,
	}
	var ctx = Context{
		table:      table,
		col:        col,
		fieldValue: fieldValue,
		indexNames: make(map[string]int),
		parser:     parser,
	}
	for j, tag := range tags {
		if ctx.ignoreNext {
			ctx.ignoreNext = false
			continue
		}
		ctx.tag = tag
		ctx.tagUname = strings.ToUpper(tag.name)
		if j > 0 {
			ctx.preTag = strings.ToUpper(tags[j-1].name)
		}
		if j < len(tags)-1 {
			ctx.nextTag = tags[j+1].name
		} else {
			ctx.nextTag = ""
		}
		if h, ok := parser.handlers[ctx.tagUname]; ok {
			if err := h(&ctx); err != nil {
				return nil, err
			}
		} else {
			if strings.HasPrefix(ctx.tag.name, "'") && strings.HasSuffix(ctx.tag.name, "'") {
				col.Name = ctx.tag.name[1 : len(ctx.tag.name)-1]
			} else {
				col.Name = ctx.tag.name
			}
		}
		if ctx.hasCacheTag {
			if parser.cacherMgr.GetDefaultCacher() != nil {
				parser.cacherMgr.SetCacher(table.Name, parser.cacherMgr.GetDefaultCacher())
			} else {
				parser.cacherMgr.SetCacher(table.Name, cache.NewLRUCacher2(cache.NewMemoryStore(), time.Hour, 10000))
			}
		}
		if ctx.hasNoCacheTag {
			parser.cacherMgr.SetCacher(table.Name, nil)
		}
	}
	if col.SQLType.Name == "" {
		var err error
		col.SQLType, err = parser.getSQLTypeByType(field.Type)
		if err != nil {
			return nil, err
		}
	}
	if ctx.isUnsigned && col.SQLType.IsNumeric() && !strings.HasPrefix(col.SQLType.Name, "UNSIGNED") {
		col.SQLType.Name = "UNSIGNED " + col.SQLType.Name
	}
	parser.dialect.SQLType(col)
	if col.Length == 0 {
		col.Length = col.SQLType.DefaultLength
	}
	if col.Length2 == 0 {
		col.Length2 = col.SQLType.DefaultLength2
	}
	if col.Name == "" {
		col.Name = parser.columnMapper.Obj2Table(field.Name)
	}
	if ctx.isUnique {
		ctx.indexNames[col.Name] = schema.UniqueType
	} else if ctx.isIndex {
		ctx.indexNames[col.Name] = schema.IndexType
	}
	for indexName, indexType := range ctx.indexNames {
		addIndex(indexName, table, col, indexType)
	}
	return col, nil
}
func (parser *Parser) parseField(table *schema.Table, fieldIndex int, field reflect.StructField, fieldValue reflect.Value) (*schema.Column, error) {
	if isNotTitle(field.Name) {
		return nil, ErrIgnoreField
	}
	var (
		tag       = field.Tag
		ormTagStr = strings.TrimSpace(tag.Get(parser.identifier))
	)
	if ormTagStr == "-" {
		return nil, ErrIgnoreField
	}
	if ormTagStr == "" {
		return parser.parseFieldWithNoTag(fieldIndex, field, fieldValue)
	}
	tags, err := splitTag(ormTagStr)
	if err != nil {
		return nil, err
	}
	return parser.parseFieldWithTags(table, fieldIndex, field, fieldValue, tags)
}
func isNotTitle(n string) bool {
	for _, c := range n {
		return unicode.IsLower(c)
	}
	return true
}

// Parse parses a struct as a table information
func (parser *Parser) Parse(v reflect.Value) (*schema.Table, error) {
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, ErrUnsupportedType
	}
	table := schema.NewEmptyTable()
	table.Type = t
	table.Name = name.GetTableName(parser.tableMapper, v)
	table.Comment = name.GetTableComment(v)
	for i := 0; i < t.NumField(); i++ {
		col, err := parser.parseField(table, i, t.Field(i), v.Field(i))
		if err == ErrIgnoreField {
			continue
		} else if err != nil {
			return nil, err
		}
		table.AddColumn(col)
	} // end for
	return table, nil
}
