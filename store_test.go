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
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"
)

var testingAssetStore *assetStore = nil

const (
	testFileDataContent     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	testFileDataContentB64  = "QUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVo="
	testFileDataContentHash = "D6EC6898DE87DDAC6E5B3611708A7AA1C2D298293349CC1A6C299A1DB7149D38"

	emptyTestFileDataContentHash = "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855"
	emptyTestFileDataContentB64  = ""
	emptyTestFileDataContent     = ""
)

func init() {
	rand.Seed(time.Now().UnixNano())
	randData := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	randName := make([]rune, 8)
	for i := range randName {
		randName[i] = randData[rand.Intn(len(randData))]
	}

	tempBaseDir := os.TempDir()
	dataDir := path.Join(tempBaseDir, "testing", "store", string(randName))
	testingAssetStore = &assetStore{
		dataDir:  path.Join(dataDir, "data"),
		spoolDir: path.Join(dataDir, "tmp"),
	}

}

func TestAssetStore_makePath(t *testing.T) {
	result := testingAssetStore.makePath("0011223344")
	expected := path.Join(testingAssetStore.dataDir, "001", "122", "0011223344")
	if result != expected {
		t.Fail()
		t.Logf("TestAssetStoremakePath: Expected: %v Got: %v", expected, result)
	}
}

func TestAssetStore_MakeHash(t *testing.T) {
	expected := "84D89877F0D4041EFB6BF91A16F0248F2FD573E6AF05C19F96BEDB9F882F7882"
	result := testingAssetStore.makeHash([]byte("0123456789"))
	if expected != result {
		t.Fail()
		t.Logf("TestAssetStoreMakeHash: Expected: %v Got: %v", expected, result)
	}
}

func TestAssetStore_preparePath_NonExisting(t *testing.T) {
	hash := "84D89877F0D4041EFB6BF91A16F0248F2FD573E6AF05C19F96BEDB9F882F7882"
	expected := path.Join(testingAssetStore.dataDir, "84D", "898", "84D89877F0D4041EFB6BF91A16F0248F2FD573E6AF05C19F96BEDB9F882F7882") + ".snappy"
	result, resultExists := testingAssetStore.preparePath(hash)
	if os.IsExist(resultExists) {
		t.Fail()
		t.Logf("TestAssetStorePreparePathNonExisting: Path exists, failure or non clean testing environment")
	} else if result != expected {
		t.Fail()
		t.Logf("TestAssetStorePreparePathNonExisting: Expected: %v Got: %v", expected, result)
	}
}

func TestAssetStore_preparePath_Existing(t *testing.T) {
	hash := "84D89877F0D4041EFB6BF91A16F0248F2FD573E6AF05C19F96BEDB9F882F7882"
	expected := path.Join(testingAssetStore.dataDir, "84D", "898", "84D89877F0D4041EFB6BF91A16F0248F2FD573E6AF05C19F96BEDB9F882F7882")
	os.MkdirAll(path.Dir(expected), 0664)
	f, err := os.Create(expected)
	if err != nil {
	} else {
		f.Close()
		result, resultExists := testingAssetStore.preparePath(hash)
		if os.IsNotExist(resultExists) {
			t.Fail()
			t.Logf("TestAssetStorePreparePathExisting: Path does not exists, failure or non clean testing environment")
		} else if result != "" {
			t.Fail()
			t.Logf("TestAssetStorePreparePathExisting: Expected: %v Got: %v", expected, result)
		}
	}
}

func TestAssetStore_exists_NonExisting(t *testing.T) {
	hash := "DEADBEEF"
	expected := path.Join(testingAssetStore.dataDir, "DEA", "DBE", "DEADBEEF") + ".snappy"
	result, resultExists := testingAssetStore.exists(hash)
	if resultExists {
		t.Fail()
		t.Logf("TestAssetStore_exists: Path does not exists, failure or non clean testing environment")
	} else if result != expected {
		t.Fail()
		t.Logf("TestAssetStorePreparePathExisting: Expected: %v Got: %v", expected, result)
	}
}

