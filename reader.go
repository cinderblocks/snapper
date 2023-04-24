// Copyright (c) 2015-2018 Cinderblocks Design Co.
//
// This file is part of snapper
// (see https://bitbucket.org/cinderblocks/snapper).
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"compress/gzip"
	"os"

	"github.com/golang/snappy"
)

type assetReader struct {
	f      *os.File
	gzip   *gzip.Reader
	snappy *snappy.Reader
}

func (a assetReader) Read(p []byte) (n int, err error) {
	if a.f == nil || (a.gzip == nil && a.snappy == nil) {
		return 0, os.ErrInvalid
	}
	if a.gzip != nil {
		return a.gzip.Read(p)
	}
	return a.snappy.Read(p)
}

func (a assetReader) Close() error {
	if a.gzip != nil {
		a.gzip.Close()
	}
	if a.f != nil {
		a.f.Close()
	}
	// snappy.Reader has no Close
	return nil
}
