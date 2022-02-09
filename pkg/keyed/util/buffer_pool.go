package util

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
	"fmt"
	"sync"
	"sync/atomic"
)

// BufferPool is a 'buffer pool'.
type BufferPool struct {
	pool     [6]sync.Pool
	baseline [5]int

	get     uint32
	put     uint32
	less    uint32
	equal   uint32
	greater uint32
	miss    uint32
}

func (p *BufferPool) poolNum(n int) int {
	for i, x := range p.baseline {
		if n <= x {
			return i
		}
	}
	return len(p.baseline)
}

// Get returns buffer with length of n.
func (p *BufferPool) Get(n int) []byte {
	if p == nil {
		return make([]byte, n)
	}
	atomic.AddUint32(&p.get, 1)

	poolNum := p.poolNum(n)

	b := p.pool[poolNum].Get().(*[]byte)

	if cap(*b) == 0 {
		// If we grabbed nothing, increment the miss stats.
		atomic.AddUint32(&p.miss, 1)

		if poolNum == len(p.baseline) {
			*b = make([]byte, n)
			return *b
		}

		*b = make([]byte, p.baseline[poolNum])
		*b = (*b)[:n]
		return *b
	} else {
		// If there is enough capacity in the bytes grabbed, resize the length
		// to n and return.
		if n < cap(*b) {
			atomic.AddUint32(&p.less, 1)
			*b = (*b)[:n]
			return *b
		} else if n == cap(*b) {
			atomic.AddUint32(&p.equal, 1)
			*b = (*b)[:n]
			return *b
		} else if n > cap(*b) {
			atomic.AddUint32(&p.greater, 1)
		}
	}

	if poolNum == len(p.baseline) {
		*b = make([]byte, n)
		return *b
	}
	*b = make([]byte, p.baseline[poolNum])
	*b = (*b)[:n]
	return *b
}

// Put adds given buffer to the pool.
func (p *BufferPool) Put(b []byte) {
	if p == nil {
		return
	}

	poolNum := p.poolNum(cap(b))

	atomic.AddUint32(&p.put, 1)
	p.pool[poolNum].Put(&b)
}

func (p *BufferPool) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("BufferPool{B·%d G·%d P·%d <·%d =·%d >·%d M·%d}",
		p.baseline, p.get, p.put, p.less, p.equal, p.greater, p.miss)
}

// NewBufferPool creates a new initialized 'buffer pool'.
func NewBufferPool(baseline int) *BufferPool {
	if baseline <= 0 {
		panic("baseline can't be <= 0")
	}
	bufPool := &BufferPool{
		baseline: [...]int{baseline / 4, baseline / 2, baseline, baseline * 2, baseline * 4},
		pool: [6]sync.Pool{
			sync.Pool{
				New: func() interface{} { return new([]byte) },
			},
			sync.Pool{
				New: func() interface{} { return new([]byte) },
			},
			sync.Pool{
				New: func() interface{} { return new([]byte) },
			},
			sync.Pool{
				New: func() interface{} { return new([]byte) },
			},
			sync.Pool{
				New: func() interface{} { return new([]byte) },
			},
			sync.Pool{
				New: func() interface{} { return new([]byte) },
			},
		},
	}

	return bufPool
}
