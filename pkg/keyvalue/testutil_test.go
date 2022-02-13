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
	. "github.com/onsi/gomega"

	"github.com/bhojpur/dbm/pkg/keyvalue/iterator"
	"github.com/bhojpur/dbm/pkg/keyvalue/opt"
	"github.com/bhojpur/dbm/pkg/keyvalue/testutil"
	"github.com/bhojpur/dbm/pkg/keyvalue/util"
)

type testingDB struct {
	*DB
	ro   *opt.ReadOptions
	wo   *opt.WriteOptions
	stor *testutil.Storage
}

func (t *testingDB) TestPut(key []byte, value []byte) error {
	return t.Put(key, value, t.wo)
}

func (t *testingDB) TestDelete(key []byte) error {
	return t.Delete(key, t.wo)
}

func (t *testingDB) TestGet(key []byte) (value []byte, err error) {
	return t.Get(key, t.ro)
}

func (t *testingDB) TestHas(key []byte) (ret bool, err error) {
	return t.Has(key, t.ro)
}

func (t *testingDB) TestNewIterator(slice *util.Range) iterator.Iterator {
	return t.NewIterator(slice, t.ro)
}

func (t *testingDB) TestClose() {
	err := t.Close()
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	err = t.stor.Close()
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
}

func newTestingDB(o *opt.Options, ro *opt.ReadOptions, wo *opt.WriteOptions) *testingDB {
	stor := testutil.NewStorage()
	db, err := Open(stor, o)
	// FIXME: This may be called from outside It, which may cause panic.
	Expect(err).NotTo(HaveOccurred())
	return &testingDB{
		DB:   db,
		ro:   ro,
		wo:   wo,
		stor: stor,
	}
}

type testingTransaction struct {
	*Transaction
	ro *opt.ReadOptions
	wo *opt.WriteOptions
}

func (t *testingTransaction) TestPut(key []byte, value []byte) error {
	return t.Put(key, value, t.wo)
}

func (t *testingTransaction) TestDelete(key []byte) error {
	return t.Delete(key, t.wo)
}

func (t *testingTransaction) TestGet(key []byte) (value []byte, err error) {
	return t.Get(key, t.ro)
}

func (t *testingTransaction) TestHas(key []byte) (ret bool, err error) {
	return t.Has(key, t.ro)
}

func (t *testingTransaction) TestNewIterator(slice *util.Range) iterator.Iterator {
	return t.NewIterator(slice, t.ro)
}

func (t *testingTransaction) TestClose() {}
