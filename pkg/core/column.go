package core

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
	"strings"
	"time"
)

const (
	TWOSIDES = iota + 1
	ONLYTODB
	ONLYFROMDB
)

// Column defines database column
type Column struct {
	Name            string
	TableName       string
	FieldName       string
	SQLType         SQLType
	IsJSON          bool
	Length          int
	Length2         int
	Nullable        bool
	Default         string
	Indexes         map[string]int
	IsPrimaryKey    bool
	IsAutoIncrement bool
	MapType         int
	IsCreated       bool
	IsUpdated       bool
	IsDeleted       bool
	IsCascade       bool
	IsVersion       bool
	DefaultIsEmpty  bool // false means column has no default set, but not default value is empty
	EnumOptions     map[string]int
	SetOptions      map[string]int
	DisableTimeZone bool
	TimeZone        *time.Location // column specified time zone
	Comment         string
}

// NewColumn creates a new column
func NewColumn(name, fieldName string, sqlType SQLType, len1, len2 int, nullable bool) *Column {
	return &Column{
		Name:            name,
		TableName:       "",
		FieldName:       fieldName,
		SQLType:         sqlType,
		Length:          len1,
		Length2:         len2,
		Nullable:        nullable,
		Default:         "",
		Indexes:         make(map[string]int),
		IsPrimaryKey:    false,
		IsAutoIncrement: false,
		MapType:         TWOSIDES,
		IsCreated:       false,
		IsUpdated:       false,
		IsDeleted:       false,
		IsCascade:       false,
		IsVersion:       false,
		DefaultIsEmpty:  true, // default should be no default
		EnumOptions:     make(map[string]int),
		Comment:         "",
	}
}

// String generate column description string according dialect
func (col *Column) String(d Dialect) string {
	sql := d.Quote(col.Name) + " "
	sql += d.SqlType(col) + " "
	if col.IsPrimaryKey {
		sql += "PRIMARY KEY "
		if col.IsAutoIncrement {
			sql += d.AutoIncrStr() + " "
		}
	}
	if col.Default != "" {
		sql += "DEFAULT " + col.Default + " "
	}
	if d.ShowCreateNull() {
		if col.Nullable {
			sql += "NULL "
		} else {
			sql += "NOT NULL "
		}
	}
	return sql
}

// StringNoPk generate column description string according dialect without primary keys
func (col *Column) StringNoPk(d Dialect) string {
	sql := d.Quote(col.Name) + " "
	sql += d.SqlType(col) + " "
	if col.Default != "" {
		sql += "DEFAULT " + col.Default + " "
	}
	if d.ShowCreateNull() {
		if col.Nullable {
			sql += "NULL "
		} else {
			sql += "NOT NULL "
		}
	}
	return sql
}

// ValueOf returns column's filed of struct's value
func (col *Column) ValueOf(bean interface{}) (*reflect.Value, error) {
	dataStruct := reflect.Indirect(reflect.ValueOf(bean))
	return col.ValueOfV(&dataStruct)
}

// ValueOfV returns column's filed of struct's value accept reflevt value
func (col *Column) ValueOfV(dataStruct *reflect.Value) (*reflect.Value, error) {
	var fieldValue reflect.Value
	fieldPath := strings.Split(col.FieldName, ".")
	if dataStruct.Type().Kind() == reflect.Map {
		keyValue := reflect.ValueOf(fieldPath[len(fieldPath)-1])
		fieldValue = dataStruct.MapIndex(keyValue)
		return &fieldValue, nil
	} else if dataStruct.Type().Kind() == reflect.Interface {
		structValue := reflect.ValueOf(dataStruct.Interface())
		dataStruct = &structValue
	}
	level := len(fieldPath)
	fieldValue = dataStruct.FieldByName(fieldPath[0])
	for i := 0; i < level-1; i++ {
		if !fieldValue.IsValid() {
			break
		}
		if fieldValue.Kind() == reflect.Struct {
			fieldValue = fieldValue.FieldByName(fieldPath[i+1])
		} else if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			}
			fieldValue = fieldValue.Elem().FieldByName(fieldPath[i+1])
		} else {
			return nil, fmt.Errorf("field %v is not valid", col.FieldName)
		}
	}
	if !fieldValue.IsValid() {
		return nil, fmt.Errorf("field %v is not valid", col.FieldName)
	}
	return &fieldValue, nil
}
