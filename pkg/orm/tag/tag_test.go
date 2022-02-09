package tags

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

func TestSplitTag(t *testing.T) {
	var cases = []struct {
		tag  string
		tags []tag
	}{
		{"not null default '2000-01-01 00:00:00' TIMESTAMP", []tag{
			{
				name: "not",
			},
			{
				name: "null",
			},
			{
				name: "default",
			},
			{
				name: "'2000-01-01 00:00:00'",
			},
			{
				name: "TIMESTAMP",
			},
		},
		},
		{"TEXT", []tag{
			{
				name: "TEXT",
			},
		},
		},
		{"default('2000-01-01 00:00:00')", []tag{
			{
				name: "default",
				params: []string{
					"'2000-01-01 00:00:00'",
				},
			},
		},
		},
		{"json  binary", []tag{
			{
				name: "json",
			},
			{
				name: "binary",
			},
		},
		},
		{"numeric(10, 2)", []tag{
			{
				name:   "numeric",
				params: []string{"10", "2"},
			},
		},
		},
		{"numeric(10, 2) notnull", []tag{
			{
				name:   "numeric",
				params: []string{"10", "2"},
			},
			{
				name: "notnull",
			},
		},
		},
	}
	for _, kase := range cases {
		t.Run(kase.tag, func(t *testing.T) {
			tags, err := splitTag(kase.tag)
			assert.NoError(t, err)
			assert.EqualValues(t, len(tags), len(kase.tags))
			for i := 0; i < len(tags); i++ {
				assert.Equal(t, tags[i], kase.tags[i])
			}
		})
	}
}
