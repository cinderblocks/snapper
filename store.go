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
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/golang/snappy"
)

type AssetStore interface {
	Load(hash string) (io.ReadCloser, error)
	Exists(hash string) bool
	Store(data string) (string, error)
	GetAsBase64(hash string) (string, error)
}

type assetStore struct {
	dataDir  string
	spoolDir string
}

func CreateAssetStore(dataDir, spoolDir string) AssetStore {
	return &assetStore{
		dataDir:  dataDir,
		spoolDir: spoolDir,
	}
}

func (a assetStore) makePath(hash string) string {
	return path.Join(a.dataDir, hash[0:3], hash[3:6], hash)
}

func (a assetStore) Load(hash string) (io.ReadCloser, error) {
	spath := a.makePath(hash)
	zipped := false
	snap := false
	f, e := os.Open(spath)
	if e != nil {

		f, e = os.Open(spath + ".snappy")
		if e != nil {
			f, e = os.Open(spath + ".gz")
			zipped = true
		} else {
			snap = true
		}
	}
	if e != nil {
		return nil, e
	}

	if !zipped && !snap {
		return f, nil
	}

	if zipped {
		gzipreader, e := gzip.NewReader(f)
		if e != nil {
			return nil, e
		}
		return &assetReader{
			f:      f,
			gzip:   gzipreader,
			snappy: nil,
		}, nil
	}
	return &assetReader{
		f:      f,
		gzip:   nil,
		snappy: snappy.NewReader(f),
	}, nil
}

func (a assetStore) makeHash(data []byte) string {
	shabuf := sha256.Sum256(data)
	return strings.ToUpper(hex.EncodeToString(shabuf[0:len(shabuf)]))
}

func (a assetStore) Exists(hash string) bool {
	_, exists := a.exists(hash)
	return exists
}

func (a assetStore) exists(hash string) (string, bool) {
	spath := a.makePath(hash)
	if _, err := os.Stat(spath); err == nil {
		// already exists
		return spath, true
	}

	gzPath := spath + ".gz"
	if _, err := os.Stat(gzPath); err == nil {
		// already exists
		return gzPath, true
	}

	snappyPath := spath + ".snappy"
	if _, err := os.Stat(snappyPath); err == nil {
		// already exists
		return snappyPath, true
	}
	return snappyPath, false
}

func (a assetStore) preparePath(hash string) (string, error) {
	spath, exists := a.exists(hash)
	if exists {
		return "", os.ErrExist
	}

	err := os.MkdirAll(path.Dir(spath), 0773)
	if err != nil {
		if err != os.ErrExist {
			return "", err
		}
	}

	return spath, nil
}

func (a assetStore) Store(data string) (string, error) {
	buffer, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	hash := a.makeHash(buffer)
	spath, err := a.preparePath(hash)

	if os.IsExist(err) {
		return hash, nil
	}

	// All checks done, now create temporary file instead of real one
	// To have some transaction safety
	tempPath := strings.Replace(spath, a.dataDir, a.spoolDir, 1)
	err = os.MkdirAll(path.Dir(tempPath), 0773)
	if err != nil {
		return hash, err
	}

	f, err := os.Create(tempPath)
	if err != nil {
		return hash, err
	}
	defer f.Close()

	if path.Ext(spath) == ".snappy" {
		snap := snappy.NewWriter(f)
		_, err = io.Copy(snap, bytes.NewReader(buffer))
	} else {
		gzw := gzip.NewWriter(f)
		defer gzw.Close()
		_, err = io.Copy(gzw, bytes.NewReader(buffer))
	}

	if err == nil {
		// File writing is done now move the temp file to the real location
		err = os.Rename(tempPath, spath)
		if err != nil {
			os.Remove(tempPath)
			if os.IsExist(err) {
				err = nil
			}
		}
	}
	return hash, err
}

func (a assetStore) GetAsBase64(hash string) (string, error) {
	reader, err := a.Load(hash)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(reader)
	if err == nil {
		return base64.StdEncoding.EncodeToString(data), nil
	}
	return "", err
}
