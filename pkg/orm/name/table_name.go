package name

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
	"reflect"
	"sync"
)

// TableName table name interface to define customerize table name
type TableName interface {
	TableName() string
}
type TableComment interface {
	TableComment() string
}

var (
	tpTableName    = reflect.TypeOf((*TableName)(nil)).Elem()
	tpTableComment = reflect.TypeOf((*TableComment)(nil)).Elem()
	tvCache        sync.Map
	tcCache        sync.Map
)

// GetTableName returns table name
func GetTableName(mapper Mapper, v reflect.Value) string {
	if v.Type().Implements(tpTableName) {
		return v.Interface().(TableName).TableName()
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		if v.Type().Implements(tpTableName) {
			return v.Interface().(TableName).TableName()
		}
	} else if v.CanAddr() {
		v1 := v.Addr()
		if v1.Type().Implements(tpTableName) {
			return v1.Interface().(TableName).TableName()
		}
	} else {
		name, ok := tvCache.Load(v.Type())
		if ok {
			if name.(string) != "" {
				return name.(string)
			}
		} else {
			v2 := reflect.New(v.Type())
			if v2.Type().Implements(tpTableName) {
				tableName := v2.Interface().(TableName).TableName()
				tvCache.Store(v.Type(), tableName)
				return tableName
			}
			tvCache.Store(v.Type(), "")
		}
	}
	return mapper.Obj2Table(v.Type().Name())
}

// GetTableComment returns table comment
func GetTableComment(v reflect.Value) string {
	if v.Type().Implements(tpTableComment) {
		return v.Interface().(TableComment).TableComment()
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		if v.Type().Implements(tpTableComment) {
			return v.Interface().(TableComment).TableComment()
		}
	} else if v.CanAddr() {
		v1 := v.Addr()
		if v1.Type().Implements(tpTableComment) {
			return v1.Interface().(TableComment).TableComment()
		}
	} else {
		comment, ok := tcCache.Load(v.Type())
		if ok {
			if comment.(string) != "" {
				return comment.(string)
			}
		} else {
			v2 := reflect.New(v.Type())
			if v2.Type().Implements(tpTableComment) {
				tableComment := v2.Interface().(TableComment).TableComment()
				tcCache.Store(v.Type(), tableComment)
				return tableComment
			}
			tcCache.Store(v.Type(), "")
		}
	}
	return ""
}
