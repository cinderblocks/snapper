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
	"io"
	"io/ioutil"
	"os"
	"testing"
)

const (
	testContentId = "83d4692b-fc3e-4b8f-a331-da07b2166250"
)

func testServiceAssetInstance() AssetBase {
	return AssetBase{
		FullId:      testContentId,
		Id:          testContentId,
		Name:        "TestAsset",
		Description: "Test Asset Description",
		Flags:       "Maptile,Collectable",
		DBFlags:     Maptile | Collectable,
		Hash:        testFileDataContentHash,
		Type:        4,
	}

}

type mockModel struct {
	GetCalls, GetHashCalls, GetHashAndTypeCalls, PutCalls int
}

func (m *mockModel) Get(id string) (asset AssetBase, err error) {
	m.GetCalls++
	return testServiceAssetInstance(), nil
}

func (m *mockModel) GetHash(id string) (hash string, err error) {
	m.GetHashCalls++
	return testServiceAssetInstance().Hash, nil
}

func (m *mockModel) GetHashAndType(id string) (hash string, assetType int8, err error) {
	m.GetHashAndTypeCalls++
	inst := testServiceAssetInstance()
	return inst.Hash, inst.Type, nil
}

func (m *mockModel) Put(asset AssetBase) error {
	m.PutCalls++
	inst := testServiceAssetInstance()
	if asset.Type == inst.Type &&
		asset.Name == inst.Name &&
		asset.Id == inst.Id &&
		asset.Description == inst.Description {

	}
	return nil
}

type mockStore struct {
	testData     string
	testDataB64  string
	expectedHash string
}

type mockDataSource struct {
	reader io.Reader
}

func (m *mockDataSource) Read(b []byte) (int, error) {
	return m.reader.Read(b)
}

func (m *mockDataSource) Close() error {
	return nil
}

func (m *mockStore) Load(hash string) (io.ReadCloser, error) {
	if hash == m.expectedHash {
		return &mockDataSource{bytes.NewReader([]byte(m.testData))}, nil
	}
	return nil, os.ErrNotExist
}

func (m *mockStore) Exists(hash string) bool {
	return m.expectedHash == hash
}

func (m *mockStore) Store(data string) (string, error) {
	if data == m.testData {
		return m.expectedHash, nil
	}
	return "", os.ErrInvalid
}

func (m *mockStore) GetAsBase64(hash string) (string, error) {
	if m.expectedHash == hash {
		return m.testDataB64, nil
	}
	return "", os.ErrNotExist
}

func validateMeta(t *testing.T, recv AssetBase) {
	expec := testServiceAssetInstance()
	if expec.Id != recv.Id ||
		expec.Name != recv.Name ||
		expec.Description != recv.Description ||
		expec.Flags != recv.Flags ||
		expec.Type != recv.Type ||
		expec.Hash != recv.Hash {

		t.Fail()
		t.Logf("Data mismatch: %v Got: %v ", expec, recv)
	}
}

func TestService_GetFullAssetData(t *testing.T) {
	svc := &service{
		model: &mockModel{},
		store: &mockStore{
			testData:     testFileDataContent,
			testDataB64:  testFileDataContentB64,
			expectedHash: testFileDataContentHash,
		},
	}

	fullData, err := svc.GetFullAssetData(testContentId)
	if err != nil {
		t.Fail()
		t.Logf("Unexpected error on GetFullAssetData: %v", err)
	}
	if fullData.Data != testFileDataContentB64 {
		t.Fail()
		t.Logf("Expected data: %s Got: %v", testFileDataContentB64, fullData.Data)
	}
	mmodel := svc.model.(*mockModel)
	if mmodel.GetCalls != 1 || mmodel.GetHashAndTypeCalls != 0 || mmodel.GetHashCalls != 0 || mmodel.PutCalls != 0 {
		t.Fail()
		t.Log("Expected one call on Model to Get(id)")
	}
	validateMeta(t, fullData.AssetBase)
}

func TestService_GetAssetMetaData(t *testing.T) {
	svc := &service{
		model: &mockModel{},
		store: &mockStore{
			testData:     testFileDataContent,
			testDataB64:  testFileDataContentB64,
			expectedHash: testFileDataContentHash,
		},
	}

	metaData, err := svc.GetAssetMetaData(testContentId)
	if err != nil {
		t.Fail()
		t.Logf("Unexpected error on GetAssetMetaData: %v", err)
	}

	mmodel := svc.model.(*mockModel)
	if mmodel.GetCalls != 1 || mmodel.GetHashAndTypeCalls != 0 || mmodel.GetHashCalls != 0 || mmodel.PutCalls != 0 {
		t.Fail()
		t.Log("Expected one call on Model to Get(id)")
	}
	validateMeta(t, metaData)
}

