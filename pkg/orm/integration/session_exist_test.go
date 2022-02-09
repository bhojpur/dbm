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
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExistStruct(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type RecordExist struct {
		Id   int64
		Name string
	}
	assertSync(t, new(RecordExist))
	has, err := testEngine.Exist(new(RecordExist))
	assert.NoError(t, err)
	assert.False(t, has)
	cnt, err := testEngine.Insert(&RecordExist{
		Name: "test1",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	has, err = testEngine.Exist(new(RecordExist))
	assert.NoError(t, err)
	assert.True(t, has)
	has, err = testEngine.Exist(&RecordExist{
		Name: "test1",
	})
	assert.NoError(t, err)
	assert.True(t, has)
	has, err = testEngine.Exist(&RecordExist{
		Name: "test2",
	})
	assert.NoError(t, err)
	assert.False(t, has)
	has, err = testEngine.Where("`name` = ?", "test1").Exist(&RecordExist{})
	assert.NoError(t, err)
	assert.True(t, has)
	has, err = testEngine.Where("`name` = ?", "test2").Exist(&RecordExist{})
	assert.NoError(t, err)
	assert.False(t, has)
	has, err = testEngine.SQL("select * from "+testEngine.Quote(testEngine.TableName("record_exist", true))+" where `name` = ?", "test1").Exist()
	assert.NoError(t, err)
	assert.True(t, has)
	has, err = testEngine.SQL("select * from "+testEngine.Quote(testEngine.TableName("record_exist", true))+" where `name` = ?", "test2").Exist()
	assert.NoError(t, err)
	assert.False(t, has)
	has, err = testEngine.Table("record_exist").Exist()
	assert.NoError(t, err)
	assert.True(t, has)
	has, err = testEngine.Table("record_exist").Where("`name` = ?", "test1").Exist()
	assert.NoError(t, err)
	assert.True(t, has)
	has, err = testEngine.Table("record_exist").Where("`name` = ?", "test2").Exist()
	assert.NoError(t, err)
	assert.False(t, has)
	has, err = testEngine.Table(new(RecordExist)).ID(1).Cols("id").Exist()
	assert.NoError(t, err)
	assert.True(t, has)
}
func TestExistStructForJoin(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type Number struct {
		Id  int64
		Lid int64
	}
	type OrderList struct {
		Id  int64
		Eid int64
	}
	type Player struct {
		Id   int64
		Name string
	}
	assert.NoError(t, testEngine.Sync(new(Number), new(OrderList), new(Player)))
	var ply Player
	cnt, err := testEngine.Insert(&ply)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	var orderlist = OrderList{
		Eid: ply.Id,
	}
	cnt, err = testEngine.Insert(&orderlist)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	var um = Number{
		Lid: orderlist.Id,
	}
	cnt, err = testEngine.Insert(&um)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	session := testEngine.NewSession()
	defer session.Close()
	session.Table("number").
		Join("INNER", "order_list", "`order_list`.`id` = `number`.`lid`").
		Join("LEFT", "player", "`player`.`id` = `order_list`.`eid`").
		Where("`number`.`lid` = ?", 1)
	has, err := session.Exist()
	assert.NoError(t, err)
	assert.True(t, has)
	session.Table("number").
		Join("INNER", "order_list", "`order_list`.`id` = `number`.`lid`").
		Join("LEFT", "player", "`player`.`id` = `order_list`.`eid`").
		Where("`number`.`lid` = ?", 2)
	has, err = session.Exist()
	assert.NoError(t, err)
	assert.False(t, has)
	session.Table("number").
		Select("`order_list`.`id`").
		Join("INNER", "order_list", "`order_list`.`id` = `number`.`lid`").
		Join("LEFT", "player", "`player`.`id` = `order_list`.`eid`").
		Where("`order_list`.`id` = ?", 1)
	has, err = session.Exist()
	assert.NoError(t, err)
	assert.True(t, has)
	session.Table("number").
		Select("player.id").
		Join("INNER", "order_list", "`order_list`.`id` = `number`.`lid`").
		Join("LEFT", "player", "`player`.`id` = `order_list`.`eid`").
		Where("`player`.`id` = ?", 2)
	has, err = session.Exist()
	assert.NoError(t, err)
	assert.False(t, has)
	session.Table("number").
		Select("player.id").
		Join("INNER", "order_list", "`order_list`.`id` = `number`.`lid`").
		Join("LEFT", "player", "`player`.`id` = `order_list`.`eid`")
	has, err = session.Exist()
	assert.NoError(t, err)
	assert.True(t, has)
	err = session.DropTable("order_list")
	assert.NoError(t, err)
	exist, err := session.IsTableExist("order_list")
	assert.NoError(t, err)
	assert.False(t, exist)
	session.Table("number").
		Select("player.id").
		Join("INNER", "order_list", "`order_list`.`id` = `number`.`lid`").
		Join("LEFT", "player", "`player`.`id` = `order_list`.`eid`")
	has, err = session.Exist()
	assert.Error(t, err)
	assert.False(t, has)
	session.Table("number").
		Select("player.id").
		Join("LEFT", "player", "`player`.`id` = `number`.`lid`")
	has, err = session.Exist()
	assert.NoError(t, err)
	assert.True(t, has)
}
func TestExistContext(t *testing.T) {
	type ContextQueryStruct struct {
		Id   int64
		Name string
	}
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(ContextQueryStruct))
	_, err := testEngine.Insert(&ContextQueryStruct{Name: "1"})
	assert.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	time.Sleep(time.Nanosecond)
	has, err := testEngine.Context(ctx).Exist(&ContextQueryStruct{Name: "1"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
	assert.False(t, has)
}
