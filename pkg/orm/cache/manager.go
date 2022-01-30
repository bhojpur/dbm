package cache

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
