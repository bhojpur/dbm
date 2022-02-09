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
	"errors"
	"fmt"
	"reflect"

	"github.com/bhojpur/dbm/pkg/orm/core"
	"github.com/bhojpur/sql/pkg/builder"
)

// Rows rows wrapper a rows to
type Rows struct {
	session  *Session
	rows     *core.Rows
	beanType reflect.Type
}

func newRows(session *Session, bean interface{}) (*Rows, error) {
	rows := new(Rows)
	rows.session = session
	rows.beanType = reflect.Indirect(reflect.ValueOf(bean)).Type()
	var sqlStr string
	var args []interface{}
	var err error
	beanValue := reflect.ValueOf(bean)
	if beanValue.Kind() != reflect.Ptr {
		return nil, errors.New("needs a pointer to a value")
	} else if beanValue.Elem().Kind() == reflect.Ptr {
		return nil, errors.New("a pointer to a pointer is not allowed")
	}
	if err = rows.session.statement.SetRefBean(bean); err != nil {
		return nil, err
	}
	if len(session.statement.TableName()) == 0 {
		return nil, ErrTableNotFound
	}
	if rows.session.statement.RawSQL == "" {
		var autoCond builder.Cond
		var addedTableName = (len(session.statement.JoinStr) > 0)
		var table = rows.session.statement.RefTable
		if !session.statement.NoAutoCondition {
			var err error
			autoCond, err = session.statement.BuildConds(table, bean, true, true, false, true, addedTableName)
			if err != nil {
				return nil, err
			}
		} else {
			// !oinume! Add "<col> IS NULL" to WHERE whatever condiBean is given.
			if col := table.DeletedColumn(); col != nil && !session.statement.GetUnscoped() { // tag "deleted" is enabled
				autoCond = session.statement.CondDeleted(col)
			}
		}
		sqlStr, args, err = rows.session.statement.GenFindSQL(autoCond)
		if err != nil {
			return nil, err
		}
	} else {
		sqlStr = rows.session.statement.GenRawSQL()
		args = rows.session.statement.RawParams
	}
	rows.rows, err = rows.session.queryRows(sqlStr, args...)
	if err != nil {
		rows.Close()
		return nil, err
	}
	return rows, nil
}

// Next move cursor to next record, return false if end has reached
func (rows *Rows) Next() bool {
	if rows.rows != nil {
		return rows.rows.Next()
	}
	return false
}

// Err returns the error, if any, that was encountered during iteration. Err may be called after an explicit or implicit Close.
func (rows *Rows) Err() error {
	if rows.rows != nil {
		return rows.rows.Err()
	}
	return nil
}

// Scan row record to bean properties
func (rows *Rows) Scan(beans ...interface{}) error {
	if rows.Err() != nil {
		return rows.Err()
	}
	var bean = beans[0]
	var tp = reflect.TypeOf(bean)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	var beanKind = tp.Kind()
	if len(beans) == 1 {
		if reflect.Indirect(reflect.ValueOf(bean)).Type() != rows.beanType {
			return fmt.Errorf("scan arg is incompatible type to [%v]", rows.beanType)
		}
		if err := rows.session.statement.SetRefBean(bean); err != nil {
			return err
		}
	}
	fields, err := rows.rows.Columns()
	if err != nil {
		return err
	}
	types, err := rows.rows.ColumnTypes()
	if err != nil {
		return err
	}
	if err := rows.session.scan(rows.rows, rows.session.statement.RefTable, beanKind, beans, types, fields); err != nil {
		return err
	}
	return rows.session.executeProcessors()
}

// Close session if session.IsAutoClose is true, and claimed any opened resources
func (rows *Rows) Close() error {
	if rows.session.isAutoClose {
		defer rows.session.Close()
	}
	if rows.rows != nil {
		return rows.rows.Close()
	}
	return nil
}
