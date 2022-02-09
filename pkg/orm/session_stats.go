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
	"database/sql"
	"errors"
	"reflect"
)

// Count counts the records. bean's non-empty fields
// are conditions.
func (session *Session) Count(bean ...interface{}) (int64, error) {
	if session.isAutoClose {
		defer session.Close()
	}
	sqlStr, args, err := session.statement.GenCountSQL(bean...)
	if err != nil {
		return 0, err
	}
	var total int64
	err = session.queryRow(sqlStr, args...).Scan(&total)
	if err == sql.ErrNoRows || err == nil {
		return total, nil
	}
	return 0, err
}

// sum call sum some column. bean's non-empty fields are conditions.
func (session *Session) sum(res interface{}, bean interface{}, columnNames ...string) error {
	if session.isAutoClose {
		defer session.Close()
	}
	v := reflect.ValueOf(res)
	if v.Kind() != reflect.Ptr {
		return errors.New("need a pointer to a variable")
	}
	sqlStr, args, err := session.statement.GenSumSQL(bean, columnNames...)
	if err != nil {
		return err
	}
	if v.Elem().Kind() == reflect.Slice {
		err = session.queryRow(sqlStr, args...).ScanSlice(res)
	} else {
		err = session.queryRow(sqlStr, args...).Scan(res)
	}
	if err == sql.ErrNoRows || err == nil {
		return nil
	}
	return err
}

// Sum call sum some column. bean's non-empty fields are conditions.
func (session *Session) Sum(bean interface{}, columnName string) (res float64, err error) {
	return res, session.sum(&res, bean, columnName)
}

// SumInt call sum some column. bean's non-empty fields are conditions.
func (session *Session) SumInt(bean interface{}, columnName string) (res int64, err error) {
	return res, session.sum(&res, bean, columnName)
}

// Sums call sum some columns. bean's non-empty fields are conditions.
func (session *Session) Sums(bean interface{}, columnNames ...string) ([]float64, error) {
	var res = make([]float64, len(columnNames))
	return res, session.sum(&res, bean, columnNames...)
}

// SumsInt sum specify columns and return as []int64 instead of []float64
func (session *Session) SumsInt(bean interface{}, columnNames ...string) ([]int64, error) {
	var res = make([]int64, len(columnNames))
	return res, session.sum(&res, bean, columnNames...)
}
