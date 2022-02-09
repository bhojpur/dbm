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

	"github.com/stretchr/testify/assert"
)

func TestParsePostgres(t *testing.T) {
	tests := []struct {
		in       string
		expected string
		valid    bool
	}{
		{"postgres://auser:password@localhost:5432/db?sslmode=disable", "db", true},
		{"postgresql://auser:password@localhost:5432/db?sslmode=disable", "db", true},
		{"postg://auser:password@localhost:5432/db?sslmode=disable", "db", false},
		//{"postgres://auser:pass with space@localhost:5432/db?sslmode=disable", "db", true},
		//{"postgres:// auser : password@localhost:5432/db?sslmode=disable", "db", true},
		{"postgres://%20auser%20:pass%20with%20space@localhost:5432/db?sslmode=disable", "db", true},
		//{"postgres://auser:パスワード@localhost:5432/データベース?sslmode=disable", "データベース", true},
		{"dbname=db sslmode=disable", "db", true},
		{"user=auser password=password dbname=db sslmode=disable", "db", true},
		{"user=auser password='pass word' dbname=db sslmode=disable", "db", true},
		{"user=auser password='pass word' sslmode=disable dbname='db'", "db", true},
		{"user=auser password='pass word' sslmode='disable dbname=db'", "db", false},
		{"", "db", false},
		{"dbname=db =disable", "db", false},
	}
	driver := QueryDriver("postgres")
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			uri, err := driver.Parse("postgres", test.in)
			if err != nil && test.valid {
				t.Errorf("%q got unexpected error: %s", test.in, err)
			} else if err == nil && !reflect.DeepEqual(test.expected, uri.DBName) {
				t.Errorf("%q got: %#v want: %#v", test.in, uri.DBName, test.expected)
			}
		})
	}
}
func TestParsePgx(t *testing.T) {
	tests := []struct {
		in       string
		expected string
		valid    bool
	}{
		{"postgres://auser:password@localhost:5432/db?sslmode=disable", "db", true},
		{"postgresql://auser:password@localhost:5432/db?sslmode=disable", "db", true},
		{"postg://auser:password@localhost:5432/db?sslmode=disable", "db", false},
		//{"postgres://auser:pass with space@localhost:5432/db?sslmode=disable", "db", true},
		//{"postgres:// auser : password@localhost:5432/db?sslmode=disable", "db", true},
		{"postgres://%20auser%20:pass%20with%20space@localhost:5432/db?sslmode=disable", "db", true},
		//{"postgres://auser:パスワード@localhost:5432/データベース?sslmode=disable", "データベース", true},
		{"dbname=db sslmode=disable", "db", true},
		{"user=auser password=password dbname=db sslmode=disable", "db", true},
		{"", "db", false},
		{"dbname=db =disable", "db", false},
	}
	driver := QueryDriver("pgx")
	for _, test := range tests {
		uri, err := driver.Parse("pgx", test.in)
		if err != nil && test.valid {
			t.Errorf("%q got unexpected error: %s", test.in, err)
		} else if err == nil && !reflect.DeepEqual(test.expected, uri.DBName) {
			t.Errorf("%q got: %#v want: %#v", test.in, uri.DBName, test.expected)
		}
		// Register DriverConfig
		uri, err = driver.Parse("pgx", test.in)
		if err != nil && test.valid {
			t.Errorf("%q got unexpected error: %s", test.in, err)
		} else if err == nil && !reflect.DeepEqual(test.expected, uri.DBName) {
			t.Errorf("%q got: %#v want: %#v", test.in, uri.DBName, test.expected)
		}
	}
}
func TestGetIndexColName(t *testing.T) {
	t.Run("Index", func(t *testing.T) {
		s := "CREATE INDEX test2_mm_idx ON test2 (major);"
		colNames := getIndexColName(s)
		assert.Equal(t, []string{"major"}, colNames)
	})
	t.Run("Multicolumn indexes", func(t *testing.T) {
		s := "CREATE INDEX test2_mm_idx ON test2 (major, minor);"
		colNames := getIndexColName(s)
		assert.Equal(t, []string{"major", "minor"}, colNames)
	})
	t.Run("Indexes and ORDER BY", func(t *testing.T) {
		s := "CREATE INDEX test2_mm_idx ON test2 (major  NULLS FIRST, minor DESC NULLS LAST);"
		colNames := getIndexColName(s)
		assert.Equal(t, []string{"major", "minor"}, colNames)
	})
	t.Run("Combining Multiple Indexes", func(t *testing.T) {
		s := "CREATE INDEX test2_mm_cm_idx ON public.test2 USING btree (major, minor) WHERE ((major <> 5) AND (minor <> 6))"
		colNames := getIndexColName(s)
		assert.Equal(t, []string{"major", "minor"}, colNames)
	})
	t.Run("unique", func(t *testing.T) {
		s := "CREATE UNIQUE INDEX test2_mm_uidx ON test2 (major);"
		colNames := getIndexColName(s)
		assert.Equal(t, []string{"major"}, colNames)
	})
	t.Run("Indexes on Expressions", func(t *testing.T) {})
}
