package keyvalue

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
	"encoding/binary"
	"fmt"

	"github.com/bhojpur/dbm/pkg/keyvalue/errors"
	"github.com/bhojpur/dbm/pkg/keyvalue/storage"
)

// ErrInternalKeyCorrupted records internal key corruption.
type ErrInternalKeyCorrupted struct {
	Ikey   []byte
	Reason string
}

func (e *ErrInternalKeyCorrupted) Error() string {
	return fmt.Sprintf("keyvalue: internal key %q corrupted: %s", e.Ikey, e.Reason)
}

func newErrInternalKeyCorrupted(ikey []byte, reason string) error {
	return errors.NewErrCorrupted(storage.FileDesc{}, &ErrInternalKeyCorrupted{append([]byte{}, ikey...), reason})
}

type keyType uint

func (kt keyType) String() string {
	switch kt {
	case keyTypeDel:
		return "d"
	case keyTypeVal:
		return "v"
	}
	return fmt.Sprintf("<invalid:%#x>", uint(kt))
}

// Value types encoded as the last component of internal keys.
// Don't modify; this value are saved to disk.
const (
	keyTypeDel = keyType(0)
	keyTypeVal = keyType(1)
)

// keyTypeSeek defines the keyType that should be passed when constructing an
// internal key for seeking to a particular sequence number (since we
// sort sequence numbers in decreasing order and the value type is
// embedded as the low 8 bits in the sequence number in internal keys,
// we need to use the highest-numbered ValueType, not the lowest).
const keyTypeSeek = keyTypeVal

const (
	// Maximum value possible for sequence number; the 8-bits are
	// used by value type, so its can packed together in single
	// 64-bit integer.
	keyMaxSeq = (uint64(1) << 56) - 1
	// Maximum value possible for packed sequence number and type.
	keyMaxNum = (keyMaxSeq << 8) | uint64(keyTypeSeek)
)

// Maximum number encoded in bytes.
var keyMaxNumBytes = make([]byte, 8)

func init() {
	binary.LittleEndian.PutUint64(keyMaxNumBytes, keyMaxNum)
}

type internalKey []byte

func makeInternalKey(dst, ukey []byte, seq uint64, kt keyType) internalKey {
	if seq > keyMaxSeq {
		panic("keyvalue: invalid sequence number")
	} else if kt > keyTypeVal {
		panic("keyvalue: invalid type")
	}

	dst = ensureBuffer(dst, len(ukey)+8)
	copy(dst, ukey)
	binary.LittleEndian.PutUint64(dst[len(ukey):], (seq<<8)|uint64(kt))
	return internalKey(dst)
}

func parseInternalKey(ik []byte) (ukey []byte, seq uint64, kt keyType, err error) {
	if len(ik) < 8 {
		return nil, 0, 0, newErrInternalKeyCorrupted(ik, "invalid length")
	}
	num := binary.LittleEndian.Uint64(ik[len(ik)-8:])
	seq, kt = uint64(num>>8), keyType(num&0xff)
	if kt > keyTypeVal {
		return nil, 0, 0, newErrInternalKeyCorrupted(ik, "invalid type")
	}
	ukey = ik[:len(ik)-8]
	return
}

func validInternalKey(ik []byte) bool {
	_, _, _, err := parseInternalKey(ik)
	return err == nil
}

func (ik internalKey) assert() {
	if ik == nil {
		panic("keyvalue: nil internalKey")
	}
	if len(ik) < 8 {
		panic(fmt.Sprintf("keyvalue: internal key %q, len=%d: invalid length", []byte(ik), len(ik)))
	}
}

func (ik internalKey) ukey() []byte {
	ik.assert()
	return ik[:len(ik)-8]
}

func (ik internalKey) num() uint64 {
	ik.assert()
	return binary.LittleEndian.Uint64(ik[len(ik)-8:])
}

func (ik internalKey) parseNum() (seq uint64, kt keyType) {
	num := ik.num()
	seq, kt = uint64(num>>8), keyType(num&0xff)
	if kt > keyTypeVal {
		panic(fmt.Sprintf("keyvalue: internal key %q, len=%d: invalid type %#x", []byte(ik), len(ik), kt))
	}
	return
}

func (ik internalKey) String() string {
	if ik == nil {
		return "<nil>"
	}

	if ukey, seq, kt, err := parseInternalKey(ik); err == nil {
		return fmt.Sprintf("%s,%s%d", shorten(string(ukey)), kt, seq)
	}
	return fmt.Sprintf("<invalid:%#x>", []byte(ik))
}
