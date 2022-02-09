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

	"github.com/bhojpur/dbm/pkg/keyed"
)

// KeyedDBStore implements CacheStore provide local machine
type KeyedDBStore struct {
	store *keyed.DB
	Debug bool
	v     interface{}
}

var _ CacheStore = &KeyedDBStore{}

// NewKeyedDBStore creates a Keyed DB store
func NewKeyedDBStore(dbfile string) (*KeyedDBStore, error) {
	db := &KeyedDBStore{}
	h, err := keyed.OpenFile(dbfile, nil)
	if err != nil {
		return nil, err
	}
	db.store = h
	return db, nil
}

// Put implements CacheStore
func (s *KeyedDBStore) Put(key string, value interface{}) error {
	val, err := Encode(value)
	if err != nil {
		if s.Debug {
			log.Println("[KeyedDB]EncodeErr: ", err, "Key:", key)
		}
		return err
	}
	err = s.store.Put([]byte(key), val, nil)
	if err != nil {
		if s.Debug {
			log.Println("[KeyedDB]PutErr: ", err, "Key:", key)
		}
		return err
	}
	if s.Debug {
		log.Println("[KeyedDB]Put: ", key)
	}
	return err
}

// Get implements CacheStore
func (s *KeyedDBStore) Get(key string) (interface{}, error) {
	data, err := s.store.Get([]byte(key), nil)
	if err != nil {
		if s.Debug {
			log.Println("[KeyedDB]GetErr: ", err, "Key:", key)
		}
		if err == keyed.ErrNotFound {
			return nil, ErrNotExist
		}
		return nil, err
	}
	err = Decode(data, &s.v)
	if err != nil {
		if s.Debug {
			log.Println("[KeyedDB]DecodeErr: ", err, "Key:", key)
		}
		return nil, err
	}
	if s.Debug {
		log.Println("[KeyedDB]Get: ", key, s.v)
	}
	return s.v, err
}

// Del implements CacheStore
func (s *KeyedDBStore) Del(key string) error {
	err := s.store.Delete([]byte(key), nil)
	if err != nil {
		if s.Debug {
			log.Println("[KeyedDB]DelErr: ", err, "Key:", key)
		}
		return err
	}
	if s.Debug {
		log.Println("[KeyedDB]Del: ", key)
	}
	return err
}

// Close implements CacheStore
func (s *KeyedDBStore) Close() {
	s.store.Close()
}
