// Copyright 2016 The npyio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package npyio provides read/write access to files following the NumPy data file format:
//  http://docs.scipy.org/doc/numpy-1.10.1/neps/npy-format.html
//
// Supported types
//
// npyio supports r/w of scalars, arrays, slices and mat64.Dense.
// Supported scalars are:
//  - bool,
//  - (u)int{,8,16,32,64},
//  - float{32,64},
//  - complex{64,128}
//
// Reading
//
// Reading from a NumPy data file can be performed like so:
//
//  f, err := os.Open("data.npy")
//  var m mat64.Dense
//  err = npyio.Read(f, &m)
//  fmt.Printf("data = %v\n", mat64.Formatted(&m, mat64.Prefix("       ")))
//
// npyio can also read data directly into slices, arrays or scalars, provided
// the on-disk data type and the provided one match.
//
// Example:
//  var data []float64
//  err = npyio.Read(f, &data)
//
//  var data uint64
//  err = npyio.Read(f, &data)
//
// Writing
//
// Writing into a NumPy data file can be done like so:
//
//  f, err := os.Create("data.npy")
//  var m mat64.Dense = ...
//  err = npyio.Write(f, m)
//
// Scalars, arrays and slices are also supported:
//
//  var data []float64 = ...
//  err = npyio.Write(f, data)
//
//  var data int64 = 42
//  err = npyio.Write(f, data)
//
//  var data [42]complex128 = ...
//  err = npyio.Write(f, data)
package npyio

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

var (
	errNilPtr = errors.New("npyio: nil pointer")
	errNotPtr = errors.New("npyio: expected a pointer to a value")
	errDims   = errors.New("npyio: invalid dimensions")
	errNoConv = errors.New("npyio: no legal type conversion")

	ble = binary.LittleEndian

	// ErrInvalidNumPyFormat is the error returned by NewReader when
	// the underlying io.Reader is not a valid or recognized NumPy data
	// file format.
	ErrInvalidNumPyFormat = errors.New("npyio: not a valid NumPy file format")

	// ErrTypeMismatch is the error returned by Reader when the on-disk
	// data type and the user provided one do NOT match.
	ErrTypeMismatch = errors.New("npyio: types don't match")

	// Magic header present at the start of a NumPy data file format.
	// See http://docs.scipy.org/doc/numpy-1.10.1/neps/npy-format.html
	Magic = [6]byte{'\x93', 'N', 'U', 'M', 'P', 'Y'}
)

// Header describes the data content of a NumPy data file.
type Header struct {
	Major byte // data file major version
	Minor byte // data file minor version
	Descr struct {
		Type    string // data type of array elements ('<i8', '<f4', ...)
		Fortran bool   // whether the array data is stored in Fortran-order (col-major)
		Shape   []int  // array shape (e.g. [2,3] a 2-rows, 3-cols array
	}
}

// newHeader creates a new Header with the major/minor version numbers that npyio currently supports.
func newHeader() Header {
	return Header{
		Major: 2,
		Minor: 0,
	}
}

func (h Header) String() string {
	return fmt.Sprintf("Header{Major:%v, Minor:%v, Descr:{Type:%v, Fortran:%v, Shape:%v}}",
		int(h.Major),
		int(h.Minor),
		h.Descr.Type,
		h.Descr.Fortran,
		h.Descr.Shape,
	)
}

var (
	uint8Type      = reflect.TypeOf((*uint8)(nil)).Elem()
	uint16Type     = reflect.TypeOf((*uint16)(nil)).Elem()
	uint32Type     = reflect.TypeOf((*uint32)(nil)).Elem()
	uint64Type     = reflect.TypeOf((*uint64)(nil)).Elem()
	int8Type       = reflect.TypeOf((*int8)(nil)).Elem()
	int16Type      = reflect.TypeOf((*int16)(nil)).Elem()
	int32Type      = reflect.TypeOf((*int32)(nil)).Elem()
	int64Type      = reflect.TypeOf((*int64)(nil)).Elem()
	float32Type    = reflect.TypeOf((*float32)(nil)).Elem()
	float64Type    = reflect.TypeOf((*float64)(nil)).Elem()
	complex64Type  = reflect.TypeOf((*complex64)(nil)).Elem()
	complex128Type = reflect.TypeOf((*complex128)(nil)).Elem()
)
