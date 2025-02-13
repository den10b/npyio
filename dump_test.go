// Copyright 2020 The npyio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package npyio

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestDump(t *testing.T) {
	for _, tc := range []struct {
		name string
		want string
	}{
		{
			name: "testdata/data_float32_2x3_corder.npy",
			want: "testdata/data_float32_2x3_corder.npy.txt",
		},
		{
			name: "testdata/data_float32_2x3_forder.npy",
			want: "testdata/data_float32_2x3_forder.npy.txt",
		},
		{
			name: "testdata/data_float64_2x3x4_corder.npy",
			want: "testdata/data_float64_2x3x4_corder.npy.txt",
		},
		{
			name: "testdata/data_float64_corder.npz",
			want: "testdata/data_float64_corder.npz.txt",
		},
		{
			name: "testdata/data_float64_forder.npz",
			want: "testdata/data_float64_forder.npz.txt",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.name)
			if err != nil {
				t.Fatalf("could not open %q: %+v", tc.name, err)
			}
			defer f.Close()

			o := new(strings.Builder)
			err = Dump(o, f)
			if err != nil {
				t.Fatalf("could not dump %q: %+v", tc.name, err)
			}

			want, err := os.ReadFile(tc.want)
			if err != nil {
				t.Fatalf("could not read reference file %q: %+v", tc.want, err)
			}

			if got, want := o.String(), string(want); got != want {
				t.Fatalf(
					"invalid dump:\ngot:\n%s\nwant:\n%s\n",
					got, want,
				)
			}
		})
	}
}

func TestDumpSeeker(t *testing.T) {
	for _, tc := range []struct {
		name string
		want string
	}{
		{
			name: "testdata/data_float32_2x3_corder.npy",
			want: "testdata/data_float32_2x3_corder.npy.txt",
		},
		{
			name: "testdata/data_float32_2x3_forder.npy",
			want: "testdata/data_float32_2x3_forder.npy.txt",
		},
		{
			name: "testdata/data_float64_2x3x4_corder.npy",
			want: "testdata/data_float64_2x3x4_corder.npy.txt",
		},
		{
			name: "testdata/data_float64_corder.npz",
			want: "testdata/data_float64_corder.npz.txt",
		},
		{
			name: "testdata/data_float64_forder.npz",
			want: "testdata/data_float64_forder.npz.txt",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.name)
			if err != nil {
				t.Fatalf("could not open %q: %+v", tc.name, err)
			}
			defer f.Close()

			type namer interface{ Name() string }

			r := struct {
				io.Seeker
				io.ReaderAt
				namer
			}{
				Seeker:   f,
				ReaderAt: f,
				namer:    f,
			}
			o := new(strings.Builder)
			err = Dump(o, r)
			if err != nil {
				t.Fatalf("could not dump %q: %+v", tc.name, err)
			}

			want, err := os.ReadFile(tc.want)
			if err != nil {
				t.Fatalf("could not read reference file %q: %+v", tc.want, err)
			}

			if got, want := o.String(), string(want); got != want {
				t.Fatalf(
					"invalid dump:\ngot:\n%s\nwant:\n%s\n",
					got, want,
				)
			}
		})
	}
}
