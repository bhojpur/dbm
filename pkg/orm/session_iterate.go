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

	"github.com/bhojpur/dbm/pkg/orm/internal/utils"
)

// IterFunc only use by Iterate
type IterFunc func(idx int, bean interface{}) error

// Rows return sql.Rows compatible Rows obj, as a forward Iterator object for iterating record by record, bean's non-empty fields
// are conditions.
func (session *Session) Rows(bean interface{}) (*Rows, error) {
	return newRows(session, bean)
}

// Iterate record by record handle records from table, condiBeans's non-empty fields
// are conditions. beans could be []Struct, []*Struct, map[int64]Struct
// map[int64]*Struct
func (session *Session) Iterate(bean interface{}, fun IterFunc) error {
	if session.isAutoClose {
		defer session.Close()
	}
	if session.statement.LastError != nil {
		return session.statement.LastError
	}
	if session.statement.BufferSize > 0 {
		return session.bufferIterate(bean, fun)
	}
	rows, err := session.Rows(bean)
	if err != nil {
		return err
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		b := reflect.New(rows.beanType).Interface()
		err = rows.Scan(b)
		if err != nil {
			return err
		}
		err = fun(i, b)
		if err != nil {
			return err
		}
		i++
	}
	return rows.Err()
}

// BufferSize sets the buffersize for iterate
func (session *Session) BufferSize(size int) *Session {
	session.statement.BufferSize = size
	return session
}
func (session *Session) bufferIterate(bean interface{}, fun IterFunc) error {
	var bufferSize = session.statement.BufferSize
	var pLimitN = session.statement.LimitN
	if pLimitN != nil && bufferSize > *pLimitN {
		bufferSize = *pLimitN
	}
	var start = session.statement.Start
	v := utils.ReflectValue(bean)
	sliceType := reflect.SliceOf(v.Type())
	var idx = 0
	session.autoResetStatement = false
	defer func() {
		session.autoResetStatement = true
	}()
	for bufferSize > 0 {
		slice := reflect.New(sliceType)
		if err := session.NoCache().Limit(bufferSize, start).find(slice.Interface(), bean); err != nil {
			return err
		}
		for i := 0; i < slice.Elem().Len(); i++ {
			if err := fun(idx, slice.Elem().Index(i).Addr().Interface()); err != nil {
				return err
			}
			idx++
		}
		if bufferSize > slice.Elem().Len() {
			break
		}
		start += slice.Elem().Len()
		if pLimitN != nil && start+bufferSize > *pLimitN {
			bufferSize = *pLimitN - start
		}
	}
	return nil
}
