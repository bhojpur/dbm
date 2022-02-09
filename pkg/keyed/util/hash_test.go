package util

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

var hashTests = []struct {
	data []byte
	seed uint32
	hash uint32
}{
	{nil, 0xbc9f1d34, 0xbc9f1d34},
	{[]byte{0x62}, 0xbc9f1d34, 0xef1345c4},
	{[]byte{0xc3, 0x97}, 0xbc9f1d34, 0x5b663814},
	{[]byte{0xe2, 0x99, 0xa5}, 0xbc9f1d34, 0x323c078f},
	{[]byte{0xe1, 0x80, 0xb9, 0x32}, 0xbc9f1d34, 0xed21633a},
	{[]byte{
		0x01, 0xc0, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x14, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x04, 0x00,
		0x00, 0x00, 0x00, 0x14,
		0x00, 0x00, 0x00, 0x18,
		0x28, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}, 0x12345678, 0xf333dabb},
}

func TestHash(t *testing.T) {
	for i, x := range hashTests {
		h := Hash(x.data, x.seed)
		if h != x.hash {
			t.Fatalf("test-%d: invalid hash, %#x vs %#x", i, h, x.hash)
		}
	}
}
