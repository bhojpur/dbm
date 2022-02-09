package convert

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
	"database/sql"
	"fmt"
	"time"
)

// Interface2Interface converts interface of pointer as interface of value
func Interface2Interface(userLocation *time.Location, v interface{}) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	switch vv := v.(type) {
	case *int64:
		return *vv, nil
	case *int8:
		return *vv, nil
	case *sql.NullString:
		return vv.String, nil
	case *sql.RawBytes:
		if len([]byte(*vv)) > 0 {
			return []byte(*vv), nil
		}
		return nil, nil
	case *sql.NullInt32:
		return vv.Int32, nil
	case *sql.NullInt64:
		return vv.Int64, nil
	case *sql.NullFloat64:
		return vv.Float64, nil
	case *sql.NullBool:
		if vv.Valid {
			return vv.Bool, nil
		}
		return nil, nil
	case *sql.NullTime:
		if vv.Valid {
			return vv.Time.In(userLocation).Format("2006-01-02 15:04:05"), nil
		}
		return "", nil
	default:
		return "", fmt.Errorf("convert assign string unsupported type: %#v", vv)
	}
}
