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

// Filter is an interface to filter SQL
type Filter interface {
	Do(sql string, dialect Dialect, table *Table) string
}

// QuoteFilter filter SQL replace ` to database's own quote character
type QuoteFilter struct {
}

func (s *QuoteFilter) Do(sql string, dialect Dialect, table *Table) string {
	dummy := dialect.Quote("")
	if len(dummy) != 2 {
		return sql
	}
	prefix, suffix := dummy[0], dummy[1]
	raw := []byte(sql)
	for i, cnt := 0, 0; i < len(raw); i = i + 1 {
		if raw[i] == '`' {
			if cnt%2 == 0 {
				raw[i] = prefix
			} else {
				raw[i] = suffix
			}
			cnt++
		}
	}
	return string(raw)
}

// IdFilter filter SQL replace (id) to primary key column name
type IdFilter struct {
}
type Quoter struct {
	dialect Dialect
}

func NewQuoter(dialect Dialect) *Quoter {
	return &Quoter{dialect}
}
func (q *Quoter) Quote(content string) string {
	return q.dialect.Quote(content)
}
func (i *IdFilter) Do(sql string, dialect Dialect, table *Table) string {
	quoter := NewQuoter(dialect)
	if table != nil && len(table.PrimaryKeys) == 1 {
		sql = strings.Replace(sql, " `(id)` ", " "+quoter.Quote(table.PrimaryKeys[0])+" ", -1)
		sql = strings.Replace(sql, " "+quoter.Quote("(id)")+" ", " "+quoter.Quote(table.PrimaryKeys[0])+" ", -1)
		return strings.Replace(sql, " (id) ", " "+quoter.Quote(table.PrimaryKeys[0])+" ", -1)
	}
	return sql
}

// SeqFilter filter SQL replace ?, ? ... to $1, $2 ...
type SeqFilter struct {
	Prefix string
	Start  int
}

func convertQuestionMark(sql, prefix string, start int) string {
	var buf strings.Builder
	var beginSingleQuote bool
	var index = start
	for _, c := range sql {
		if !beginSingleQuote && c == '?' {
			buf.WriteString(fmt.Sprintf("%s%v", prefix, index))
			index++
		} else {
			if c == '\'' {
				beginSingleQuote = !beginSingleQuote
			}
			buf.WriteRune(c)
		}
	}
	return buf.String()
}
func (s *SeqFilter) Do(sql string, dialect Dialect, table *Table) string {
	return convertQuestionMark(sql, s.Prefix, s.Start)
}
