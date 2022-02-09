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
	"time"

	"github.com/bhojpur/dbm/pkg/orm/cache"
	"github.com/stretchr/testify/assert"
)

func TestCacheFind(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type MailBox struct {
		Id       int64 `orm:"pk"`
		Username string
		Password string
	}
	oldCacher := testEngine.GetDefaultCacher()
	cacher := cache.NewLRUCacher2(cache.NewMemoryStore(), time.Hour, 10000)
	testEngine.SetDefaultCacher(cacher)
	assert.NoError(t, testEngine.Sync(new(MailBox)))
	var inserts = []*MailBox{
		{
			Id:       0,
			Username: "user1",
			Password: "pass1",
		},
		{
			Id:       1,
			Username: "user2",
			Password: "pass2",
		},
	}
	_, err := testEngine.Insert(inserts[0], inserts[1])
	assert.NoError(t, err)
	var boxes []MailBox
	assert.NoError(t, testEngine.Find(&boxes))
	assert.EqualValues(t, 2, len(boxes))
	for i, box := range boxes {
		assert.Equal(t, inserts[i].Id, box.Id)
		assert.Equal(t, inserts[i].Username, box.Username)
		assert.Equal(t, inserts[i].Password, box.Password)
	}
	boxes = make([]MailBox, 0, 2)
	assert.NoError(t, testEngine.Find(&boxes))
	assert.EqualValues(t, 2, len(boxes))
	for i, box := range boxes {
		assert.Equal(t, inserts[i].Id, box.Id)
		assert.Equal(t, inserts[i].Username, box.Username)
		assert.Equal(t, inserts[i].Password, box.Password)
	}
	boxes = make([]MailBox, 0, 2)
	assert.NoError(t, testEngine.Alias("a").Where("`a`.`id`> -1").
		Asc("`a`.`id`").Find(&boxes))
	assert.EqualValues(t, 2, len(boxes))
	for i, box := range boxes {
		assert.Equal(t, inserts[i].Id, box.Id)
		assert.Equal(t, inserts[i].Username, box.Username)
		assert.Equal(t, inserts[i].Password, box.Password)
	}
	type MailBox4 struct {
		Id       int64
		Username string
		Password string
	}
	boxes2 := make([]MailBox4, 0, 2)
	assert.NoError(t, testEngine.Table("mail_box").Where("`mail_box`.`id` > -1").
		Asc("mail_box.id").Find(&boxes2))
	assert.EqualValues(t, 2, len(boxes2))
	for i, box := range boxes2 {
		assert.Equal(t, inserts[i].Id, box.Id)
		assert.Equal(t, inserts[i].Username, box.Username)
		assert.Equal(t, inserts[i].Password, box.Password)
	}
	testEngine.SetDefaultCacher(oldCacher)
}
func TestCacheFind2(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type MailBox2 struct {
		Id       uint64 `orm:"pk"`
		Username string
		Password string
	}
	oldCacher := testEngine.GetDefaultCacher()
	cacher := cache.NewLRUCacher2(cache.NewMemoryStore(), time.Hour, 10000)
	testEngine.SetDefaultCacher(cacher)
	assert.NoError(t, testEngine.Sync(new(MailBox2)))
	var inserts = []*MailBox2{
		{
			Id:       0,
			Username: "user1",
			Password: "pass1",
		},
		{
			Id:       1,
			Username: "user2",
			Password: "pass2",
		},
	}
	_, err := testEngine.Insert(inserts[0], inserts[1])
	assert.NoError(t, err)
	var boxes []MailBox2
	assert.NoError(t, testEngine.Find(&boxes))
	assert.EqualValues(t, 2, len(boxes))
	for i, box := range boxes {
		assert.Equal(t, inserts[i].Id, box.Id)
		assert.Equal(t, inserts[i].Username, box.Username)
		assert.Equal(t, inserts[i].Password, box.Password)
	}
	boxes = make([]MailBox2, 0, 2)
	assert.NoError(t, testEngine.Find(&boxes))
	assert.EqualValues(t, 2, len(boxes))
	for i, box := range boxes {
		assert.Equal(t, inserts[i].Id, box.Id)
		assert.Equal(t, inserts[i].Username, box.Username)
		assert.Equal(t, inserts[i].Password, box.Password)
	}
	testEngine.SetDefaultCacher(oldCacher)
}
func TestCacheGet(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type MailBox3 struct {
		Id       uint64
		Username string
		Password string
	}
	oldCacher := testEngine.GetDefaultCacher()
	cacher := cache.NewLRUCacher2(cache.NewMemoryStore(), time.Hour, 10000)
	testEngine.SetDefaultCacher(cacher)
	assert.NoError(t, testEngine.Sync(new(MailBox3)))
	var inserts = []*MailBox3{
		{
			Username: "user1",
			Password: "pass1",
		},
	}
	_, err := testEngine.Insert(inserts[0])
	assert.NoError(t, err)
	var box1 MailBox3
	has, err := testEngine.Where("`id` = ?", inserts[0].Id).Get(&box1)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, "user1", box1.Username)
	assert.EqualValues(t, "pass1", box1.Password)
	var box2 MailBox3
	has, err = testEngine.Where("`id` = ?", inserts[0].Id).Get(&box2)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, "user1", box2.Username)
	assert.EqualValues(t, "pass1", box2.Password)
	testEngine.SetDefaultCacher(oldCacher)
}
