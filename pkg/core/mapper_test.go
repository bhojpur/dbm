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
	"testing"
)

func TestGonicMapperFromObj(t *testing.T) {
	testCases := map[string]string{
		"HTTPLib":             "http_lib",
		"id":                  "id",
		"ID":                  "id",
		"IDa":                 "i_da",
		"iDa":                 "i_da",
		"IDAa":                "id_aa",
		"aID":                 "a_id",
		"aaID":                "aa_id",
		"aaaID":               "aaa_id",
		"MyREalFunkYLONgNAME": "my_r_eal_funk_ylo_ng_name",
	}
	for in, expected := range testCases {
		out := gonicCasedName(in)
		if out != expected {
			t.Errorf("Given %s, expected %s but got %s", in, expected, out)
		}
	}
}
func TestGonicMapperToObj(t *testing.T) {
	testCases := map[string]string{
		"http_lib":                  "HTTPLib",
		"id":                        "ID",
		"ida":                       "Ida",
		"id_aa":                     "IDAa",
		"aa_id":                     "AaID",
		"my_r_eal_funk_ylo_ng_name": "MyREalFunkYloNgName",
	}
	for in, expected := range testCases {
		out := LintGonicMapper.Table2Obj(in)
		if out != expected {
			t.Errorf("Given %s, expected %s but got %s", in, expected, out)
		}
	}
}
