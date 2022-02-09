package statement

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

	"github.com/bhojpur/dbm/pkg/orm/internal/utils"
	"github.com/bhojpur/dbm/pkg/orm/schema"
)

// ConvertIDSQL converts SQL with id
func (statement *Statement) ConvertIDSQL(sqlStr string) string {
	if statement.RefTable != nil {
		cols := statement.RefTable.PKColumns()
		if len(cols) == 0 {
			return ""
		}
		colstrs := statement.joinColumns(cols, false)
		sqls := utils.SplitNNoCase(sqlStr, " from ", 2)
		if len(sqls) != 2 {
			return ""
		}
		var top string
		pLimitN := statement.LimitN
		if pLimitN != nil && statement.dialect.URI().DBType == schema.MSSQL {
			top = fmt.Sprintf("TOP %d ", *pLimitN)
		}
		newsql := fmt.Sprintf("SELECT %s%s FROM %v", top, colstrs, sqls[1])
		return newsql
	}
	return ""
}

// ConvertUpdateSQL converts update SQL
func (statement *Statement) ConvertUpdateSQL(sqlStr string) (string, string) {
	if statement.RefTable == nil || len(statement.RefTable.PrimaryKeys) != 1 {
		return "", ""
	}
	colstrs := statement.joinColumns(statement.RefTable.PKColumns(), true)
	sqls := utils.SplitNNoCase(sqlStr, "where", 2)
	if len(sqls) != 2 {
		if len(sqls) == 1 {
			return sqls[0], fmt.Sprintf("SELECT %v FROM %v",
				colstrs, statement.quote(statement.TableName()))
		}
		return "", ""
	}
	var whereStr = sqls[1]
	// TODO: for postgres only, if any other database?
	var paraStr string
	if statement.dialect.URI().DBType == schema.POSTGRES {
		paraStr = "$"
	} else if statement.dialect.URI().DBType == schema.MSSQL {
		paraStr = ":"
	}
	if paraStr != "" {
		if strings.Contains(sqls[1], paraStr) {
			dollers := strings.Split(sqls[1], paraStr)
			whereStr = dollers[0]
			for i, c := range dollers[1:] {
				ccs := strings.SplitN(c, " ", 2)
				whereStr += fmt.Sprintf(paraStr+"%v %v", i+1, ccs[1])
			}
		}
	}
	return sqls[0], fmt.Sprintf("SELECT %v FROM %v WHERE %v",
		colstrs, statement.quote(statement.TableName()),
		whereStr)
}
