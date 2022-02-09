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
	"os"
)

type plan9FileLock struct {
	f *os.File
}

func (fl *plan9FileLock) release() error {
	return fl.f.Close()
}

func newFileLock(path string, readOnly bool) (fl fileLock, err error) {
	var (
		flag int
		perm os.FileMode
	)
	if readOnly {
		flag = os.O_RDONLY
	} else {
		flag = os.O_RDWR
		perm = os.ModeExclusive
	}
	f, err := os.OpenFile(path, flag, perm)
	if os.IsNotExist(err) {
		f, err = os.OpenFile(path, flag|os.O_CREATE, perm|0644)
	}
	if err != nil {
		return
	}
	fl = &plan9FileLock{f: f}
	return
}

func rename(oldpath, newpath string) error {
	if _, err := os.Stat(newpath); err == nil {
		if err := os.Remove(newpath); err != nil {
			return err
		}
	}

	return os.Rename(oldpath, newpath)
}

func syncDir(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := f.Sync(); err != nil {
		return err
	}
	return nil
}
