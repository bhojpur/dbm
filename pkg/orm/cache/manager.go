package cache

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

import "sync"

// Manager represents a cache manager
type Manager struct {
	cacher             Cacher
	disableGlobalCache bool
	cachers            map[string]Cacher
	cacherLock         sync.RWMutex
}

// NewManager creates a cache manager
func NewManager() *Manager {
	return &Manager{
		cachers: make(map[string]Cacher),
	}
}

// SetDisableGlobalCache disable global cache or not
func (mgr *Manager) SetDisableGlobalCache(disable bool) {
	if mgr.disableGlobalCache != disable {
		mgr.disableGlobalCache = disable
	}
}

// SetCacher set cacher of table
func (mgr *Manager) SetCacher(tableName string, cacher Cacher) {
	mgr.cacherLock.Lock()
	mgr.cachers[tableName] = cacher
	mgr.cacherLock.Unlock()
}

// GetCacher returns a cache of a table
func (mgr *Manager) GetCacher(tableName string) Cacher {
	var cacher Cacher
	var ok bool
	mgr.cacherLock.RLock()
	cacher, ok = mgr.cachers[tableName]
	mgr.cacherLock.RUnlock()
	if !ok && !mgr.disableGlobalCache {
		cacher = mgr.cacher
	}
	return cacher
}

// SetDefaultCacher set the default cacher. ORM's default not enable cacher.
func (mgr *Manager) SetDefaultCacher(cacher Cacher) {
	mgr.cacher = cacher
}

// GetDefaultCacher returns the default cacher
func (mgr *Manager) GetDefaultCacher() Cacher {
	return mgr.cacher
}
