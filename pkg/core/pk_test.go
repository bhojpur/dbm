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
	"reflect"
	"testing"
)

func TestPK(t *testing.T) {
	p := NewPK(1, 3, "string")
	str, err := p.ToString()
	if err != nil {
		t.Error(err)
	}
	t.Log(str)
	s := &PK{}
	err = s.FromString(str)
	if err != nil {
		t.Error(err)
	}
	t.Log(s)
	if len(*p) != len(*s) {
		t.Fatal("p", *p, "should be equal", *s)
	}
	for i, ori := range *p {
		if ori != (*s)[i] {
			t.Fatal("ori", ori, reflect.ValueOf(ori), "should be equal", (*s)[i], reflect.ValueOf((*s)[i]))
		}
	}
}