func TestAssetStore_exists_NonGZ(t *testing.T) {
	hash := "DEADBEEF"
	expected := path.Join(testingAssetStore.dataDir, "DEA", "DBE", "DEADBEEF")
	os.MkdirAll(path.Dir(expected), 0664)
	f, err := os.Create(expected)
	if err != nil {
	} else {
		f.Close()
		result, resultExists := testingAssetStore.exists(hash)
		if !resultExists {
			t.Fail()
			t.Logf("TestAssetStore_exists: Path does not exists, failure or non clean testing environment")
		} else if result != expected {
			t.Fail()
			t.Logf("TestAssetStorePreparePathExisting: Expected: %v Got: %v", expected, result)
		}
	}
}

func TestAssetStore_Store(t *testing.T) {
	dataHash := testFileDataContentHash
	hashResult, err := testingAssetStore.Store(testFileDataContentB64)
	if hashResult != dataHash {
		t.Fail()
		t.Logf("Expected hash: %v Got: %v", dataHash, hashResult)
	} else if err != nil {
		t.Fail()
		t.Logf("Store failed: %v", err)
	}
	path := testingAssetStore.makePath(dataHash) + ".snappy"
	if _, err := os.Stat(path); err != nil {
		t.Fail()
		t.Logf("Store did not store data: %v => %v", path, err)
	}
}

func TestAssetStore_Load(t *testing.T) {
	readCloser, err := testingAssetStore.Load(testFileDataContentHash)
	if err != nil {
		t.Fail()
		t.Logf("Failed to load file from previous test: %v", err)
	} else if readCloser == nil {
		t.Fail()
		t.Logf("Failed to get read closer without error")
	} else {
		receivedData, err := ioutil.ReadAll(readCloser)
		if err != nil {
			t.Logf("Failure on reading received data stream: %v", err)
			t.Fail()
		} else if string(receivedData) != testFileDataContent {
			t.Fail()
			t.Logf("Received data is not what was written to the file: Expected: %v Got: %v", testFileDataContent, string(receivedData))
		}
	}
	if readCloser != nil {
		readCloser.Close()
	}
}

func TestAssetStore_GetAsBase64(t *testing.T) {
	data, err := testingAssetStore.GetAsBase64(testFileDataContentHash)
	if err != nil {
		t.Fail()
		t.Logf("Failed to get previously stored data file content: %v", err)
	} else if data != testFileDataContentB64 {
		t.Fail()
		t.Logf("Data received does not equal stored data: Expected: %v Got: %v", testFileDataContentB64, data)
	}
}

func TestAssetStore_StoreEmpty(t *testing.T) {
	dataHash := emptyTestFileDataContentHash
	hashResult, err := testingAssetStore.Store(emptyTestFileDataContentB64)
	if hashResult != dataHash {
		t.Fail()
		t.Logf("Expected hash: %v Got: %v", dataHash, hashResult)
	} else if err != nil {
		t.Fail()
		t.Logf("Store failed: %v", err)
	}
	path := testingAssetStore.makePath(dataHash) + ".snappy"
	if _, err := os.Stat(path); err != nil {
		t.Fail()
		t.Logf("Store did not store data: %v => %v", path, err)
	}
}

func TestAssetStore_GetAsBase64Empty(t *testing.T) {
	data, err := testingAssetStore.GetAsBase64(emptyTestFileDataContentHash)
	if err != nil {
		t.Fail()
		t.Logf("Failed to get previously stored data file content: %v", err)
	} else if data != emptyTestFileDataContentB64 {
		t.Fail()
		t.Logf("Data received does not equal stored data: Expected: %v Got: %v", emptyTestFileDataContentB64, data)
	}
}

func TestAssetStore_LoadEmpty(t *testing.T) {
	readCloser, err := testingAssetStore.Load(emptyTestFileDataContentHash)
	if err != nil {
		t.Fail()
		t.Logf("Failed to load file from previous test: %v", err)
	} else if readCloser == nil {
		t.Fail()
		t.Logf("Failed to get read closer without error")
	} else {
		receivedData, err := ioutil.ReadAll(readCloser)
		if err != nil {
			t.Logf("Failure on reading received data stream: %v", err)
			t.Fail()
		} else if string(receivedData) != emptyTestFileDataContent {
			t.Fail()
			t.Logf("Received data is not what was written to the file: Expected: %v Got: %v", emptyTestFileDataContent, string(receivedData))
		}
	}
	if readCloser != nil {
		readCloser.Close()
	}
}
