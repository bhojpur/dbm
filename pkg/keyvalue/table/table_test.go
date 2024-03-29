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
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bhojpur/dbm/pkg/keyvalue/iterator"
	"github.com/bhojpur/dbm/pkg/keyvalue/opt"
	"github.com/bhojpur/dbm/pkg/keyvalue/storage"
	"github.com/bhojpur/dbm/pkg/keyvalue/testutil"
	"github.com/bhojpur/dbm/pkg/keyvalue/util"
)

type tableWrapper struct {
	*Reader
}

func (t tableWrapper) TestFind(key []byte) (rkey, rvalue []byte, err error) {
	return t.Reader.Find(key, false, nil)
}

func (t tableWrapper) TestGet(key []byte) (value []byte, err error) {
	return t.Reader.Get(key, nil)
}

func (t tableWrapper) TestNewIterator(slice *util.Range) iterator.Iterator {
	return t.Reader.NewIterator(slice, nil)
}

var _ = testutil.Defer(func() {
	Describe("Table", func() {
		Describe("approximate offset test", func() {
			var (
				buf = &bytes.Buffer{}
				o   = &opt.Options{
					BlockSize:   1024,
					Compression: opt.NoCompression,
				}
			)

			// Building the table.
			tw := NewWriter(buf, o, nil, 0)
			tw.Append([]byte("k01"), []byte("hello"))
			tw.Append([]byte("k02"), []byte("hello2"))
			tw.Append([]byte("k03"), bytes.Repeat([]byte{'x'}, 10000))
			tw.Append([]byte("k04"), bytes.Repeat([]byte{'x'}, 200000))
			tw.Append([]byte("k05"), bytes.Repeat([]byte{'x'}, 300000))
			tw.Append([]byte("k06"), []byte("hello3"))
			tw.Append([]byte("k07"), bytes.Repeat([]byte{'x'}, 100000))
			err := tw.Close()

			It("Should be able to approximate offset of a key correctly", func() {
				Expect(err).ShouldNot(HaveOccurred())

				tr, err := NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()), storage.FileDesc{}, nil, nil, o)
				Expect(err).ShouldNot(HaveOccurred())
				CheckOffset := func(key string, expect, threshold int) {
					offset, err := tr.OffsetOf([]byte(key))
					Expect(err).ShouldNot(HaveOccurred())
					Expect(offset).Should(BeNumerically("~", expect, threshold), "Offset of key %q", key)
				}

				CheckOffset("k0", 0, 0)
				CheckOffset("k01a", 0, 0)
				CheckOffset("k02", 0, 0)
				CheckOffset("k03", 0, 0)
				CheckOffset("k04", 10000, 1000)
				CheckOffset("k04a", 210000, 1000)
				CheckOffset("k05", 210000, 1000)
				CheckOffset("k06", 510000, 1000)
				CheckOffset("k07", 510000, 1000)
				CheckOffset("xyz", 610000, 2000)
			})
		})

		Describe("read test", func() {
			Build := func(kv testutil.KeyValue) testutil.DB {
				o := &opt.Options{
					BlockSize:            512,
					BlockRestartInterval: 3,
				}
				buf := &bytes.Buffer{}

				// Building the table.
				tw := NewWriter(buf, o, nil, 0)
				kv.Iterate(func(i int, key, value []byte) {
					tw.Append(key, value)
				})
				tw.Close()

				// Opening the table.
				tr, _ := NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()), storage.FileDesc{}, nil, nil, o)
				return tableWrapper{tr}
			}
			Test := func(kv *testutil.KeyValue, body func(r *Reader)) func() {
				return func() {
					db := Build(*kv)
					if body != nil {
						body(db.(tableWrapper).Reader)
					}
					testutil.KeyValueTesting(nil, *kv, db, nil, nil)
				}
			}

			testutil.AllKeyValueTesting(nil, Build, nil, nil)
			Describe("with one key per block", Test(testutil.KeyValue_Generate(nil, 9, 1, 1, 10, 512, 512), func(r *Reader) {
				It("should have correct blocks number", func() {
					indexBlock, err := r.readBlock(r.indexBH, true)
					Expect(err).To(BeNil())
					Expect(indexBlock.restartsLen).Should(Equal(9))
				})
			}))
		})
	})
})
