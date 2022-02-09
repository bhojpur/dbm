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
	"reflect"
	"testing"
)

func TestParseOracleConnStr(t *testing.T) {
	tests := []struct {
		in       string
		expected string
		valid    bool
	}{
		{"user/pass@tcp(server:1521)/db", "db", true},
		{"user/pass@server:1521/db", "db", true},
		{"user/pass@server:1521", "", true},
		{"user/pass@", "", false},
		{"user/pass", "", false},
		{"", "", false},
	}
	driver := QueryDriver("oci8")
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			driver := driver
			uri, err := driver.Parse("oci8", test.in)
			if err != nil && test.valid {
				t.Errorf("%q got unexpected error: %s", test.in, err)
			} else if err == nil && !reflect.DeepEqual(test.expected, uri.DBName) {
				t.Errorf("%q got: %#v want: %#v", test.in, uri.DBName, test.expected)
			}
		})
	}
}
