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
	"fmt"
	"strings"
)

// enumerate all index types
const (
	IndexType = iota + 1
	UniqueType
)

// Index represents a database index
type Index struct {
	IsRegular bool
	Name      string
	Type      int
	Cols      []string
}

func (index *Index) XName(tableName string) string {
	if !strings.HasPrefix(index.Name, "UQE_") &&
		!strings.HasPrefix(index.Name, "IDX_") {
		tableParts := strings.Split(strings.Replace(tableName, `"`, "", -1), ".")
		tableName = tableParts[len(tableParts)-1]
		if index.Type == UniqueType {
			return fmt.Sprintf("UQE_%v_%v", tableName, index.Name)
		}
		return fmt.Sprintf("IDX_%v_%v", tableName, index.Name)
	}
	return index.Name
}

// AddColumn add columns which will be composite index
func (index *Index) AddColumn(cols ...string) {
	for _, col := range cols {
		index.Cols = append(index.Cols, col)
	}
}
func (index *Index) Equal(dst *Index) bool {
	if index.Type != dst.Type {
		return false
	}
	if len(index.Cols) != len(dst.Cols) {
		return false
	}
	for i := 0; i < len(index.Cols); i++ {
		var found bool
		for j := 0; j < len(dst.Cols); j++ {
			if index.Cols[i] == dst.Cols[j] {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// NewIndex new an index object
func NewIndex(name string, indexType int) *Index {
	return &Index{true, name, indexType, make([]string, 0)}
}
