package cache

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

	"github.com/bhojpur/dbm/pkg/orm/schema"
	"github.com/stretchr/testify/assert"
)

func TestLRUCache(t *testing.T) {
	type CacheObject1 struct {
		Id int64
	}
	store := NewMemoryStore()
	cacher := NewLRUCacher(store, 10000)
	tableName := "cache_object1"
	pks := []schema.PK{
		{1},
		{2},
	}
	for _, pk := range pks {
		sid, err := pk.ToString()
		assert.NoError(t, err)
		cacher.PutIds(tableName, "select * from cache_object1", sid)
		ids := cacher.GetIds(tableName, "select * from cache_object1")
		assert.EqualValues(t, sid, ids)
		cacher.ClearIds(tableName)
		ids2 := cacher.GetIds(tableName, "select * from cache_object1")
		assert.Nil(t, ids2)
		obj2 := cacher.GetBean(tableName, sid)
		assert.Nil(t, obj2)
		var obj = new(CacheObject1)
		cacher.PutBean(tableName, sid, obj)
		obj3 := cacher.GetBean(tableName, sid)
		assert.EqualValues(t, obj, obj3)
		cacher.DelBean(tableName, sid)
		obj4 := cacher.GetBean(tableName, sid)
		assert.Nil(t, obj4)
	}
}
