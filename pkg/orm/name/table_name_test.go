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
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Userinfo struct {
	Uid        int64  `orm:"id pk not null autoincr"`
	Username   string `orm:"unique"`
	Departname string
	Alias      string `orm:"-"`
	Created    time.Time
	Detail     Userdetail `orm:"detail_id int(11)"`
	Height     float64
	Avatar     []byte
	IsMan      bool
}
type Userdetail struct {
	Id      int64
	Intro   string `orm:"text"`
	Profile string `orm:"varchar(2000)"`
}
type MyGetCustomTableImpletation struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

const getCustomTableName = "GetCustomTableInterface"

func (MyGetCustomTableImpletation) TableName() string {
	return getCustomTableName
}

type TestTableNameStruct struct{}

const getTestTableName = "my_test_table_name_struct"

func (t *TestTableNameStruct) TableName() string {
	return getTestTableName
}
func TestGetTableName(t *testing.T) {
	var kases = []struct {
		mapper            Mapper
		v                 reflect.Value
		expectedTableName string
	}{
		{
			SnakeMapper{},
			reflect.ValueOf(new(Userinfo)),
			"userinfo",
		},
		{
			SnakeMapper{},
			reflect.ValueOf(Userinfo{}),
			"userinfo",
		},
		{
			SameMapper{},
			reflect.ValueOf(new(Userinfo)),
			"Userinfo",
		},
		{
			SameMapper{},
			reflect.ValueOf(Userinfo{}),
			"Userinfo",
		},
		{
			SnakeMapper{},
			reflect.ValueOf(new(MyGetCustomTableImpletation)),
			getCustomTableName,
		},
		{
			SnakeMapper{},
			reflect.ValueOf(MyGetCustomTableImpletation{}),
			getCustomTableName,
		},
		{
			SnakeMapper{},
			reflect.ValueOf(new(TestTableNameStruct)),
			new(TestTableNameStruct).TableName(),
		},
		{
			SnakeMapper{},
			reflect.ValueOf(new(TestTableNameStruct)),
			getTestTableName,
		},
		{
			SnakeMapper{},
			reflect.ValueOf(TestTableNameStruct{}),
			getTestTableName,
		},
	}
	for _, kase := range kases {
		assert.EqualValues(t, kase.expectedTableName, GetTableName(kase.mapper, kase.v))
	}
}

type OAuth2Application struct {
}

// TableName sets the table name to `oauth2_application`
func (app *OAuth2Application) TableName() string {
	return "oauth2_application"
}
func TestGonicMapperCustomTable(t *testing.T) {
	assert.EqualValues(t, "oauth2_application",
		GetTableName(LintGonicMapper, reflect.ValueOf(new(OAuth2Application))))
	assert.EqualValues(t, "oauth2_application",
		GetTableName(LintGonicMapper, reflect.ValueOf(OAuth2Application{})))
}

type MyTable struct {
	Idx int
}

func (t *MyTable) TableName() string {
	return fmt.Sprintf("mytable_%d", t.Idx)
}
func TestMyTable(t *testing.T) {
	var table MyTable
	for i := 0; i < 10; i++ {
		table.Idx = i
		assert.EqualValues(t, fmt.Sprintf("mytable_%d", i), GetTableName(SameMapper{}, reflect.ValueOf(&table)))
	}
}
