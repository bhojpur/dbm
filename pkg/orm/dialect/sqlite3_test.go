package dialect

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

func TestSplitColStr(t *testing.T) {
	var kases = []struct {
		colStr string
		fields []string
	}{
		{
			colStr: "`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL",
			fields: []string{
				"`id`", "INTEGER", "PRIMARY", "KEY", "AUTOINCREMENT", "NOT", "NULL",
			},
		},
		{
			colStr: "`created` DATETIME DEFAULT '2006-01-02 15:04:05' NULL",
			fields: []string{
				"`created`", "DATETIME", "DEFAULT", "'2006-01-02 15:04:05'", "NULL",
			},
		},
	}
	for _, kase := range kases {
		assert.EqualValues(t, kase.fields, splitColStr(kase.colStr))
	}
}
