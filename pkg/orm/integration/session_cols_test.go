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

	schemasvr "github.com/bhojpur/dbm/pkg/orm/schema"
	"github.com/bhojpur/sql/pkg/builder"
	"github.com/stretchr/testify/assert"
)

func TestSetExpr(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type UserExprIssue struct {
		Id    int64
		Title string
	}
	assert.NoError(t, testEngine.Sync(new(UserExprIssue)))
	var issue = UserExprIssue{
		Title: "my issue",
	}
	cnt, err := testEngine.Insert(&issue)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	assert.EqualValues(t, 1, issue.Id)
	type UserExpr struct {
		Id      int64
		IssueId int64 `orm:"index"`
		Show    bool
	}
	assert.NoError(t, testEngine.Sync(new(UserExpr)))
	cnt, err = testEngine.Insert(&UserExpr{
		Show: true,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	var not = "NOT"
	if testEngine.Dialect().URI().DBType == schemasvr.MSSQL || testEngine.Dialect().URI().DBType == schemasvr.DAMENG {
		not = "~"
	}
	cnt, err = testEngine.SetExpr("show", not+" `show`").ID(1).Update(new(UserExpr))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	tableName := testEngine.TableName(new(UserExprIssue), true)
	cnt, err = testEngine.SetExpr("issue_id",
		builder.Select("`id`").
			From(testEngine.Quote(tableName)).
			Where(builder.Eq{"`id`": issue.Id})).
		ID(1).
		Update(new(UserExpr))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}
func TestCols(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type ColsTable struct {
		Id   int64
		Col1 string
		Col2 string
	}
	assertSync(t, new(ColsTable))
	_, err := testEngine.Insert(&ColsTable{
		Col1: "1",
		Col2: "2",
	})
	assert.NoError(t, err)
	sess := testEngine.ID(1)
	_, err = sess.Cols("col1").Cols("col2").Update(&ColsTable{
		Col1: "",
		Col2: "",
	})
	assert.NoError(t, err)
	var tb ColsTable
	has, err := testEngine.ID(1).Get(&tb)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, "", tb.Col1)
	assert.EqualValues(t, "", tb.Col2)
}
func TestMustCol(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type CustomerUpdate struct {
		Id        int64  `form:"id" json:"id"`
		Username  string `form:"username" json:"username" binding:"required"`
		Email     string `form:"email" json:"email"`
		Sex       int    `form:"sex" json:"sex"`
		Name      string `form:"name" json:"name" binding:"required"`
		Telephone string `form:"telephone" json:"telephone"`
		Type      int    `form:"type" json:"type" binding:"required"`
		ParentId  int64  `form:"parent_id" json:"parent_id" orm:"int null"`
		Remark    string `form:"remark" json:"remark"`
		Status    int    `form:"status" json:"status" binding:"required"`
		Age       int    `form:"age" json:"age"`
		CreatedAt int64  `orm:"created" form:"created_at" json:"created_at"`
		UpdatedAt int64  `orm:"updated" form:"updated_at" json:"updated_at"`
		BirthDate int64  `form:"birth_date" json:"birth_date"`
		Password  string `orm:"varchar(200)" form:"password" json:"password"`
	}
	assertSync(t, new(CustomerUpdate))
	var customer = CustomerUpdate{
		ParentId: 1,
	}
	cnt, err := testEngine.Insert(&customer)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	type CustomerOnlyId struct {
		Id int64
	}
	customer.ParentId = 0
	affected, err := testEngine.MustCols("parent_id").Update(&customer, &CustomerOnlyId{Id: customer.Id})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affected)
}
