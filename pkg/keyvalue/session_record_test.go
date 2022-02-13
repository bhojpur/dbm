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
	"bytes"
	"testing"
)

func decodeEncode(v *sessionRecord) (res bool, err error) {
	b := new(bytes.Buffer)
	err = v.encode(b)
	if err != nil {
		return
	}
	v2 := &sessionRecord{}
	err = v.decode(b)
	if err != nil {
		return
	}
	b2 := new(bytes.Buffer)
	err = v2.encode(b2)
	if err != nil {
		return
	}
	return bytes.Equal(b.Bytes(), b2.Bytes()), nil
}

func TestSessionRecord_EncodeDecode(t *testing.T) {
	big := int64(1) << 50
	v := &sessionRecord{}
	i := int64(0)
	test := func() {
		res, err := decodeEncode(v)
		if err != nil {
			t.Fatalf("error when testing encode/decode sessionRecord: %v", err)
		}
		if !res {
			t.Error("encode/decode test failed at iteration:", i)
		}
	}

	for ; i < 4; i++ {
		test()
		v.addTable(3, big+300+i, big+400+i,
			makeInternalKey(nil, []byte("foo"), uint64(big+500+1), keyTypeVal),
			makeInternalKey(nil, []byte("zoo"), uint64(big+600+1), keyTypeDel))
		v.delTable(4, big+700+i)
		v.addCompPtr(int(i), makeInternalKey(nil, []byte("x"), uint64(big+900+1), keyTypeVal))
	}

	v.setComparer("foo")
	v.setJournalNum(big + 100)
	v.setPrevJournalNum(big + 99)
	v.setNextFileNum(big + 200)
	v.setSeqNum(uint64(big + 1000))
	test()
}
