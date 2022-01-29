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
	"database/sql/driver"
	"fmt"
	"time"
)

type NullTime time.Time

var (
	_ driver.Valuer = NullTime{}
)

func (ns *NullTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return convertTime(ns, value)
}

// Value implements the driver Valuer interface.
func (ns NullTime) Value() (driver.Value, error) {
	if (time.Time)(ns).IsZero() {
		return nil, nil
	}
	return (time.Time)(ns).Format("2006-01-02 15:04:05"), nil
}
func convertTime(dest *NullTime, src interface{}) error {
	// Common cases, without reflect.
	switch s := src.(type) {
	case string:
		t, err := time.Parse("2006-01-02 15:04:05", s)
		if err != nil {
			return err
		}
		*dest = NullTime(t)
		return nil
	case []uint8:
		t, err := time.Parse("2006-01-02 15:04:05", string(s))
		if err != nil {
			return err
		}
		*dest = NullTime(t)
		return nil
	case time.Time:
		*dest = NullTime(s)
		return nil
	case nil:
	default:
		return fmt.Errorf("unsupported driver -> Scan pair: %T -> %T", src, dest)
	}
	return nil
}

type EmptyScanner struct {
}

func (EmptyScanner) Scan(src interface{}) error {
	return nil
}
