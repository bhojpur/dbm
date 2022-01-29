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
	"context"
	"database/sql"
	"errors"
	"reflect"
)

// Stmt reprents a stmt objects
type Stmt struct {
	*sql.Stmt
	db    *DB
	names map[string]int
}

func (db *DB) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	names := make(map[string]int)
	var i int
	query = re.ReplaceAllStringFunc(query, func(src string) string {
		names[src[1:]] = i
		i += 1
		return "?"
	})
	stmt, err := db.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Stmt{stmt, db, names}, nil
}
func (db *DB) Prepare(query string) (*Stmt, error) {
	return db.PrepareContext(context.Background(), query)
}
func (s *Stmt) ExecMapContext(ctx context.Context, mp interface{}) (sql.Result, error) {
	vv := reflect.ValueOf(mp)
	if vv.Kind() != reflect.Ptr || vv.Elem().Kind() != reflect.Map {
		return nil, errors.New("mp should be a map's pointer")
	}
	args := make([]interface{}, len(s.names))
	for k, i := range s.names {
		args[i] = vv.Elem().MapIndex(reflect.ValueOf(k)).Interface()
	}
	return s.Stmt.ExecContext(ctx, args...)
}
func (s *Stmt) ExecMap(mp interface{}) (sql.Result, error) {
	return s.ExecMapContext(context.Background(), mp)
}
func (s *Stmt) ExecStructContext(ctx context.Context, st interface{}) (sql.Result, error) {
	vv := reflect.ValueOf(st)
	if vv.Kind() != reflect.Ptr || vv.Elem().Kind() != reflect.Struct {
		return nil, errors.New("mp should be a map's pointer")
	}
	args := make([]interface{}, len(s.names))
	for k, i := range s.names {
		args[i] = vv.Elem().FieldByName(k).Interface()
	}
	return s.Stmt.ExecContext(ctx, args...)
}
func (s *Stmt) ExecStruct(st interface{}) (sql.Result, error) {
	return s.ExecStructContext(context.Background(), st)
}
func (s *Stmt) QueryContext(ctx context.Context, args ...interface{}) (*Rows, error) {
	rows, err := s.Stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{rows, s.db}, nil
}
func (s *Stmt) Query(args ...interface{}) (*Rows, error) {
	return s.QueryContext(context.Background(), args...)
}
func (s *Stmt) QueryMapContext(ctx context.Context, mp interface{}) (*Rows, error) {
	vv := reflect.ValueOf(mp)
	if vv.Kind() != reflect.Ptr || vv.Elem().Kind() != reflect.Map {
		return nil, errors.New("mp should be a map's pointer")
	}
	args := make([]interface{}, len(s.names))
	for k, i := range s.names {
		args[i] = vv.Elem().MapIndex(reflect.ValueOf(k)).Interface()
	}
	return s.QueryContext(ctx, args...)
}
func (s *Stmt) QueryMap(mp interface{}) (*Rows, error) {
	return s.QueryMapContext(context.Background(), mp)
}
func (s *Stmt) QueryStructContext(ctx context.Context, st interface{}) (*Rows, error) {
	vv := reflect.ValueOf(st)
	if vv.Kind() != reflect.Ptr || vv.Elem().Kind() != reflect.Struct {
		return nil, errors.New("mp should be a map's pointer")
	}
	args := make([]interface{}, len(s.names))
	for k, i := range s.names {
		args[i] = vv.Elem().FieldByName(k).Interface()
	}
	return s.Query(args...)
}
func (s *Stmt) QueryStruct(st interface{}) (*Rows, error) {
	return s.QueryStructContext(context.Background(), st)
}
func (s *Stmt) QueryRowContext(ctx context.Context, args ...interface{}) *Row {
	rows, err := s.QueryContext(ctx, args...)
	return &Row{rows, err}
}
func (s *Stmt) QueryRow(args ...interface{}) *Row {
	return s.QueryRowContext(context.Background(), args...)
}
func (s *Stmt) QueryRowMapContext(ctx context.Context, mp interface{}) *Row {
	vv := reflect.ValueOf(mp)
	if vv.Kind() != reflect.Ptr || vv.Elem().Kind() != reflect.Map {
		return &Row{nil, errors.New("mp should be a map's pointer")}
	}
	args := make([]interface{}, len(s.names))
	for k, i := range s.names {
		args[i] = vv.Elem().MapIndex(reflect.ValueOf(k)).Interface()
	}
	return s.QueryRowContext(ctx, args...)
}
func (s *Stmt) QueryRowMap(mp interface{}) *Row {
	return s.QueryRowMapContext(context.Background(), mp)
}
func (s *Stmt) QueryRowStructContext(ctx context.Context, st interface{}) *Row {
	vv := reflect.ValueOf(st)
	if vv.Kind() != reflect.Ptr || vv.Elem().Kind() != reflect.Struct {
		return &Row{nil, errors.New("st should be a struct's pointer")}
	}
	args := make([]interface{}, len(s.names))
	for k, i := range s.names {
		args[i] = vv.Elem().FieldByName(k).Interface()
	}
	return s.QueryRowContext(ctx, args...)
}
func (s *Stmt) QueryRowStruct(st interface{}) *Row {
	return s.QueryRowStructContext(context.Background(), st)
}
