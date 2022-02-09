package integration

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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRows(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type UserRows struct {
		Id    int64
		IsMan bool
	}
	assert.NoError(t, testEngine.Sync(new(UserRows)))
	cnt, err := testEngine.Insert(&UserRows{
		IsMan: true,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	rows, err := testEngine.Rows(new(UserRows))
	assert.NoError(t, err)
	defer rows.Close()
	cnt = 0
	user := new(UserRows)
	for rows.Next() {
		err = rows.Scan(user)
		assert.NoError(t, err)
		cnt++
	}
	assert.EqualValues(t, 1, cnt)
	assert.False(t, rows.Next())
	assert.NoError(t, rows.Close())
	rows0, err := testEngine.Where("1>1").Rows(new(UserRows))
	assert.NoError(t, err)
	defer rows0.Close()
	cnt = 0
	user0 := new(UserRows)
	for rows0.Next() {
		err = rows0.Scan(user0)
		assert.NoError(t, err)
		cnt++
	}
	assert.EqualValues(t, 0, cnt)
	assert.NoError(t, rows0.Close())
	sess := testEngine.NewSession()
	defer sess.Close()
	rows1, err := sess.Prepare().Rows(new(UserRows))
	assert.NoError(t, err)
	defer rows1.Close()
	cnt = 0
	for rows1.Next() {
		err = rows1.Scan(user)
		assert.NoError(t, err)
		cnt++
	}
	assert.EqualValues(t, 1, cnt)
	var tbName = testEngine.Quote(testEngine.TableName(user, true))
	rows2, err := testEngine.SQL("SELECT * FROM " + tbName).Rows(new(UserRows))
	assert.NoError(t, err)
	defer rows2.Close()
	cnt = 0
	for rows2.Next() {
		err = rows2.Scan(user)
		assert.NoError(t, err)
		cnt++
	}
	assert.EqualValues(t, 1, cnt)
}
func TestRowsMyTableName(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type UserRowsMyTable struct {
		Id    int64
		IsMan bool
	}
	var tableName = "user_rows_my_table_name"
	assert.NoError(t, testEngine.Table(tableName).Sync(new(UserRowsMyTable)))
	cnt, err := testEngine.Table(tableName).Insert(&UserRowsMyTable{
		IsMan: true,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	rows, err := testEngine.Table(tableName).Rows(new(UserRowsMyTable))
	assert.NoError(t, err)
	cnt = 0
	user := new(UserRowsMyTable)
	for rows.Next() {
		err = rows.Scan(user)
		assert.NoError(t, err)
		cnt++
	}
	assert.EqualValues(t, 1, cnt)
	rows.Close()
	rows, err = testEngine.Table(tableName).Rows(&UserRowsMyTable{
		Id: 2,
	})
	assert.NoError(t, err)
	cnt = 0
	user = new(UserRowsMyTable)
	for rows.Next() {
		err = rows.Scan(user)
		assert.NoError(t, err)
		cnt++
	}
	assert.EqualValues(t, 0, cnt)
}

type UserRowsSpecTable struct {
	Id    int64
	IsMan bool
}

func (UserRowsSpecTable) TableName() string {
	return "user_rows_my_table_name"
}
func TestRowsSpecTableName(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assert.NoError(t, testEngine.Sync(new(UserRowsSpecTable)))
	cnt, err := testEngine.Insert(&UserRowsSpecTable{
		IsMan: true,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	rows, err := testEngine.Rows(new(UserRowsSpecTable))
	assert.NoError(t, err)
	defer rows.Close()
	cnt = 0
	user := new(UserRowsSpecTable)
	for rows.Next() {
		err = rows.Scan(user)
		assert.NoError(t, err)
		cnt++
	}
	assert.NoError(t, rows.Err())
	assert.EqualValues(t, 1, cnt)
}
func TestRowsScanVars(t *testing.T) {
	type RowsScanVars struct {
		Id   int64
		Name string
		Age  int
	}
	assert.NoError(t, PrepareEngine())
	assert.NoError(t, testEngine.Sync2(new(RowsScanVars)))
	cnt, err := testEngine.Insert(&RowsScanVars{
		Name: "xlw",
		Age:  42,
	}, &RowsScanVars{
		Name: "xlw2",
		Age:  24,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 2, cnt)
	rows, err := testEngine.Cols("name", "age").Rows(new(RowsScanVars))
	assert.NoError(t, err)
	defer rows.Close()
	cnt = 0
	for rows.Next() {
		var name string
		var age int
		err = rows.Scan(&name, &age)
		assert.NoError(t, err)
		if cnt == 0 {
			assert.EqualValues(t, "xlw", name)
			assert.EqualValues(t, 42, age)
		} else if cnt == 1 {
			assert.EqualValues(t, "xlw2", name)
			assert.EqualValues(t, 24, age)
		}
		cnt++
	}
	assert.NoError(t, rows.Err())
	assert.EqualValues(t, 2, cnt)
}
