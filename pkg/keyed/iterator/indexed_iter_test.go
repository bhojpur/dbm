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
	"sort"

	. "github.com/onsi/ginkgo"

	"github.com/bhojpur/dbm/pkg/keyed/comparer"
	. "github.com/bhojpur/dbm/pkg/keyed/iterator"
	"github.com/bhojpur/dbm/pkg/keyed/testutil"
)

type keyValue struct {
	key []byte
	testutil.KeyValue
}

type keyValueIndex []keyValue

func (x keyValueIndex) Search(key []byte) int {
	return sort.Search(x.Len(), func(i int) bool {
		return comparer.DefaultComparer.Compare(x[i].key, key) >= 0
	})
}

func (x keyValueIndex) Len() int                        { return len(x) }
func (x keyValueIndex) Index(i int) (key, value []byte) { return x[i].key, nil }
func (x keyValueIndex) Get(i int) Iterator              { return NewArrayIterator(x[i]) }

var _ = testutil.Defer(func() {
	Describe("Indexed iterator", func() {
		Test := func(n ...int) func() {
			if len(n) == 0 {
				rnd := testutil.NewRand()
				n = make([]int, rnd.Intn(17)+3)
				for i := range n {
					n[i] = rnd.Intn(19) + 1
				}
			}

			return func() {
				It("Should iterates and seeks correctly", func(done Done) {
					// Build key/value.
					index := make(keyValueIndex, len(n))
					sum := 0
					for _, x := range n {
						sum += x
					}
					kv := testutil.KeyValue_Generate(nil, sum, 1, 1, 10, 4, 4)
					for i, j := 0, 0; i < len(n); i++ {
						for x := n[i]; x > 0; x-- {
							key, value := kv.Index(j)
							index[i].key = key
							index[i].Put(key, value)
							j++
						}
					}

					// Test the iterator.
					t := testutil.IteratorTesting{
						KeyValue: kv.Clone(),
						Iter:     NewIndexedIterator(NewArrayIndexer(index), true),
					}
					testutil.DoIteratorTesting(&t)
					done <- true
				}, 15.0)
			}
		}

		Describe("with 100 keys", Test(100))
		Describe("with 50-50 keys", Test(50, 50))
		Describe("with 50-1 keys", Test(50, 1))
		Describe("with 50-1-50 keys", Test(50, 1, 50))
		Describe("with 1-50 keys", Test(1, 50))
		Describe("with random N-keys", Test())
	})
})
