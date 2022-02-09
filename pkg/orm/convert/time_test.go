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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestString2Time(t *testing.T) {
	expectedLoc, err := time.LoadLocation("Asia/Kolkata")
	assert.NoError(t, err)
	var kases = map[string]time.Time{
		"2021-08-10":                time.Date(2021, 8, 10, 8, 0, 0, 0, expectedLoc),
		"2021-06-06T22:58:20+08:00": time.Date(2021, 6, 6, 22, 58, 20, 0, expectedLoc),
		"2021-07-11 10:44:00":       time.Date(2021, 7, 11, 18, 44, 0, 0, expectedLoc),
		"2021-08-10T10:33:04Z":      time.Date(2021, 8, 10, 18, 33, 04, 0, expectedLoc),
	}
	for layout, tm := range kases {
		t.Run(layout, func(t *testing.T) {
			target, err := String2Time(layout, time.UTC, expectedLoc)
			assert.NoError(t, err)
			assert.EqualValues(t, tm, *target)
		})
	}
}
