package iterator

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
	"github.com/bhojpur/dbm/pkg/keyvalue/errors"
	"github.com/bhojpur/dbm/pkg/keyvalue/util"
)

// IteratorIndexer is the interface that wraps CommonIterator and basic Get
// method. IteratorIndexer provides index for indexed iterator.
type IteratorIndexer interface {
	CommonIterator

	// Get returns a new data iterator for the current position, or nil if
	// done.
	Get() Iterator
}

type indexedIterator struct {
	util.BasicReleaser
	index  IteratorIndexer
	strict bool

	data   Iterator
	err    error
	errf   func(err error)
	closed bool
}

func (i *indexedIterator) setData() {
	if i.data != nil {
		i.data.Release()
	}
	i.data = i.index.Get()
}

func (i *indexedIterator) clearData() {
	if i.data != nil {
		i.data.Release()
	}
	i.data = nil
}

func (i *indexedIterator) indexErr() {
	if err := i.index.Error(); err != nil {
		if i.errf != nil {
			i.errf(err)
		}
		i.err = err
	}
}

func (i *indexedIterator) dataErr() bool {
	if err := i.data.Error(); err != nil {
		if i.errf != nil {
			i.errf(err)
		}
		if i.strict || !errors.IsCorrupted(err) {
			i.err = err
			return true
		}
	}
	return false
}

func (i *indexedIterator) Valid() bool {
	return i.data != nil && i.data.Valid()
}

func (i *indexedIterator) First() bool {
	if i.err != nil {
		return false
	} else if i.Released() {
		i.err = ErrIterReleased
		return false
	}

	if !i.index.First() {
		i.indexErr()
		i.clearData()
		return false
	}
	i.setData()
	return i.Next()
}

func (i *indexedIterator) Last() bool {
	if i.err != nil {
		return false
	} else if i.Released() {
		i.err = ErrIterReleased
		return false
	}

	if !i.index.Last() {
		i.indexErr()
		i.clearData()
		return false
	}
	i.setData()
	if !i.data.Last() {
		if i.dataErr() {
			return false
		}
		i.clearData()
		return i.Prev()
	}
	return true
}

func (i *indexedIterator) Seek(key []byte) bool {
	if i.err != nil {
		return false
	} else if i.Released() {
		i.err = ErrIterReleased
		return false
	}

	if !i.index.Seek(key) {
		i.indexErr()
		i.clearData()
		return false
	}
	i.setData()
	if !i.data.Seek(key) {
		if i.dataErr() {
			return false
		}
		i.clearData()
		return i.Next()
	}
	return true
}

func (i *indexedIterator) Next() bool {
	if i.err != nil {
		return false
	} else if i.Released() {
		i.err = ErrIterReleased
		return false
	}

	switch {
	case i.data != nil && !i.data.Next():
		if i.dataErr() {
			return false
		}
		i.clearData()
		fallthrough
	case i.data == nil:
		if !i.index.Next() {
			i.indexErr()
			return false
		}
		i.setData()
		return i.Next()
	}
	return true
}

func (i *indexedIterator) Prev() bool {
	if i.err != nil {
		return false
	} else if i.Released() {
		i.err = ErrIterReleased
		return false
	}

	switch {
	case i.data != nil && !i.data.Prev():
		if i.dataErr() {
			return false
		}
		i.clearData()
		fallthrough
	case i.data == nil:
		if !i.index.Prev() {
			i.indexErr()
			return false
		}
		i.setData()
		if !i.data.Last() {
			if i.dataErr() {
				return false
			}
			i.clearData()
			return i.Prev()
		}
	}
	return true
}

func (i *indexedIterator) Key() []byte {
	if i.data == nil {
		return nil
	}
	return i.data.Key()
}

func (i *indexedIterator) Value() []byte {
	if i.data == nil {
		return nil
	}
	return i.data.Value()
}

func (i *indexedIterator) Release() {
	i.clearData()
	i.index.Release()
	i.BasicReleaser.Release()
}

func (i *indexedIterator) Error() error {
	if i.err != nil {
		return i.err
	}
	if err := i.index.Error(); err != nil {
		return err
	}
	return nil
}

func (i *indexedIterator) SetErrorCallback(f func(err error)) {
	i.errf = f
}

// NewIndexedIterator returns an 'indexed iterator'. An index is iterator
// that returns another iterator, a 'data iterator'. A 'data iterator' is the
// iterator that contains actual key/value pairs.
//
// If strict is true the any 'corruption errors' (i.e errors.IsCorrupted(err) == true)
// won't be ignored and will halt 'indexed iterator', otherwise the iterator will
// continue to the next 'data iterator'. Corruption on 'index iterator' will not be
// ignored and will halt the iterator.
func NewIndexedIterator(index IteratorIndexer, strict bool) Iterator {
	return &indexedIterator{index: index, strict: strict}
}
