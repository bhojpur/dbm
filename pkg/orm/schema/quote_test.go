package schema

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

	"github.com/stretchr/testify/assert"
)

func TestAlwaysQuoteTo(t *testing.T) {
	var (
		quoter = Quoter{'[', ']', AlwaysReserve}
		kases  = []struct {
			expected string
			value    string
		}{
			{"[mytable]", "mytable"},
			{"[mytable]", "`mytable`"},
			{"[mytable]", `[mytable]`},
			{`["mytable"]`, `"mytable"`},
			{`[mytable].*`, `[mytable].*`},
			{"[myschema].[mytable]", "myschema.mytable"},
			{"[myschema].[mytable]", "`myschema`.mytable"},
			{"[myschema].[mytable]", "myschema.`mytable`"},
			{"[myschema].[mytable]", "`myschema`.`mytable`"},
			{"[myschema].[mytable]", `[myschema].mytable`},
			{"[myschema].[mytable]", `myschema.[mytable]`},
			{"[myschema].[mytable]", `[myschema].[mytable]`},
			{`["myschema].[mytable"]`, `"myschema.mytable"`},
			{"[message_user] AS [sender]", "`message_user` AS `sender`"},
			{"[myschema].[mytable] AS [table]", "myschema.mytable AS table"},
			{" [mytable]", " mytable"},
			{"  [mytable]", "  mytable"},
			{"[mytable] ", "mytable "},
			{"[mytable]  ", "mytable  "},
			{" [mytable] ", " mytable "},
			{"  [mytable]  ", "  mytable  "},
		}
	)
	for _, v := range kases {
		t.Run(v.value, func(t *testing.T) {
			buf := &strings.Builder{}
			err := quoter.QuoteTo(buf, v.value)
			assert.NoError(t, err)
			assert.EqualValues(t, v.expected, buf.String())
		})
	}
}
func TestReversedQuoteTo(t *testing.T) {
	var (
		quoter = Quoter{'[', ']', func(s string) bool {
			return s == "mytable"
		}}
		kases = []struct {
			expected string
			value    string
		}{
			{"[mytable]", "mytable"},
			{"[mytable]", "`mytable`"},
			{"[mytable]", `[mytable]`},
			{"[mytable].*", `[mytable].*`},
			{`"mytable"`, `"mytable"`},
			{"myschema.[mytable]", "myschema.mytable"},
			{"myschema.[mytable]", "`myschema`.mytable"},
			{"myschema.[mytable]", "myschema.`mytable`"},
			{"myschema.[mytable]", "`myschema`.`mytable`"},
			{"myschema.[mytable]", `[myschema].mytable`},
			{"myschema.[mytable]", `myschema.[mytable]`},
			{"myschema.[mytable]", `[myschema].[mytable]`},
			{`"myschema.mytable"`, `"myschema.mytable"`},
			{"message_user AS sender", "`message_user` AS `sender`"},
			{"myschema.[mytable] AS table", "myschema.mytable AS table"},
		}
	)
	for _, v := range kases {
		t.Run(v.value, func(t *testing.T) {
			buf := &strings.Builder{}
			quoter.QuoteTo(buf, v.value)
			assert.EqualValues(t, v.expected, buf.String())
		})
	}
}
func TestNoQuoteTo(t *testing.T) {
	var (
		quoter = Quoter{'[', ']', AlwaysNoReserve}
		kases  = []struct {
			expected string
			value    string
		}{
			{"mytable", "mytable"},
			{"mytable", "`mytable`"},
			{"mytable", `[mytable]`},
			{"mytable.*", `[mytable].*`},
			{`"mytable"`, `"mytable"`},
			{"myschema.mytable", "myschema.mytable"},
			{"myschema.mytable", "`myschema`.mytable"},
			{"myschema.mytable", "myschema.`mytable`"},
			{"myschema.mytable", "`myschema`.`mytable`"},
			{"myschema.mytable", `[myschema].mytable`},
			{"myschema.mytable", `myschema.[mytable]`},
			{"myschema.mytable", `[myschema].[mytable]`},
			{`"myschema.mytable"`, `"myschema.mytable"`},
			{"message_user AS sender", "`message_user` AS `sender`"},
			{"myschema.mytable AS table", "myschema.mytable AS table"},
		}
	)
	for _, v := range kases {
		t.Run(v.value, func(t *testing.T) {
			buf := &strings.Builder{}
			err := quoter.QuoteTo(buf, v.value)
			assert.NoError(t, err)
			assert.EqualValues(t, v.expected, buf.String())
		})
	}
}
func TestJoin(t *testing.T) {
	cols := []string{"f1", "f2", "f3"}
	quoter := Quoter{'[', ']', AlwaysReserve}
	assert.EqualValues(t, "[a],[b]", quoter.Join([]string{"a", " b"}, ","))
	assert.EqualValues(t, "[a].*,[b].[c]", quoter.Join([]string{"a.*", " b.c"}, ","))
	assert.EqualValues(t, "[f1], [f2], [f3]", quoter.Join(cols, ", "))
	quoter.IsReserved = AlwaysNoReserve
	assert.EqualValues(t, "f1, f2, f3", quoter.Join(cols, ", "))
}
func TestStrings(t *testing.T) {
	cols := []string{"f1", "f2", "t3.f3", "t4.*"}
	quoter := Quoter{'[', ']', AlwaysReserve}
	quotedCols := quoter.Strings(cols)
	assert.EqualValues(t, []string{"[f1]", "[f2]", "[t3].[f3]", "[t4].*"}, quotedCols)
}
func TestTrim(t *testing.T) {
	var kases = map[string]string{
		"[table_name]":          "table_name",
		"[schema].[table_name]": "schema.table_name",
	}
	for src, dst := range kases {
		assert.EqualValues(t, src, CommonQuoter.Trim(src))
		assert.EqualValues(t, dst, Quoter{'[', ']', AlwaysReserve}.Trim(src))
	}
}
func TestReplace(t *testing.T) {
	q := Quoter{'[', ']', AlwaysReserve}
	var kases = []struct {
		source   string
		expected string
	}{
		{
			"SELECT `COLUMN_NAME` FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ? AND `COLUMN_NAME` = ?",
			"SELECT [COLUMN_NAME] FROM [INFORMATION_SCHEMA].[COLUMNS] WHERE [TABLE_SCHEMA] = ? AND [TABLE_NAME] = ? AND [COLUMN_NAME] = ?",
		},
		{
			"SELECT 'abc```test```''', `a` FROM b",
			"SELECT 'abc```test```''', [a] FROM b",
		},
		{
			"UPDATE table SET `a` = ~ `a`, `b`='abc`'",
			"UPDATE table SET [a] = ~ [a], [b]='abc`'",
		},
		{
			"INSERT INTO `insert_where` (`height`,`name`,`repo_id`,`width`,`index`) SELECT $1,$2,$3,$4,coalesce(MAX(`index`),0)+1 FROM `insert_where` WHERE (`repo_id`=$5)",
			"INSERT INTO [insert_where] ([height],[name],[repo_id],[width],[index]) SELECT $1,$2,$3,$4,coalesce(MAX([index]),0)+1 FROM [insert_where] WHERE ([repo_id]=$5)",
		},
	}
	for _, kase := range kases {
		t.Run(kase.source, func(t *testing.T) {
			assert.EqualValues(t, kase.expected, q.Replace(kase.source))
		})
	}
}
