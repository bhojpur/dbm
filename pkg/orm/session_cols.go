package orm

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
	"strings"
	"time"

	schemasvr "github.com/bhojpur/dbm/pkg/orm/schema"
)

func setColumnInt(bean interface{}, col *schemasvr.Column, t int64) {
	v, err := col.ValueOf(bean)
	if err != nil {
		return
	}
	if v.CanSet() {
		switch v.Type().Kind() {
		case reflect.Int, reflect.Int64, reflect.Int32:
			v.SetInt(t)
		case reflect.Uint, reflect.Uint64, reflect.Uint32:
			v.SetUint(uint64(t))
		}
	}
}
func setColumnTime(bean interface{}, col *schemasvr.Column, t time.Time) {
	v, err := col.ValueOf(bean)
	if err != nil {
		return
	}
	if v.CanSet() {
		switch v.Type().Kind() {
		case reflect.Struct:
			v.Set(reflect.ValueOf(t).Convert(v.Type()))
		case reflect.Int, reflect.Int64, reflect.Int32:
			v.SetInt(t.Unix())
		case reflect.Uint, reflect.Uint64, reflect.Uint32:
			v.SetUint(uint64(t.Unix()))
		}
	}
}
func getFlagForColumn(m map[string]bool, col *schemasvr.Column) (val bool, has bool) {
	if len(m) == 0 {
		return false, false
	}
	n := len(col.Name)
	for mk := range m {
		if len(mk) != n {
			continue
		}
		if strings.EqualFold(mk, col.Name) {
			return m[mk], true
		}
	}
	return false, false
}

// Incr provides a query string like "count = count + 1"
func (session *Session) Incr(column string, arg ...interface{}) *Session {
	session.statement.Incr(column, arg...)
	return session
}

// Decr provides a query string like "count = count - 1"
func (session *Session) Decr(column string, arg ...interface{}) *Session {
	session.statement.Decr(column, arg...)
	return session
}

// SetExpr provides a query string like "column = {expression}"
func (session *Session) SetExpr(column string, expression interface{}) *Session {
	session.statement.SetExpr(column, expression)
	return session
}

// Select provides some columns to special
func (session *Session) Select(str string) *Session {
	session.statement.Select(str)
	return session
}

// Cols provides some columns to special
func (session *Session) Cols(columns ...string) *Session {
	session.statement.Cols(columns...)
	return session
}

// AllCols ask all columns
func (session *Session) AllCols() *Session {
	session.statement.AllCols()
	return session
}

// MustCols specify some columns must use even if they are empty
func (session *Session) MustCols(columns ...string) *Session {
	session.statement.MustCols(columns...)
	return session
}

// UseBool automatically retrieve condition according struct, but
// if struct has bool field, it will ignore them. So use UseBool
// to tell system to do not ignore them.
// If no parameters, it will use all the bool field of struct, or
// it will use parameters's columns
func (session *Session) UseBool(columns ...string) *Session {
	session.statement.UseBool(columns...)
	return session
}

// Distinct use for distinct columns. Caution: when you are using cache,
// distinct will not be cached because cache system need id,
// but distinct will not provide id
func (session *Session) Distinct(columns ...string) *Session {
	session.statement.Distinct(columns...)
	return session
}

// Omit Only not use the parameters as select or update columns
func (session *Session) Omit(columns ...string) *Session {
	session.statement.Omit(columns...)
	return session
}

// Nullable Set null when column is zero-value and nullable for update
func (session *Session) Nullable(columns ...string) *Session {
	session.statement.Nullable(columns...)
	return session
}

// NoAutoTime means do not automatically give created field and updated field
// the current time on the current session temporarily
func (session *Session) NoAutoTime() *Session {
	session.statement.UseAutoTime = false
	return session
}
