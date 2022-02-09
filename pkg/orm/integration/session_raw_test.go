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
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExecAndQuery(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type UserinfoQuery struct {
		Uid  int
		Name string
	}
	assert.NoError(t, testEngine.Sync(new(UserinfoQuery)))
	res, err := testEngine.Exec("INSERT INTO "+testEngine.TableName("`userinfo_query`", true)+" (`uid`, `name`) VALUES (?, ?)", 1, "user")
	assert.NoError(t, err)
	cnt, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	results, err := testEngine.Query("select * from " + testEngine.Quote(testEngine.TableName("userinfo_query", true)))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(results))
	id, err := strconv.Atoi(string(results[0]["uid"]))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, id)
	assert.Equal(t, "user", string(results[0]["name"]))
}
func TestExecTime(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type UserinfoExecTime struct {
		Uid     int
		Name    string
		Created time.Time
	}
	assert.NoError(t, testEngine.Sync(new(UserinfoExecTime)))
	now := time.Now()
	res, err := testEngine.Exec("INSERT INTO "+testEngine.TableName("`userinfo_exec_time`", true)+" (`uid`, `name`, `created`) VALUES (?, ?, ?)", 1, "user", now)
	assert.NoError(t, err)
	cnt, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	results, err := testEngine.QueryString("SELECT * FROM " + testEngine.Quote(testEngine.TableName("userinfo_exec_time", true)))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(results))
	assert.EqualValues(t, now.In(testEngine.GetTZLocation()).Format("2006-01-02 15:04:05"), results[0]["created"])
	var uet UserinfoExecTime
	has, err := testEngine.Where("`uid`=?", 1).Get(&uet)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, now.In(testEngine.GetTZLocation()).Format("2006-01-02 15:04:05"), uet.Created.Format("2006-01-02 15:04:05"))
}
