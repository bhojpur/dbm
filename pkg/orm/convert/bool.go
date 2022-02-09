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
	"strconv"
)

// AsBool convert interface as bool
func AsBool(src interface{}) (bool, error) {
	switch v := src.(type) {
	case bool:
		return v, nil
	case *bool:
		return *v, nil
	case *sql.NullBool:
		return v.Bool, nil
	case int64:
		return v > 0, nil
	case int:
		return v > 0, nil
	case int8:
		return v > 0, nil
	case int16:
		return v > 0, nil
	case int32:
		return v > 0, nil
	case []byte:
		if len(v) == 0 {
			return false, nil
		}
		if v[0] == 0x00 {
			return false, nil
		} else if v[0] == 0x01 {
			return true, nil
		}
		return strconv.ParseBool(string(v))
	case string:
		return strconv.ParseBool(v)
	case *sql.NullInt64:
		return v.Int64 > 0, nil
	case *sql.NullInt32:
		return v.Int32 > 0, nil
	default:
		return false, fmt.Errorf("unknown type %T as bool", src)
	}
}
