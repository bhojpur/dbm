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

func TestIterate(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type UserIterate struct {
		Id    int64
		IsMan bool
	}
	assert.NoError(t, testEngine.Sync(new(UserIterate)))
	cnt, err := testEngine.Insert(&UserIterate{
		IsMan: true,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	cnt, err = testEngine.Insert(&UserIterate{
		IsMan: false,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	cnt = 0
	err = testEngine.Iterate(new(UserIterate), func(i int, bean interface{}) error {
		user := bean.(*UserIterate)
		if cnt == 0 {
			assert.EqualValues(t, 1, user.Id)
			assert.EqualValues(t, true, user.IsMan)
		} else {
			assert.EqualValues(t, 2, user.Id)
			assert.EqualValues(t, false, user.IsMan)
		}
		cnt++
		return nil
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 2, cnt)
}
func TestBufferIterate(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type UserBufferIterate struct {
		Id    int64
		IsMan bool
	}
	assert.NoError(t, testEngine.Sync(new(UserBufferIterate)))
	var size = 20
	for i := 0; i < size; i++ {
		cnt, err := testEngine.Insert(&UserBufferIterate{
			IsMan: true,
		})
		assert.NoError(t, err)
		assert.EqualValues(t, 1, cnt)
	}
	var cnt = 0
	err := testEngine.BufferSize(9).Iterate(new(UserBufferIterate), func(i int, bean interface{}) error {
		user := bean.(*UserBufferIterate)
		assert.EqualValues(t, cnt+1, user.Id)
		assert.EqualValues(t, true, user.IsMan)
		cnt++
		return nil
	})
	assert.NoError(t, err)
	assert.EqualValues(t, size, cnt)
	cnt = 0
	err = testEngine.Limit(20).BufferSize(9).Iterate(new(UserBufferIterate), func(i int, bean interface{}) error {
		user := bean.(*UserBufferIterate)
		assert.EqualValues(t, cnt+1, user.Id)
		assert.EqualValues(t, true, user.IsMan)
		cnt++
		return nil
	})
	assert.NoError(t, err)
	assert.EqualValues(t, size, cnt)
	cnt = 0
	err = testEngine.Limit(7).BufferSize(9).Iterate(new(UserBufferIterate), func(i int, bean interface{}) error {
		user := bean.(*UserBufferIterate)
		assert.EqualValues(t, cnt+1, user.Id)
		assert.EqualValues(t, true, user.IsMan)
		cnt++
		return nil
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 7, cnt)
	cnt = 0
	err = testEngine.Where("`id` <= 10").BufferSize(2).Iterate(new(UserBufferIterate), func(i int, bean interface{}) error {
		user := bean.(*UserBufferIterate)
		assert.EqualValues(t, cnt+1, user.Id)
		assert.EqualValues(t, true, user.IsMan)
		cnt++
		return nil
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 10, cnt)
}
