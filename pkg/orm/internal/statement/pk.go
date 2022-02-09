package statement

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
	"fmt"
	"reflect"

	"github.com/bhojpur/dbm/pkg/orm/schema"
	"github.com/bhojpur/sql/pkg/builder"
)

var (
	ptrPkType  = reflect.TypeOf(&schema.PK{})
	pkType     = reflect.TypeOf(schema.PK{})
	stringType = reflect.TypeOf("")
	intType    = reflect.TypeOf(int64(0))
	uintType   = reflect.TypeOf(uint64(0))
)

// ErrIDConditionWithNoTable represents an error there is no reference table with an ID condition
type ErrIDConditionWithNoTable struct {
	ID schema.PK
}

func (err ErrIDConditionWithNoTable) Error() string {
	return fmt.Sprintf("ID condition %#v need reference table", err.ID)
}

// IsIDConditionWithNoTableErr return true if the err is ErrIDConditionWithNoTable
func IsIDConditionWithNoTableErr(err error) bool {
	_, ok := err.(ErrIDConditionWithNoTable)
	return ok
}

// ID generate "where id = ? " statement or for composite key "where key1 = ? and key2 = ?"
func (statement *Statement) ID(id interface{}) *Statement {
	switch t := id.(type) {
	case *schema.PK:
		statement.idParam = *t
	case schema.PK:
		statement.idParam = t
	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		statement.idParam = schema.PK{id}
	default:
		idValue := reflect.ValueOf(id)
		idType := idValue.Type()
		switch idType.Kind() {
		case reflect.String:
			statement.idParam = schema.PK{idValue.Convert(stringType).Interface()}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			statement.idParam = schema.PK{idValue.Convert(intType).Interface()}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			statement.idParam = schema.PK{idValue.Convert(uintType).Interface()}
		case reflect.Slice:
			if idType.ConvertibleTo(pkType) {
				statement.idParam = idValue.Convert(pkType).Interface().(schema.PK)
			}
		case reflect.Ptr:
			if idType.ConvertibleTo(ptrPkType) {
				statement.idParam = idValue.Convert(ptrPkType).Elem().Interface().(schema.PK)
			}
		}
	}
	if statement.idParam == nil {
		statement.LastError = fmt.Errorf("ID param %#v is not supported", id)
	}
	return statement
}

// ProcessIDParam handles the process of id condition
func (statement *Statement) ProcessIDParam() error {
	if statement.idParam == nil {
		return nil
	}
	if statement.RefTable == nil {
		return ErrIDConditionWithNoTable{statement.idParam}
	}
	if len(statement.RefTable.PrimaryKeys) != len(statement.idParam) {
		return fmt.Errorf("ID condition is error, expect %d primarykeys, there are %d",
			len(statement.RefTable.PrimaryKeys),
			len(statement.idParam),
		)
	}
	for i, col := range statement.RefTable.PKColumns() {
		var colName = statement.colName(col, statement.TableName())
		statement.cond = statement.cond.And(builder.Eq{colName: statement.idParam[i]})
	}
	return nil
}
