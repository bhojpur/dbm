package iterator_test

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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bhojpur/dbm/pkg/keyed/comparer"
	. "github.com/bhojpur/dbm/pkg/keyed/iterator"
	"github.com/bhojpur/dbm/pkg/keyed/testutil"
)

var _ = testutil.Defer(func() {
	Describe("Merged iterator", func() {
		Test := func(filled int, empty int) func() {
			return func() {
				It("Should iterates and seeks correctly", func(done Done) {
					rnd := testutil.NewRand()

					// Build key/value.
					filledKV := make([]testutil.KeyValue, filled)
					kv := testutil.KeyValue_Generate(nil, 100, 1, 1, 10, 4, 4)
					kv.Iterate(func(i int, key, value []byte) {
						filledKV[rnd.Intn(filled)].Put(key, value)
					})

					// Create itearators.
					iters := make([]Iterator, filled+empty)
					for i := range iters {
						if empty == 0 || (rnd.Int()%2 == 0 && filled > 0) {
							filled--
							Expect(filledKV[filled].Len()).ShouldNot(BeZero())
							iters[i] = NewArrayIterator(filledKV[filled])
						} else {
							empty--
							iters[i] = NewEmptyIterator(nil)
						}
					}

					// Test the iterator.
					t := testutil.IteratorTesting{
						KeyValue: kv.Clone(),
						Iter:     NewMergedIterator(iters, comparer.DefaultComparer, true),
					}
					testutil.DoIteratorTesting(&t)
					done <- true
				}, 15.0)
			}
		}

		Describe("with three, all filled iterators", Test(3, 0))
		Describe("with one filled, one empty iterators", Test(1, 1))
		Describe("with one filled, two empty iterators", Test(1, 2))
	})
})
