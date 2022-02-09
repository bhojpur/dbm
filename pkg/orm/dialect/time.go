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
	"strings"
	"time"

	schemasvr "github.com/bhojpur/dbm/pkg/orm/schema"
)

// FormatColumnTime format column time
func FormatColumnTime(dialect Dialect, dbLocation *time.Location, col *schemasvr.Column, t time.Time) (interface{}, error) {
	if t.IsZero() {
		if col.Nullable {
			return nil, nil
		}
		if col.SQLType.IsNumeric() {
			return 0, nil
		}
	}
	var tmZone = dbLocation
	if col.TimeZone != nil {
		tmZone = col.TimeZone
	}
	t = t.In(tmZone)
	switch col.SQLType.Name {
	case schemasvr.Date:
		return t.Format("2006-01-02"), nil
	case schemasvr.Time:
		var layout = "15:04:05"
		if col.Length > 0 {
			layout += "." + strings.Repeat("0", col.Length)
		}
		return t.Format(layout), nil
	case schemasvr.DateTime, schemasvr.TimeStamp:
		var layout = "2006-01-02 15:04:05"
		if col.Length > 0 {
			layout += "." + strings.Repeat("0", col.Length)
		}
		return t.Format(layout), nil
	case schemasvr.Varchar:
		return t.Format("2006-01-02 15:04:05"), nil
	case schemasvr.TimeStampz:
		if dialect.URI().DBType == schemasvr.MSSQL {
			return t.Format("2006-01-02T15:04:05.9999999Z07:00"), nil
		} else {
			return t.Format(time.RFC3339Nano), nil
		}
	case schemasvr.BigInt, schemasvr.Int:
		return t.Unix(), nil
	default:
		return t, nil
	}
}
