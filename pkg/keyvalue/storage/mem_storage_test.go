package storage

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
	"bytes"
	"fmt"
	"testing"
)

func TestMemStorage(t *testing.T) {
	m := NewMemStorage()

	l, err := m.Lock()
	if err != nil {
		t.Fatal("storage lock failed(1): ", err)
	}
	_, err = m.Lock()
	if err == nil {
		t.Fatal("expect error for second storage lock attempt")
	} else {
		t.Logf("storage lock got error: %s (expected)", err)
	}
	l.Unlock()
	_, err = m.Lock()
	if err != nil {
		t.Fatal("storage lock failed(2): ", err)
	}

	w, err := m.Create(FileDesc{TypeTable, 1})
	if err != nil {
		t.Fatal("Storage.Create: ", err)
	}
	w.Write([]byte("abc"))
	w.Close()
	if fds, _ := m.List(TypeAll); len(fds) != 1 {
		t.Fatal("invalid GetFiles len")
	}
	buf := new(bytes.Buffer)
	r, err := m.Open(FileDesc{TypeTable, 1})
	if err != nil {
		t.Fatal("Open: got error: ", err)
	}
	buf.ReadFrom(r)
	r.Close()
	if got := buf.String(); got != "abc" {
		t.Fatalf("Read: invalid value, want=abc got=%s", got)
	}
	if _, err := m.Open(FileDesc{TypeTable, 1}); err != nil {
		t.Fatal("Open: got error: ", err)
	}
	if _, err := m.Open(FileDesc{TypeTable, 1}); err == nil {
		t.Fatal("expecting error")
	}
	m.Remove(FileDesc{TypeTable, 1})
	if fds, _ := m.List(TypeAll); len(fds) != 0 {
		t.Fatal("invalid GetFiles len", len(fds))
	}
	if _, err := m.Open(FileDesc{TypeTable, 1}); err == nil {
		t.Fatal("expecting error")
	}
}

func TestMemStorageRename(t *testing.T) {
	fd1 := FileDesc{Type: TypeTable, Num: 1}
	fd2 := FileDesc{Type: TypeTable, Num: 2}

	m := NewMemStorage()
	w, err := m.Create(fd1)
	if err != nil {
		t.Fatalf("Storage.Create: %v", err)
	}

	fmt.Fprintf(w, "abc")
	w.Close()

	rd, err := m.Open(fd1)
	if err != nil {
		t.Fatalf("Storage.Open(%v): %v", fd1, err)
	}
	rd.Close()

	fds, err := m.List(TypeAll)
	if err != nil {
		t.Fatalf("Storage.List: %v", err)
	}
	for _, fd := range fds {
		if !FileDescOk(fd) {
			t.Errorf("Storage.List -> FileDescOk(%q)", fd)
		}
	}

	err = m.Rename(fd1, fd2)
	if err != nil {
		t.Fatalf("Storage.Rename: %v", err)
	}

	rd, err = m.Open(fd2)
	if err != nil {
		t.Fatalf("Storage.Open(%v): %v", fd2, err)
	}
	rd.Close()

	fds, err = m.List(TypeAll)
	if err != nil {
		t.Fatalf("Storage.List: %v", err)
	}
	for _, fd := range fds {
		if !FileDescOk(fd) {
			t.Errorf("Storage.List -> FileDescOk(%q)", fd)
		}
	}
}
