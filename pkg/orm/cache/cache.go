package cache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bhojpur/dbm/pkg/orm/schema"
)

const (
	// CacheExpired is default cache expired time
	CacheExpired = 60 * time.Minute
	// CacheMaxMemory is not use now
	CacheMaxMemory = 256
	// CacheGcInterval represents interval time to clear all expired nodes
	CacheGcInterval = 10 * time.Minute
	// CacheGcMaxRemoved represents max nodes removed when gc
	CacheGcMaxRemoved = 20
)

// list all the errors
var (
	ErrCacheMiss = errors.New("orm/cache: key not found")
	ErrNotStored = errors.New("orm/cache: not stored")
	// ErrNotExist record does not exist error
	ErrNotExist = errors.New("Record does not exist")
)

// CacheStore is a interface to store cache
type CacheStore interface {
	// key is primary key or composite primary key
	// value is struct's pointer
	// key format : <tablename>-p-<pk1>-<pk2>...
	Put(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Del(key string) error
}

// Cacher is an interface to provide cache
// id format : u-<pk1>-<pk2>...
type Cacher interface {
	GetIds(tableName, sql string) interface{}
	GetBean(tableName string, id string) interface{}
	PutIds(tableName, sql string, ids interface{})
	PutBean(tableName string, id string, obj interface{})
	DelIds(tableName, sql string)
	DelBean(tableName string, id string)
	ClearIds(tableName string)
	ClearBeans(tableName string)
}

func encodeIds(ids []schema.PK) (string, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(ids)
	return buf.String(), err
}
func decodeIds(s string) ([]schema.PK, error) {
	pks := make([]schema.PK, 0)
	dec := gob.NewDecoder(strings.NewReader(s))
	err := dec.Decode(&pks)
	return pks, err
}

// GetCacheSql returns cacher PKs via SQL
func GetCacheSql(m Cacher, tableName, sql string, args interface{}) ([]schema.PK, error) {
	bytes := m.GetIds(tableName, GenSqlKey(sql, args))
	if bytes == nil {
		return nil, errors.New("Not Exist")
	}
	return decodeIds(bytes.(string))
}

// PutCacheSql puts cacher SQL and PKs
func PutCacheSql(m Cacher, ids []schema.PK, tableName, sql string, args interface{}) error {
	bytes, err := encodeIds(ids)
	if err != nil {
		return err
	}
	m.PutIds(tableName, GenSqlKey(sql, args), bytes)
	return nil
}

// GenSqlKey generates cache key
func GenSqlKey(sql string, args interface{}) string {
	return fmt.Sprintf("%v-%v", sql, args)
}
