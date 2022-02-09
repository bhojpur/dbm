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

func BenchmarkGetVars(b *testing.B) {
	b.StopTimer()
	assert.NoError(b, PrepareEngine())
	testEngine.ShowSQL(false)
	type BenchmarkGetVars struct {
		Id   int64
		Name string
	}
	assert.NoError(b, testEngine.Sync(new(BenchmarkGetVars)))
	var v = BenchmarkGetVars{
		Name: "myname",
	}
	_, err := testEngine.Insert(&v)
	assert.NoError(b, err)
	b.StartTimer()
	var myname string
	for i := 0; i < b.N; i++ {
		has, err := testEngine.Cols("name").Table("benchmark_get_vars").Where("`id`=?", v.Id).Get(&myname)
		b.StopTimer()
		myname = ""
		assert.True(b, has)
		assert.NoError(b, err)
		b.StartTimer()
	}
}
func BenchmarkGetStruct(b *testing.B) {
	b.StopTimer()
	assert.NoError(b, PrepareEngine())
	testEngine.ShowSQL(false)
	type BenchmarkGetStruct struct {
		Id   int64
		Name string
	}
	assert.NoError(b, testEngine.Sync(new(BenchmarkGetStruct)))
	var v = BenchmarkGetStruct{
		Name: "myname",
	}
	_, err := testEngine.Insert(&v)
	assert.NoError(b, err)
	b.StartTimer()
	var myname BenchmarkGetStruct
	for i := 0; i < b.N; i++ {
		has, err := testEngine.ID(v.Id).Get(&myname)
		b.StopTimer()
		myname.Id = 0
		myname.Name = ""
		assert.True(b, has)
		assert.NoError(b, err)
		b.StartTimer()
	}
}
func BenchmarkFindStruct(b *testing.B) {
	b.StopTimer()
	assert.NoError(b, PrepareEngine())
	testEngine.ShowSQL(false)
	type BenchmarkFindStruct struct {
		Id   int64
		Name string
	}
	assert.NoError(b, testEngine.Sync(new(BenchmarkFindStruct)))
	var v = BenchmarkFindStruct{
		Name: "myname",
	}
	_, err := testEngine.Insert(&v)
	assert.NoError(b, err)
	var mynames = make([]BenchmarkFindStruct, 0, 1)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := testEngine.Find(&mynames)
		b.StopTimer()
		mynames = make([]BenchmarkFindStruct, 0, 1)
		assert.NoError(b, err)
		b.StartTimer()
	}
}
