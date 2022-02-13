package keyvalue

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
	"github.com/bhojpur/dbm/pkg/keyvalue/comparer"
)

type iComparer struct {
	ucmp comparer.Comparer
}

func (icmp *iComparer) uName() string {
	return icmp.ucmp.Name()
}

func (icmp *iComparer) uCompare(a, b []byte) int {
	return icmp.ucmp.Compare(a, b)
}

func (icmp *iComparer) uSeparator(dst, a, b []byte) []byte {
	return icmp.ucmp.Separator(dst, a, b)
}

func (icmp *iComparer) uSuccessor(dst, b []byte) []byte {
	return icmp.ucmp.Successor(dst, b)
}

func (icmp *iComparer) Name() string {
	return icmp.uName()
}

func (icmp *iComparer) Compare(a, b []byte) int {
	x := icmp.uCompare(internalKey(a).ukey(), internalKey(b).ukey())
	if x == 0 {
		if m, n := internalKey(a).num(), internalKey(b).num(); m > n {
			return -1
		} else if m < n {
			return 1
		}
	}
	return x
}

func (icmp *iComparer) Separator(dst, a, b []byte) []byte {
	ua, ub := internalKey(a).ukey(), internalKey(b).ukey()
	dst = icmp.uSeparator(dst, ua, ub)
	if dst != nil && len(dst) < len(ua) && icmp.uCompare(ua, dst) < 0 {
		// Append earliest possible number.
		return append(dst, keyMaxNumBytes...)
	}
	return nil
}

func (icmp *iComparer) Successor(dst, b []byte) []byte {
	ub := internalKey(b).ukey()
	dst = icmp.uSuccessor(dst, ub)
	if dst != nil && len(dst) < len(ub) && icmp.uCompare(ub, dst) < 0 {
		// Append earliest possible number.
		return append(dst, keyMaxNumBytes...)
	}
	return nil
}
