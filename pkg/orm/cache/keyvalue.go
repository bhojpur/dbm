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

import (
	"log"

	"github.com/bhojpur/dbm/pkg/keyvalue"
)

// KeyValueDBStore implements CacheStore provide local machine
type KeyValueDBStore struct {
	store *keyvalue.DB
	Debug bool
	v     interface{}
}

var _ CacheStore = &KeyValueDBStore{}

// NewKeyValueDBStore creates a KeyValue DB store
func NewKeyValueDBStore(dbfile string) (*KeyValueDBStore, error) {
	db := &KeyValueDBStore{}
	h, err := keyvalue.OpenFile(dbfile, nil)
	if err != nil {
		return nil, err
	}
	db.store = h
	return db, nil
}

// Put implements CacheStore
func (s *KeyValueDBStore) Put(key string, value interface{}) error {
	val, err := Encode(value)
	if err != nil {
		if s.Debug {
			log.Println("[KeyValueDB]EncodeErr: ", err, "Key:", key)
		}
		return err
	}
	err = s.store.Put([]byte(key), val, nil)
	if err != nil {
		if s.Debug {
			log.Println("[KeyValueDB]PutErr: ", err, "Key:", key)
		}
		return err
	}
	if s.Debug {
		log.Println("[KeyValueDB]Put: ", key)
	}
	return err
}

// Get implements CacheStore
func (s *KeyValueDBStore) Get(key string) (interface{}, error) {
	data, err := s.store.Get([]byte(key), nil)
	if err != nil {
		if s.Debug {
			log.Println("[KeyValueDB]GetErr: ", err, "Key:", key)
		}
		if err == keyvalue.ErrNotFound {
			return nil, ErrNotExist
		}
		return nil, err
	}
	err = Decode(data, &s.v)
	if err != nil {
		if s.Debug {
			log.Println("[KeyValueDB]DecodeErr: ", err, "Key:", key)
		}
		return nil, err
	}
	if s.Debug {
		log.Println("[KeyValueDB]Get: ", key, s.v)
	}
	return s.v, err
}

// Del implements CacheStore
func (s *KeyValueDBStore) Del(key string) error {
	err := s.store.Delete([]byte(key), nil)
	if err != nil {
		if s.Debug {
			log.Println("[KeyValueDB]DelErr: ", err, "Key:", key)
		}
		return err
	}
	if s.Debug {
		log.Println("[KeyValueDB]Del: ", key)
	}
	return err
}

// Close implements CacheStore
func (s *KeyValueDBStore) Close() {
	s.store.Close()
}
