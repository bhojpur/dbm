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
	"syscall"
	"unsafe"
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procMoveFileExW = modkernel32.NewProc("MoveFileExW")
)

const (
	_MOVEFILE_REPLACE_EXISTING = 1
)

type windowsFileLock struct {
	fd syscall.Handle
}

func (fl *windowsFileLock) release() error {
	return syscall.Close(fl.fd)
}

func newFileLock(path string, readOnly bool) (fl fileLock, err error) {
	pathp, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return
	}
	var access, shareMode uint32
	if readOnly {
		access = syscall.GENERIC_READ
		shareMode = syscall.FILE_SHARE_READ | syscall.FILE_SHARE_WRITE
	} else {
		access = syscall.GENERIC_READ | syscall.GENERIC_WRITE
	}
	fd, err := syscall.CreateFile(pathp, access, shareMode, nil, syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err == syscall.ERROR_FILE_NOT_FOUND {
		fd, err = syscall.CreateFile(pathp, access, shareMode, nil, syscall.OPEN_ALWAYS, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	}
	if err != nil {
		return
	}
	fl = &windowsFileLock{fd: fd}
	return
}

func moveFileEx(from *uint16, to *uint16, flags uint32) error {
	r1, _, e1 := syscall.Syscall(procMoveFileExW.Addr(), 3, uintptr(unsafe.Pointer(from)), uintptr(unsafe.Pointer(to)), uintptr(flags))
	if r1 == 0 {
		if e1 != 0 {
			return error(e1)
		}
		return syscall.EINVAL
	}
	return nil
}

func rename(oldpath, newpath string) error {
	from, err := syscall.UTF16PtrFromString(oldpath)
	if err != nil {
		return err
	}
	to, err := syscall.UTF16PtrFromString(newpath)
	if err != nil {
		return err
	}
	return moveFileEx(from, to, _MOVEFILE_REPLACE_EXISTING)
}

func syncDir(name string) error { return nil }