func TestService_GetAssetData(t *testing.T) {
	svc := &service{
		model: &mockModel{},
		store: &mockStore{
			testData:     testFileDataContent,
			testDataB64:  testFileDataContentB64,
			expectedHash: testFileDataContentHash,
		},
	}

	readerCloser, assetType, err := svc.GetAssetData(testContentId)
	if err != nil {
		t.Fail()
		t.Logf("Unexpected error on GetAssetData: %v", err)
	} else {
		defer readerCloser.Close()
	}
	data, err := ioutil.ReadAll(readerCloser)
	if err == nil {
		if string(data) != testFileDataContent {
			t.Fail()
			t.Logf("Data passed through corrupted: %v Got: %v", testFileDataContent, string(data))
		}
	}
	if assetType != testServiceAssetInstance().Type {
		t.Fail()
		t.Logf("Unexpected assetType returned: %v Got: %v", testServiceAssetInstance().Type, assetType)
	}

	mmodel := svc.model.(*mockModel)
	if mmodel.GetCalls != 0 || mmodel.GetHashAndTypeCalls != 1 || mmodel.GetHashCalls != 0 || mmodel.PutCalls != 0 {
		t.Fail()
		t.Log("Expected one call on Model to Get(id)")
	}
}

func TestService_CreateAsset(t *testing.T) {
	svc := &service{
		model: &mockModel{},
		store: &mockStore{
			testData:     testFileDataContent,
			testDataB64:  "",
			expectedHash: testFileDataContentHash,
		},
	}

	data := FullAssetData{}
	data.AssetBase = testServiceAssetInstance()
	data.Data = testFileDataContent
	err := svc.CreateAsset(&data)
	if err != nil {
		t.Fail()
		t.Logf("Unexpected error on CreateAsset: %v", err)
	}
	if data.Hash != testFileDataContentHash {
		t.Fail()
		t.Logf("Unexpected hash returned. Expected: %v Got: %v", testFileDataContentHash, data.Hash)
	}

	mmodel := svc.model.(*mockModel)
	if mmodel.GetCalls != 0 || mmodel.GetHashAndTypeCalls != 0 || mmodel.GetHashCalls != 0 || mmodel.PutCalls != 1 {
		t.Fail()
		t.Log("Expected one call on Model to Get(id)")
	}
}

func TestService_AssetExists(t *testing.T) {
	svc := &service{
		model: &mockModel{},
		store: &mockStore{
			testData:     "",
			testDataB64:  "",
			expectedHash: testFileDataContentHash,
		},
	}
	if !svc.AssetExists(testContentId) {
		t.Fail()
		t.Log("Expected asset to exist!")
	}

	mmodel := svc.model.(*mockModel)
	if mmodel.GetCalls != 0 || mmodel.GetHashAndTypeCalls != 0 || mmodel.GetHashCalls != 1 || mmodel.PutCalls != 0 {
		t.Fail()
		t.Log("Expected one call on Model to Get(id)")
	}
}

func TestService_AssetExistsNegative(t *testing.T) {
	svc := &service{
		model: &mockModel{},
		store: &mockStore{
			testData:     "",
			testDataB64:  "",
			expectedHash: "",
		},
	}
	if svc.AssetExists(testContentId) {
		t.Fail()
		t.Log("Expected asset to NOT exist!")
	}

	mmodel := svc.model.(*mockModel)
	if mmodel.GetCalls != 0 || mmodel.GetHashAndTypeCalls != 0 || mmodel.GetHashCalls != 1 || mmodel.PutCalls != 0 {
		t.Fail()
		t.Log("Expected one call on Model to Get(id)")
	}
}

func TestService_AssetsExists(t *testing.T) {
	svc := &service{
		model: &mockModel{},
		store: &mockStore{
			testData:     "",
			testDataB64:  "",
			expectedHash: testFileDataContentHash,
		},
	}
	result := svc.AssetsExist([]string{testContentId})
	if len(result) == 1 && !result[0] {
		t.Fail()
		t.Log("Expected asset to exist!")
	}

	mmodel := svc.model.(*mockModel)
	if mmodel.GetCalls != 0 || mmodel.GetHashAndTypeCalls != 0 || mmodel.GetHashCalls != 1 || mmodel.PutCalls != 0 {
		t.Fail()
		t.Log("Expected one call on Model to Get(id)")
	}
}

func TestService_AssetsExistsNegative(t *testing.T) {
	svc := &service{
		model: &mockModel{},
		store: &mockStore{
			testData:     "",
			testDataB64:  "",
			expectedHash: "",
		},
	}
	result := svc.AssetsExist([]string{testContentId})
	if len(result) == 1 && result[0] {
		t.Fail()
		t.Log("Expected asset to NOT exist!")
	}

	mmodel := svc.model.(*mockModel)
	if mmodel.GetCalls != 0 || mmodel.GetHashAndTypeCalls != 0 || mmodel.GetHashCalls != 1 || mmodel.PutCalls != 0 {
		t.Fail()
		t.Log("Expected one call on Model to Get(id)")
	}
}
