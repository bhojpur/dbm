package core

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

type quoterOnly struct {
	Dialect
}

func (q *quoterOnly) Quote(item string) string {
	return "[" + item + "]"
}
func TestQuoteFilter_Do(t *testing.T) {
	f := QuoteFilter{}
	sql := "SELECT `COLUMN_NAME` FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ? AND `COLUMN_NAME` = ?"
	res := f.Do(sql, new(quoterOnly), nil)
	assert.EqualValues(t,
		"SELECT [COLUMN_NAME] FROM [INFORMATION_SCHEMA].[COLUMNS] WHERE [TABLE_SCHEMA] = ? AND [TABLE_NAME] = ? AND [COLUMN_NAME] = ?",
		res,
	)
}
func TestSeqFilter(t *testing.T) {
	var kases = map[string]string{
		"SELECT * FROM TABLE1 WHERE a=? AND b=?":                               "SELECT * FROM TABLE1 WHERE a=$1 AND b=$2",
		"SELECT 1, '???', '2006-01-02 15:04:05' FROM TABLE1 WHERE a=? AND b=?": "SELECT 1, '???', '2006-01-02 15:04:05' FROM TABLE1 WHERE a=$1 AND b=$2",
		"select '1''?' from issue":                                             "select '1''?' from issue",
		"select '1\\??' from issue":                                            "select '1\\??' from issue",
		"select '1\\\\',? from issue":                                          "select '1\\\\',$1 from issue",
		"select '1\\''?',? from issue":                                         "select '1\\''?',$1 from issue",
	}
	for sql, result := range kases {
		assert.EqualValues(t, result, convertQuestionMark(sql, "$", 1))
	}
}
