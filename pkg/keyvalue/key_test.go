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

	"github.com/bhojpur/dbm/pkg/keyvalue/comparer"
)

var defaultIComparer = &iComparer{comparer.DefaultComparer}

func ikey(key string, seq uint64, kt keyType) internalKey {
	return makeInternalKey(nil, []byte(key), uint64(seq), kt)
}

func shortSep(a, b []byte) []byte {
	dst := make([]byte, len(a))
	dst = defaultIComparer.Separator(dst[:0], a, b)
	if dst == nil {
		return a
	}
	return dst
}

func shortSuccessor(b []byte) []byte {
	dst := make([]byte, len(b))
	dst = defaultIComparer.Successor(dst[:0], b)
	if dst == nil {
		return b
	}
	return dst
}

func testSingleKey(t *testing.T, key string, seq uint64, kt keyType) {
	ik := ikey(key, seq, kt)

	if !bytes.Equal(ik.ukey(), []byte(key)) {
		t.Errorf("user key does not equal, got %v, want %v", string(ik.ukey()), key)
	}

	rseq, rt := ik.parseNum()
	if rseq != seq {
		t.Errorf("seq number does not equal, got %v, want %v", rseq, seq)
	}
	if rt != kt {
		t.Errorf("type does not equal, got %v, want %v", rt, kt)
	}

	if rukey, rseq, rt, kerr := parseInternalKey(ik); kerr == nil {
		if !bytes.Equal(rukey, []byte(key)) {
			t.Errorf("user key does not equal, got %v, want %v", string(ik.ukey()), key)
		}
		if rseq != seq {
			t.Errorf("seq number does not equal, got %v, want %v", rseq, seq)
		}
		if rt != kt {
			t.Errorf("type does not equal, got %v, want %v", rt, kt)
		}
	} else {
		t.Errorf("key error: %v", kerr)
	}
}

func TestInternalKey_EncodeDecode(t *testing.T) {
	keys := []string{"", "k", "hello", "longggggggggggggggggggggg"}
	seqs := []uint64{
		1, 2, 3,
		(1 << 8) - 1, 1 << 8, (1 << 8) + 1,
		(1 << 16) - 1, 1 << 16, (1 << 16) + 1,
		(1 << 32) - 1, 1 << 32, (1 << 32) + 1,
	}
	for _, key := range keys {
		for _, seq := range seqs {
			testSingleKey(t, key, seq, keyTypeVal)
			testSingleKey(t, "hello", 1, keyTypeDel)
		}
	}
}

func assertBytes(t *testing.T, want, got []byte) {
	if !bytes.Equal(got, want) {
		t.Errorf("assert failed, got %v, want %v", got, want)
	}
}

func TestInternalKeyShortSeparator(t *testing.T) {
	// When user keys are same
	assertBytes(t, ikey("foo", 100, keyTypeVal),
		shortSep(ikey("foo", 100, keyTypeVal),
			ikey("foo", 99, keyTypeVal)))
	assertBytes(t, ikey("foo", 100, keyTypeVal),
		shortSep(ikey("foo", 100, keyTypeVal),
			ikey("foo", 101, keyTypeVal)))
	assertBytes(t, ikey("foo", 100, keyTypeVal),
		shortSep(ikey("foo", 100, keyTypeVal),
			ikey("foo", 100, keyTypeVal)))
	assertBytes(t, ikey("foo", 100, keyTypeVal),
		shortSep(ikey("foo", 100, keyTypeVal),
			ikey("foo", 100, keyTypeDel)))

	// When user keys are misordered
	assertBytes(t, ikey("foo", 100, keyTypeVal),
		shortSep(ikey("foo", 100, keyTypeVal),
			ikey("bar", 99, keyTypeVal)))

	// When user keys are different, but correctly ordered
	assertBytes(t, ikey("g", uint64(keyMaxSeq), keyTypeSeek),
		shortSep(ikey("foo", 100, keyTypeVal),
			ikey("hello", 200, keyTypeVal)))

	// When start user key is prefix of limit user key
	assertBytes(t, ikey("foo", 100, keyTypeVal),
		shortSep(ikey("foo", 100, keyTypeVal),
			ikey("foobar", 200, keyTypeVal)))

	// When limit user key is prefix of start user key
	assertBytes(t, ikey("foobar", 100, keyTypeVal),
		shortSep(ikey("foobar", 100, keyTypeVal),
			ikey("foo", 200, keyTypeVal)))
}

func TestInternalKeyShortestSuccessor(t *testing.T) {
	assertBytes(t, ikey("g", uint64(keyMaxSeq), keyTypeSeek),
		shortSuccessor(ikey("foo", 100, keyTypeVal)))
	assertBytes(t, ikey("\xff\xff", 100, keyTypeVal),
		shortSuccessor(ikey("\xff\xff", 100, keyTypeVal)))
}
