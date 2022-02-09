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
	"github.com/bhojpur/dbm/pkg/keyed/util"
)

// BasicArray is the interface that wraps basic Len and Search method.
type BasicArray interface {
	// Len returns length of the array.
	Len() int

	// Search finds smallest index that point to a key that is greater
	// than or equal to the given key.
	Search(key []byte) int
}

// Array is the interface that wraps BasicArray and basic Index method.
type Array interface {
	BasicArray

	// Index returns key/value pair with index of i.
	Index(i int) (key, value []byte)
}

// Array is the interface that wraps BasicArray and basic Get method.
type ArrayIndexer interface {
	BasicArray

	// Get returns a new data iterator with index of i.
	Get(i int) Iterator
}

type basicArrayIterator struct {
	util.BasicReleaser
	array BasicArray
	pos   int
	err   error
}

func (i *basicArrayIterator) Valid() bool {
	return i.pos >= 0 && i.pos < i.array.Len() && !i.Released()
}

func (i *basicArrayIterator) First() bool {
	if i.Released() {
		i.err = ErrIterReleased
		return false
	}

	if i.array.Len() == 0 {
		i.pos = -1
		return false
	}
	i.pos = 0
	return true
}

func (i *basicArrayIterator) Last() bool {
	if i.Released() {
		i.err = ErrIterReleased
		return false
	}

	n := i.array.Len()
	if n == 0 {
		i.pos = 0
		return false
	}
	i.pos = n - 1
	return true
}

func (i *basicArrayIterator) Seek(key []byte) bool {
	if i.Released() {
		i.err = ErrIterReleased
		return false
	}

	n := i.array.Len()
	if n == 0 {
		i.pos = 0
		return false
	}
	i.pos = i.array.Search(key)
	if i.pos >= n {
		return false
	}
	return true
}

func (i *basicArrayIterator) Next() bool {
	if i.Released() {
		i.err = ErrIterReleased
		return false
	}

	i.pos++
	if n := i.array.Len(); i.pos >= n {
		i.pos = n
		return false
	}
	return true
}

func (i *basicArrayIterator) Prev() bool {
	if i.Released() {
		i.err = ErrIterReleased
		return false
	}

	i.pos--
	if i.pos < 0 {
		i.pos = -1
		return false
	}
	return true
}

func (i *basicArrayIterator) Error() error { return i.err }

type arrayIterator struct {
	basicArrayIterator
	array      Array
	pos        int
	key, value []byte
}

func (i *arrayIterator) updateKV() {
	if i.pos == i.basicArrayIterator.pos {
		return
	}
	i.pos = i.basicArrayIterator.pos
	if i.Valid() {
		i.key, i.value = i.array.Index(i.pos)
	} else {
		i.key = nil
		i.value = nil
	}
}

func (i *arrayIterator) Key() []byte {
	i.updateKV()
	return i.key
}

func (i *arrayIterator) Value() []byte {
	i.updateKV()
	return i.value
}

type arrayIteratorIndexer struct {
	basicArrayIterator
	array ArrayIndexer
}

func (i *arrayIteratorIndexer) Get() Iterator {
	if i.Valid() {
		return i.array.Get(i.basicArrayIterator.pos)
	}
	return nil
}

// NewArrayIterator returns an iterator from the given array.
func NewArrayIterator(array Array) Iterator {
	return &arrayIterator{
		basicArrayIterator: basicArrayIterator{array: array, pos: -1},
		array:              array,
		pos:                -1,
	}
}

// NewArrayIndexer returns an index iterator from the given array.
func NewArrayIndexer(array ArrayIndexer) IteratorIndexer {
	return &arrayIteratorIndexer{
		basicArrayIterator: basicArrayIterator{array: array, pos: -1},
		array:              array,
	}
}
