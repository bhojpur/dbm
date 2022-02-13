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

// It  provides interface and implementation to traverse over contents of
// a database.

import (
	"errors"

	"github.com/bhojpur/dbm/pkg/keyvalue/util"
)

var (
	ErrIterReleased = errors.New("keyvalue/iterator: iterator released")
)

// IteratorSeeker is the interface that wraps the 'seeks method'.
type IteratorSeeker interface {
	// First moves the iterator to the first key/value pair. If the iterator
	// only contains one key/value pair then First and Last would moves
	// to the same key/value pair.
	// It returns whether such pair exist.
	First() bool

	// Last moves the iterator to the last key/value pair. If the iterator
	// only contains one key/value pair then First and Last would moves
	// to the same key/value pair.
	// It returns whether such pair exist.
	Last() bool

	// Seek moves the iterator to the first key/value pair whose key is greater
	// than or equal to the given key.
	// It returns whether such pair exist.
	//
	// It is safe to modify the contents of the argument after Seek returns.
	Seek(key []byte) bool

	// Next moves the iterator to the next key/value pair.
	// It returns false if the iterator is exhausted.
	Next() bool

	// Prev moves the iterator to the previous key/value pair.
	// It returns false if the iterator is exhausted.
	Prev() bool
}

// CommonIterator is the interface that wraps common iterator methods.
type CommonIterator interface {
	IteratorSeeker

	// util.Releaser is the interface that wraps basic Release method.
	// When called Release will releases any resources associated with the
	// iterator.
	util.Releaser

	// util.ReleaseSetter is the interface that wraps the basic SetReleaser
	// method.
	util.ReleaseSetter

	// TODO: Remove this when ready.
	Valid() bool

	// Error returns any accumulated error. Exhausting all the key/value pairs
	// is not considered to be an error.
	Error() error
}

// Iterator iterates over a DB's key/value pairs in key order.
//
// When encounter an error any 'seeks method' will return false and will
// yield no key/value pairs. The error can be queried by calling the Error
// method. Calling Release is still necessary.
//
// An iterator must be released after use, but it is not necessary to read
// an iterator until exhaustion.
// Also, an iterator is not necessarily safe for concurrent use, but it is
// safe to use multiple iterators concurrently, with each in a dedicated
// goroutine.
type Iterator interface {
	CommonIterator

	// Key returns the key of the current key/value pair, or nil if done.
	// The caller should not modify the contents of the returned slice, and
	// its contents may change on the next call to any 'seeks method'.
	Key() []byte

	// Value returns the value of the current key/value pair, or nil if done.
	// The caller should not modify the contents of the returned slice, and
	// its contents may change on the next call to any 'seeks method'.
	Value() []byte
}

// ErrorCallbackSetter is the interface that wraps basic SetErrorCallback
// method.
//
// ErrorCallbackSetter implemented by indexed and merged iterator.
type ErrorCallbackSetter interface {
	// SetErrorCallback allows set an error callback of the corresponding
	// iterator. Use nil to clear the callback.
	SetErrorCallback(f func(err error))
}

type emptyIterator struct {
	util.BasicReleaser
	err error
}

func (i *emptyIterator) rErr() {
	if i.err == nil && i.Released() {
		i.err = ErrIterReleased
	}
}

func (*emptyIterator) Valid() bool            { return false }
func (i *emptyIterator) First() bool          { i.rErr(); return false }
func (i *emptyIterator) Last() bool           { i.rErr(); return false }
func (i *emptyIterator) Seek(key []byte) bool { i.rErr(); return false }
func (i *emptyIterator) Next() bool           { i.rErr(); return false }
func (i *emptyIterator) Prev() bool           { i.rErr(); return false }
func (*emptyIterator) Key() []byte            { return nil }
func (*emptyIterator) Value() []byte          { return nil }
func (i *emptyIterator) Error() error         { return i.err }

// NewEmptyIterator creates an empty iterator. The err parameter can be
// nil, but if not nil the given err will be returned by Error method.
func NewEmptyIterator(err error) Iterator {
	return &emptyIterator{err: err}
}
