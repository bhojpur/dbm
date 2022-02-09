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
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClose(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	sess1 := testEngine.NewSession()
	sess1.Close()
	assert.True(t, sess1.IsClosed())
	sess2 := testEngine.Where("`a` = ?", 1)
	sess2.Close()
	assert.True(t, sess2.IsClosed())
}
func TestNullFloatStruct(t *testing.T) {
	type MyNullFloat64 sql.NullFloat64
	type MyNullFloatStruct struct {
		Uuid   string
		Amount MyNullFloat64
	}
	assert.NoError(t, PrepareEngine())
	assert.NoError(t, testEngine.Sync(new(MyNullFloatStruct)))
	_, err := testEngine.Insert(&MyNullFloatStruct{
		Uuid: "111111",
		Amount: MyNullFloat64(sql.NullFloat64{
			Float64: 0.1111,
			Valid:   true,
		}),
	})
	assert.NoError(t, err)
}
func TestMustLogSQL(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	testEngine.ShowSQL(false)
	defer testEngine.ShowSQL(true)
	assertSync(t, new(Userinfo))
	_, err := testEngine.Table("userinfo").MustLogSQL(true).Get(new(Userinfo))
	assert.NoError(t, err)
}
func TestEnableSessionId(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	testEngine.EnableSessionID(true)
	assertSync(t, new(Userinfo))
	_, err := testEngine.Table("userinfo").MustLogSQL(true).Get(new(Userinfo))
	assert.NoError(t, err)
}
