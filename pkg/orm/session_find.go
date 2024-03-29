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
	"reflect"

	"github.com/bhojpur/dbm/pkg/orm/cache"
	"github.com/bhojpur/dbm/pkg/orm/convert"
	statementsvr "github.com/bhojpur/dbm/pkg/orm/internal/statement"
	"github.com/bhojpur/dbm/pkg/orm/internal/utils"
	schemasvr "github.com/bhojpur/dbm/pkg/orm/schema"
	"github.com/bhojpur/sql/pkg/builder"
)

const (
	tpStruct = iota
	tpNonStruct
)

// Find retrieve records from table, condiBeans's non-empty fields
// are conditions. beans could be []Struct, []*Struct, map[int64]Struct
// map[int64]*Struct
func (session *Session) Find(rowsSlicePtr interface{}, condiBean ...interface{}) error {
	if session.isAutoClose {
		defer session.Close()
	}
	return session.find(rowsSlicePtr, condiBean...)
}

// FindAndCount find the results and also return the counts
func (session *Session) FindAndCount(rowsSlicePtr interface{}, condiBean ...interface{}) (int64, error) {
	if session.isAutoClose {
		defer session.Close()
	}
	session.autoResetStatement = false
	err := session.find(rowsSlicePtr, condiBean...)
	if err != nil {
		return 0, err
	}
	sliceValue := reflect.Indirect(reflect.ValueOf(rowsSlicePtr))
	if sliceValue.Kind() != reflect.Slice && sliceValue.Kind() != reflect.Map {
		return 0, errors.New("needs a pointer to a slice or a map")
	}
	sliceElementType := sliceValue.Type().Elem()
	if sliceElementType.Kind() == reflect.Ptr {
		sliceElementType = sliceElementType.Elem()
	}
	session.autoResetStatement = true
	if session.statement.SelectStr != "" {
		session.statement.SelectStr = ""
	}
	if len(session.statement.ColumnMap) > 0 && !session.statement.IsDistinct {
		session.statement.ColumnMap = []string{}
	}
	if session.statement.OrderStr != "" {
		session.statement.OrderStr = ""
	}
	if session.statement.LimitN != nil {
		session.statement.LimitN = nil
	}
	if session.statement.Start > 0 {
		session.statement.Start = 0
	}
	// session has stored the conditions so we use `unscoped` to avoid duplicated condition.
	if sliceElementType.Kind() == reflect.Struct {
		return session.Unscoped().Count(reflect.New(sliceElementType).Interface())
	}
	return session.Unscoped().Count()
}
func (session *Session) find(rowsSlicePtr interface{}, condiBean ...interface{}) error {
	defer session.resetStatement()
	if session.statement.LastError != nil {
		return session.statement.LastError
	}
	sliceValue := reflect.Indirect(reflect.ValueOf(rowsSlicePtr))
	var isSlice = sliceValue.Kind() == reflect.Slice
	var isMap = sliceValue.Kind() == reflect.Map
	if !isSlice && !isMap {
		return errors.New("needs a pointer to a slice or a map")
	}
	sliceElementType := sliceValue.Type().Elem()
	var tp = tpStruct
	if session.statement.RefTable == nil {
		if sliceElementType.Kind() == reflect.Ptr {
			if sliceElementType.Elem().Kind() == reflect.Struct {
				pv := reflect.New(sliceElementType.Elem())
				if err := session.statement.SetRefValue(pv); err != nil {
					return err
				}
			} else {
				tp = tpNonStruct
			}
		} else if sliceElementType.Kind() == reflect.Struct {
			pv := reflect.New(sliceElementType)
			if err := session.statement.SetRefValue(pv); err != nil {
				return err
			}
		} else {
			tp = tpNonStruct
		}
	}
	var (
		table          = session.statement.RefTable
		addedTableName = (len(session.statement.JoinStr) > 0)
		autoCond       builder.Cond
	)
	if tp == tpStruct {
		if !session.statement.NoAutoCondition && len(condiBean) > 0 {
			condTable, err := session.engine.tagParser.Parse(reflect.ValueOf(condiBean[0]))
			if err != nil {
				return err
			}
			autoCond, err = session.statement.BuildConds(condTable, condiBean[0], true, true, false, true, addedTableName)
			if err != nil {
				return err
			}
		} else {
			if col := table.DeletedColumn(); col != nil && !session.statement.GetUnscoped() { // tag "deleted" is enabled
				autoCond = session.statement.CondDeleted(col)
			}
		}
	}
	// if it's a map with Cols but primary key not in column list, we still need the primary key
	if isMap && !session.statement.ColumnMap.IsEmpty() {
		for _, k := range session.statement.RefTable.PrimaryKeys {
			session.statement.ColumnMap.Add(k)
		}
	}
	sqlStr, args, err := session.statement.GenFindSQL(autoCond)
	if err != nil {
		return err
	}
	if session.statement.ColumnMap.IsEmpty() && session.canCache() {
		if cacher := session.engine.GetCacher(session.statement.TableName()); cacher != nil &&
			!session.statement.IsDistinct &&
			!session.statement.GetUnscoped() {
			err = session.cacheFind(sliceElementType, sqlStr, rowsSlicePtr, args...)
			if err != ErrCacheFailed {
				return err
			}
			session.engine.logger.Warnf("Cache Find Failed")
		}
	}
	return session.noCacheFind(table, sliceValue, sqlStr, args...)
}
func (session *Session) noCacheFind(table *schemasvr.Table, containerValue reflect.Value, sqlStr string, args ...interface{}) error {
	elemType := containerValue.Type().Elem()
	var isPointer bool
	if elemType.Kind() == reflect.Ptr {
		isPointer = true
		elemType = elemType.Elem()
	}
	if elemType.Kind() == reflect.Ptr {
		return errors.New("pointer to pointer is not supported")
	}
	rows, err := session.queryRows(sqlStr, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	fields, err := rows.Columns()
	if err != nil {
		return err
	}
	types, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	var newElemFunc = func(fields []string) reflect.Value {
		return utils.New(elemType, len(fields), len(fields))
	}
	var containerValueSetFunc func(*reflect.Value, schemasvr.PK) error
	if containerValue.Kind() == reflect.Slice {
		containerValueSetFunc = func(newValue *reflect.Value, pk schemasvr.PK) error {
			if isPointer {
				containerValue.Set(reflect.Append(containerValue, newValue.Elem().Addr()))
			} else {
				containerValue.Set(reflect.Append(containerValue, newValue.Elem()))
			}
			return nil
		}
	} else {
		keyType := containerValue.Type().Key()
		if len(table.PrimaryKeys) == 0 {
			return errors.New("don't support multiple primary key's map has non-slice key type")
		}
		if len(table.PrimaryKeys) > 1 && keyType.Kind() != reflect.Slice {
			return errors.New("don't support multiple primary key's map has non-slice key type")
		}
		containerValueSetFunc = func(newValue *reflect.Value, pk schemasvr.PK) error {
			keyValue := reflect.New(keyType)
			cols := table.PKColumns()
			if len(cols) == 1 {
				if err := convert.AssignValue(keyValue, pk[0]); err != nil {
					return err
				}
			} else {
				keyValue.Set(reflect.ValueOf(&pk))
			}
			if isPointer {
				containerValue.SetMapIndex(keyValue.Elem(), newValue.Elem().Addr())
			} else {
				containerValue.SetMapIndex(keyValue.Elem(), newValue.Elem())
			}
			return nil
		}
	}
	if elemType.Kind() == reflect.Struct {
		var newValue = newElemFunc(fields)
		tb, err := session.engine.tagParser.ParseWithCache(newValue)
		if err != nil {
			return err
		}
		err = session.rows2Beans(rows, fields, types, tb, newElemFunc, containerValueSetFunc)
		rows.Close()
		if err != nil {
			return err
		}
		return session.executeProcessors()
	}
	for rows.Next() {
		var newValue = newElemFunc(fields)
		bean := newValue.Interface()
		switch elemType.Kind() {
		case reflect.Slice:
			err = session.getSlice(rows, types, fields, bean)
		case reflect.Map:
			err = session.getMap(rows, types, fields, bean)
		default:
			err = rows.Scan(bean)
		}
		if err != nil {
			return err
		}
		if err := containerValueSetFunc(&newValue, nil); err != nil {
			return err
		}
	}
	return rows.Err()
}
func (session *Session) cacheFind(t reflect.Type, sqlStr string, rowsSlicePtr interface{}, args ...interface{}) (err error) {
	if !session.canCache() ||
		utils.IndexNoCase(sqlStr, "having") != -1 ||
		utils.IndexNoCase(sqlStr, "group by") != -1 {
		return ErrCacheFailed
	}
	tableName := session.statement.TableName()
	cacher := session.engine.cacherMgr.GetCacher(tableName)
	if cacher == nil {
		return nil
	}
	for _, filter := range session.engine.dialect.Filters() {
		sqlStr = filter.Do(sqlStr)
	}
	newsql := session.statement.ConvertIDSQL(sqlStr)
	if newsql == "" {
		return ErrCacheFailed
	}
	table := session.statement.RefTable
	ids, err := cache.GetCacheSql(cacher, tableName, newsql, args)
	if err != nil {
		rows, err := session.queryRows(newsql, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		var i int
		ids = make([]schemasvr.PK, 0)
		for rows.Next() {
			i++
			if i > 500 {
				session.engine.logger.Debugf("[cacheFind] ids length > 500, no cache")
				return ErrCacheFailed
			}
			var res = make([]string, len(table.PrimaryKeys))
			err = rows.ScanSlice(&res)
			if err != nil {
				return err
			}
			var pk schemasvr.PK = make([]interface{}, len(table.PrimaryKeys))
			for i, col := range table.PKColumns() {
				pk[i], err = col.ConvertID(res[i])
				if err != nil {
					return err
				}
			}
			ids = append(ids, pk)
		}
		if rows.Err() != nil {
			return rows.Err()
		}
		session.engine.logger.Debugf("[cache] cache sql: %v, %v, %v, %v, %v", ids, tableName, sqlStr, newsql, args)
		err = cache.PutCacheSql(cacher, ids, tableName, newsql, args)
		if err != nil {
			return err
		}
	} else {
		session.engine.logger.Debugf("[cache] cache hit sql: %v, %v, %v, %v", tableName, sqlStr, newsql, args)
	}
	sliceValue := reflect.Indirect(reflect.ValueOf(rowsSlicePtr))
	ididxes := make(map[string]int)
	var ides []schemasvr.PK
	var temps = make([]interface{}, len(ids))
	for idx, id := range ids {
		sid, err := id.ToString()
		if err != nil {
			return err
		}
		bean := cacher.GetBean(tableName, sid)
		// fix issue #894
		isHit := func() (ht bool) {
			if bean == nil {
				ht = false
				return
			}
			ckb := reflect.ValueOf(bean).Elem().Type()
			ht = ckb == t
			if !ht && t.Kind() == reflect.Ptr {
				ht = t.Elem() == ckb
			}
			return
		}
		if !isHit() {
			ides = append(ides, id)
			ididxes[sid] = idx
		} else {
			session.engine.logger.Debugf("[cache] cache hit bean: %v, %v, %v", tableName, id, bean)
			pk, err := table.IDOfV(reflect.ValueOf(bean))
			if err != nil {
				return err
			}
			xid, err := pk.ToString()
			if err != nil {
				return err
			}
			if sid != xid {
				session.engine.logger.Errorf("[cache] error cache: %v, %v, %v", xid, sid, bean)
				return ErrCacheFailed
			}
			temps[idx] = bean
		}
	}
	if len(ides) > 0 {
		slices := reflect.New(reflect.SliceOf(t))
		beans := slices.Interface()
		statement := session.statement
		session.statement = statementsvr.NewStatement(
			session.engine.dialect,
			session.engine.tagParser,
			session.engine.DatabaseTZ,
		)
		if len(table.PrimaryKeys) == 1 {
			ff := make([]interface{}, 0, len(ides))
			for _, ie := range ides {
				ff = append(ff, ie[0])
			}
			session.In("`"+table.PrimaryKeys[0]+"`", ff...)
		} else {
			for _, ie := range ides {
				cond := builder.NewCond()
				for i, name := range table.PrimaryKeys {
					cond = cond.And(builder.Eq{"`" + name + "`": ie[i]})
				}
				session.Or(cond)
			}
		}
		err = session.NoCache().Table(tableName).find(beans)
		if err != nil {
			return err
		}
		session.statement = statement
		vs := reflect.Indirect(reflect.ValueOf(beans))
		for i := 0; i < vs.Len(); i++ {
			rv := vs.Index(i)
			if rv.Kind() != reflect.Ptr {
				rv = rv.Addr()
			}
			id, err := table.IDOfV(rv)
			if err != nil {
				return err
			}
			sid, err := id.ToString()
			if err != nil {
				return err
			}
			bean := rv.Interface()
			temps[ididxes[sid]] = bean
			session.engine.logger.Debugf("[cache] cache bean: %v, %v, %v, %v", tableName, id, bean, temps)
			cacher.PutBean(tableName, sid, bean)
		}
	}
	for j := 0; j < len(temps); j++ {
		bean := temps[j]
		if bean == nil {
			session.engine.logger.Warnf("[cache] cache no hit: %v, %v, %v", tableName, ids[j], temps)
			// return errors.New("cache error") // !nashtsai! no need to return error, but continue instead
			continue
		}
		if sliceValue.Kind() == reflect.Slice {
			if t.Kind() == reflect.Ptr {
				sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(bean)))
			} else {
				sliceValue.Set(reflect.Append(sliceValue, reflect.Indirect(reflect.ValueOf(bean))))
			}
		} else if sliceValue.Kind() == reflect.Map {
			var key = ids[j]
			keyType := sliceValue.Type().Key()
			keyValue := reflect.New(keyType)
			var ikey interface{}
			if len(key) == 1 {
				if err := convert.AssignValue(keyValue, key[0]); err != nil {
					return err
				}
				ikey = keyValue.Elem().Interface()
			} else {
				if keyType.Kind() != reflect.Slice {
					return errors.New("table have multiple primary keys, key is not schemas.PK or slice")
				}
				ikey = key
			}
			if t.Kind() == reflect.Ptr {
				sliceValue.SetMapIndex(reflect.ValueOf(ikey), reflect.ValueOf(bean))
			} else {
				sliceValue.SetMapIndex(reflect.ValueOf(ikey), reflect.Indirect(reflect.ValueOf(bean)))
			}
		}
	}
	return nil
}
