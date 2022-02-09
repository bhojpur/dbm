package table

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
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bhojpur/dbm/pkg/keyed/comparer"
	"github.com/bhojpur/dbm/pkg/keyed/iterator"
	"github.com/bhojpur/dbm/pkg/keyed/testutil"
	"github.com/bhojpur/dbm/pkg/keyed/util"
)

type blockTesting struct {
	tr *Reader
	b  *block
}

func (t *blockTesting) TestNewIterator(slice *util.Range) iterator.Iterator {
	return t.tr.newBlockIter(t.b, nil, slice, false)
}

var _ = testutil.Defer(func() {
	Describe("Block", func() {
		Build := func(kv *testutil.KeyValue, restartInterval int) *blockTesting {
			// Building the block.
			bw := &blockWriter{
				restartInterval: restartInterval,
				scratch:         make([]byte, 30),
			}
			kv.Iterate(func(i int, key, value []byte) {
				bw.append(key, value)
			})
			bw.finish()

			// Opening the block.
			data := bw.buf.Bytes()
			restartsLen := int(binary.LittleEndian.Uint32(data[len(data)-4:]))
			return &blockTesting{
				tr: &Reader{cmp: comparer.DefaultComparer},
				b: &block{
					data:           data,
					restartsLen:    restartsLen,
					restartsOffset: len(data) - (restartsLen+1)*4,
				},
			}
		}

		Describe("read test", func() {
			for restartInterval := 1; restartInterval <= 5; restartInterval++ {
				Describe(fmt.Sprintf("with restart interval of %d", restartInterval), func() {
					kv := &testutil.KeyValue{}
					Text := func() string {
						return fmt.Sprintf("and %d keys", kv.Len())
					}

					Test := func() {
						// Make block.
						br := Build(kv, restartInterval)
						// Do testing.
						testutil.KeyValueTesting(nil, kv.Clone(), br, nil, nil)
					}

					Describe(Text(), Test)

					kv.PutString("", "empty")
					Describe(Text(), Test)

					kv.PutString("a1", "foo")
					Describe(Text(), Test)

					kv.PutString("a2", "v")
					Describe(Text(), Test)

					kv.PutString("a3qqwrkks", "hello")
					Describe(Text(), Test)

					kv.PutString("a4", "bar")
					Describe(Text(), Test)

					kv.PutString("a5111111", "v5")
					kv.PutString("a6", "")
					kv.PutString("a7", "v7")
					kv.PutString("a8", "vvvvvvvvvvvvvvvvvvvvvv8")
					kv.PutString("b", "v9")
					kv.PutString("c9", "v9")
					kv.PutString("c91", "v9")
					kv.PutString("d0", "v9")
					Describe(Text(), Test)
				})
			}
		})

		Describe("out-of-bound slice test", func() {
			kv := &testutil.KeyValue{}
			kv.PutString("k1", "v1")
			kv.PutString("k2", "v2")
			kv.PutString("k3abcdefgg", "v3")
			kv.PutString("k4", "v4")
			kv.PutString("k5", "v5")
			for restartInterval := 1; restartInterval <= 5; restartInterval++ {
				Describe(fmt.Sprintf("with restart interval of %d", restartInterval), func() {
					// Make block.
					bt := Build(kv, restartInterval)

					Test := func(r *util.Range) func(done Done) {
						return func(done Done) {
							iter := bt.TestNewIterator(r)
							Expect(iter.Error()).ShouldNot(HaveOccurred())

							t := testutil.IteratorTesting{
								KeyValue: kv.Clone(),
								Iter:     iter,
							}

							testutil.DoIteratorTesting(&t)
							iter.Release()
							done <- true
						}
					}

					It("Should do iterations and seeks correctly #0",
						Test(&util.Range{Start: []byte("k0"), Limit: []byte("k6")}), 2.0)

					It("Should do iterations and seeks correctly #1",
						Test(&util.Range{Start: []byte(""), Limit: []byte("zzzzzzz")}), 2.0)
				})
			}
		})
	})
})
