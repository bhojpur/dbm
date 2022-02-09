package filter

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
	"encoding/binary"
	"testing"

	"github.com/bhojpur/dbm/pkg/keyed/util"
)

type harness struct {
	t *testing.T

	bloom     Filter
	generator FilterGenerator
	filter    []byte
}

func newHarness(t *testing.T) *harness {
	bloom := NewBloomFilter(10)
	return &harness{
		t:         t,
		bloom:     bloom,
		generator: bloom.NewGenerator(),
	}
}

func (h *harness) add(key []byte) {
	h.generator.Add(key)
}

func (h *harness) addNum(key uint32) {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], key)
	h.add(b[:])
}

func (h *harness) build() {
	b := &util.Buffer{}
	h.generator.Generate(b)
	h.filter = b.Bytes()
}

func (h *harness) reset() {
	h.filter = nil
}

func (h *harness) filterLen() int {
	return len(h.filter)
}

func (h *harness) assert(key []byte, want, silent bool) bool {
	got := h.bloom.Contains(h.filter, key)
	if !silent && got != want {
		h.t.Errorf("assert on '%v' failed got '%v', want '%v'", key, got, want)
	}
	return got
}

func (h *harness) assertNum(key uint32, want, silent bool) bool {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], key)
	return h.assert(b[:], want, silent)
}

func TestBloomFilter_Empty(t *testing.T) {
	h := newHarness(t)
	h.build()
	h.assert([]byte("hello"), false, false)
	h.assert([]byte("world"), false, false)
}

func TestBloomFilter_Small(t *testing.T) {
	h := newHarness(t)
	h.add([]byte("hello"))
	h.add([]byte("world"))
	h.build()
	h.assert([]byte("hello"), true, false)
	h.assert([]byte("world"), true, false)
	h.assert([]byte("x"), false, false)
	h.assert([]byte("foo"), false, false)
}

func nextN(n int) int {
	switch {
	case n < 10:
		n += 1
	case n < 100:
		n += 10
	case n < 1000:
		n += 100
	default:
		n += 1000
	}
	return n
}

func TestBloomFilter_VaryingLengths(t *testing.T) {
	h := newHarness(t)
	var mediocre, good int
	for n := 1; n < 10000; n = nextN(n) {
		h.reset()
		for i := 0; i < n; i++ {
			h.addNum(uint32(i))
		}
		h.build()

		got := h.filterLen()
		want := (n * 10 / 8) + 40
		if got > want {
			t.Errorf("filter len test failed, '%d' > '%d'", got, want)
		}

		for i := 0; i < n; i++ {
			h.assertNum(uint32(i), true, false)
		}

		var rate float32
		for i := 0; i < 10000; i++ {
			if h.assertNum(uint32(i+1000000000), true, true) {
				rate++
			}
		}
		rate /= 10000
		if rate > 0.02 {
			t.Errorf("false positive rate is more than 2%%, got %v, at len %d", rate, n)
		}
		if rate > 0.0125 {
			mediocre++
		} else {
			good++
		}
	}
	t.Logf("false positive rate: %d good, %d mediocre", good, mediocre)
	if mediocre > good/5 {
		t.Error("mediocre false positive rate is more than expected")
	}
}
