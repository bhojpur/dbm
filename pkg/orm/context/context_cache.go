package context

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

// ContextCache is the interface that operates the cache data.
type ContextCache interface {
	// Put puts value into cache with key.
	Put(key string, val interface{})
	// Get gets cached value by given key.
	Get(key string) interface{}
}
type memoryContextCache map[string]interface{}

// NewMemoryContextCache return memoryContextCache
func NewMemoryContextCache() memoryContextCache {
	return make(map[string]interface{})
}

// Put puts value into cache with key.
func (m memoryContextCache) Put(key string, val interface{}) {
	m[key] = val
}

// Get gets cached value by given key.
func (m memoryContextCache) Get(key string) interface{} {
	return m[key]
}
