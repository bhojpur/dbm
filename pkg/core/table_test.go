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
	"strings"
	"testing"
)

var testsGetColumn = []struct {
	name string
	idx  int
}{
	{"Id", 0},
	{"Deleted", 0},
	{"Caption", 0},
	{"Code_1", 0},
	{"Code_2", 0},
	{"Code_3", 0},
	{"Parent_Id", 0},
	{"Latitude", 0},
	{"Longitude", 0},
}
var table *Table

func init() {
	table = NewEmptyTable()
	var name string
	for i := 0; i < len(testsGetColumn); i++ {
		// as in Table.AddColumn func
		name = strings.ToLower(testsGetColumn[i].name)
		table.columnsMap[name] = append(table.columnsMap[name], &Column{})
	}
}
func TestGetColumn(t *testing.T) {
	for _, test := range testsGetColumn {
		if table.GetColumn(test.name) == nil {
			t.Error("Column not found!")
		}
	}
}
func TestGetColumnIdx(t *testing.T) {
	for _, test := range testsGetColumn {
		if table.GetColumnIdx(test.name, test.idx) == nil {
			t.Errorf("Column %s with idx %d not found!", test.name, test.idx)
		}
	}
}
func BenchmarkGetColumnWithToLower(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range testsGetColumn {
			if _, ok := table.columnsMap[strings.ToLower(test.name)]; !ok {
				b.Errorf("Column not found:%s", test.name)
			}
		}
	}
}
func BenchmarkGetColumnIdxWithToLower(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range testsGetColumn {
			if c, ok := table.columnsMap[strings.ToLower(test.name)]; ok {
				if test.idx < len(c) {
					continue
				} else {
					b.Errorf("Bad idx in: %s, %d", test.name, test.idx)
				}
			} else {
				b.Errorf("Column not found: %s, %d", test.name, test.idx)
			}
		}
	}
}
func BenchmarkGetColumn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range testsGetColumn {
			if table.GetColumn(test.name) == nil {
				b.Errorf("Column not found:%s", test.name)
			}
		}
	}
}
func BenchmarkGetColumnIdx(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range testsGetColumn {
			if table.GetColumnIdx(test.name, test.idx) == nil {
				b.Errorf("Column not found:%s, %d", test.name, test.idx)
			}
		}
	}
}
