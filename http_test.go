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
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type mockService struct {
}

func (m *mockService) GetFullAssetData(id string) (data FullAssetData, err error) {
	if id == testContentId {
		return FullAssetData{
			AssetBase: testServiceAssetInstance(),
			Data:      testFileDataContentB64,
		}, nil
	}
	return FullAssetData{}, os.ErrNotExist
}

func (m *mockService) GetAssetMetaData(id string) (data AssetBase, err error) {
	if id == testContentId {
		return testServiceAssetInstance(), nil
	}
	return AssetBase{}, os.ErrNotExist
}

func (m *mockService) GetAssetData(id string) (io.ReadCloser, int8, error) {
	if id == testContentId {
		return &mockDataSource{bytes.NewReader([]byte(testFileDataContent))}, testServiceAssetInstance().Type, nil
	}
	return nil, 0, os.ErrNotExist
}

func (m *mockService) CreateAsset(data *FullAssetData) error {
	if data.Id == testContentId {
		data.Hash = testFileDataContentHash
		return nil
	}
	return os.ErrInvalid
}

func (m *mockService) AssetExists(id string) bool {
	return id == testContentId
}

func (m *mockService) AssetsExist(ids []string) []bool {
	result := make([]bool, len(ids))
	for i := range ids {
		result[i] = m.AssetExists(ids[i])
	}
	return result
}

var httpTestServiceInstance *HTTPService = &HTTPService{service: &mockService{}}

func TestHTTP_GetFullData(t *testing.T) {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/assets/"+testContentId, nil)
	httpTestServiceInstance.Router().ServeHTTP(recorder, request)
	if recorder.Code != 200 {
		t.Fail()
		t.Logf("Expected non failure on request: Got Code: %v", recorder.Code)
	} else {
		if recorder.Header().Get("Content-Type") != "application/xml" {
			t.Fail()
			t.Logf("Expected Content-Type: application/xml Got: %v", recorder.Header().Get("Content-Type"))
		}
		result := FullAssetData{}
		err := xml.NewDecoder(recorder.Body).Decode(&result)
		if err != nil {
			t.Fail()
			t.Logf("Decoding response failed: %v", err)
		}
		if result.Id != testContentId ||
			result.Hash != "" ||
			result.Hash == testFileDataContentHash ||
			result.Type != testServiceAssetInstance().Type ||
			result.Flags != testServiceAssetInstance().Flags ||
			result.Data != testFileDataContentB64 {
			t.Fail()
			t.Logf("Data not correctly passed: Id: %v Hash: %v Type: %v Flags: %v Data: %v", result.Id, result.Hash, result.Type, result.Flags, result.Data)
		}
	}
}

func TestHTTP_GetFullDataNotExists(t *testing.T) {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/assets/-"+testContentId[1:len(testContentId)], nil)
	httpTestServiceInstance.Router().ServeHTTP(recorder, request)
	if recorder.Code != 404 {
		t.Fail()
		t.Logf("Expected failure on request: Got Code: %v", recorder.Code)
	}
}

func TestHTTP_GetMetaData(t *testing.T) {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/assets/"+testContentId+"/metadata", nil)
	httpTestServiceInstance.Router().ServeHTTP(recorder, request)
	if recorder.Code != 200 {
		t.Fail()
		t.Logf("Expected non failure on request: Got Code: %v", recorder.Code)
	} else {
		if recorder.Header().Get("Content-Type") != "application/xml" {
			t.Fail()
			t.Logf("Expected Content-Type: application/xml Got: %v", recorder.Header().Get("Content-Type"))
		}
		result := AssetBase{}
		err := xml.NewDecoder(recorder.Body).Decode(&result)
		if err != nil {
			t.Fail()
			t.Logf("Decoding response failed: %v", err)
		}
		if result.Id != testContentId ||
			result.Hash != "" ||
			result.Hash == testFileDataContentHash ||
			result.Type != testServiceAssetInstance().Type ||
			result.Flags != testServiceAssetInstance().Flags {
			t.Fail()
			t.Logf("Data not correctly passed: Id: %v Hash: %v Type: %v Flags: %v", result.Id, result.Hash, result.Type, result.Flags)
		}
	}
}

func TestHTTP_GetMetaDataNotExists(t *testing.T) {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/assets/-"+testContentId[1:len(testContentId)]+"/metadata", nil)
	httpTestServiceInstance.Router().ServeHTTP(recorder, request)
	if recorder.Code != 404 {
		t.Fail()
		t.Logf("Expected failure on request: Got Code: %v", recorder.Code)
	}
}

func TestHTTP_GetData(t *testing.T) {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/assets/"+testContentId+"/data", nil)
	httpTestServiceInstance.Router().ServeHTTP(recorder, request)
	if recorder.Code != 200 {
		t.Fail()
		t.Logf("Expected non failure on request: Got Code: %v", recorder.Code)
	} else {
		if recorder.Header().Get("Content-Type") != Asset2Mime(testServiceAssetInstance().Type) {
			t.Fail()
			t.Logf("Expected Content-Type: %s Got: %v", Asset2Mime(testServiceAssetInstance().Type), recorder.Header().Get("Content-Type"))
		}
		recv := string(recorder.Body.Bytes())
		if recv != testFileDataContent {
			t.Fail()
			t.Logf("Expected data: %v Got: %v", testFileDataContent, recv)
		}
	}
}

func TestHTTP_GetDataNotExists(t *testing.T) {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/assets/-"+testContentId[1:len(testContentId)]+"/data", nil)
	httpTestServiceInstance.Router().ServeHTTP(recorder, request)
	if recorder.Code != 404 {
		t.Fail()
		t.Logf("Expected failure on request: Got Code: %v", recorder.Code)
	}
}

func TestHTTP_GetAssetExists(t *testing.T) {
	recorder := httptest.NewRecorder()
	testData := ArrayOfStrings{
		Strings: []string{
			testContentId,
			testContentId,
			"invalid-id",
		},
	}
	data, err := xml.Marshal(testData)
	data = append([]byte(xml.Header), data...)
	if err != nil {
		t.Skip("TestHTTP_GetAssetExists test broken")
	} else {
		request, _ := http.NewRequest("POST", "/get_assets_exist", bytes.NewReader(data))
		httpTestServiceInstance.Router().ServeHTTP(recorder, request)
		if recorder.Code != 200 {
			t.Fail()
			t.Logf("Expected non failure on request: Got Code: %v", recorder.Code)
		} else {
			result := ArrayOfBoolean{}
			err := xml.NewDecoder(recorder.Body).Decode(&result)
			if err != nil {
				t.Fail()
				t.Logf("Failed to decode response: %v", err)
			} else {
				if len(result.Booleans) != len(testData.Strings) {
					t.Fail()
					t.Logf("Expected %d booleans got: %d", len(testData.Strings), len(result.Booleans))
				} else {
					if result.Booleans[0] != true ||
						result.Booleans[1] != true ||
						result.Booleans[2] != false {
						t.Fail()
						t.Logf("Wrong answer. Expected: true true false Got: %t %t %t", result.Booleans[0], result.Booleans[1], result.Booleans[2])
					}
				}
			}
		}

	}
}

func TestHTTP_Create(t *testing.T) {
	fullData, err := httpTestServiceInstance.service.GetFullAssetData(testContentId)
	data, err := xml.Marshal(&fullData)
	data = append([]byte(xml.Header), data...)
	if err != nil {
		t.Skip("TestHTTP_GetAssetExists test broken")
	} else {
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/assets", bytes.NewReader(data))
		httpTestServiceInstance.Router().ServeHTTP(recorder, request)
		if recorder.Code != 200 {
			t.Fail()
			t.Logf("Expected non failure on request: Got Code: %v", recorder.Code)
		} else {
			result := CreateResponseSuccess{}
			err := xml.NewDecoder(recorder.Body).Decode(&result)
			if err != nil {
				t.Fail()
				t.Logf("Failed to decode response: %v", err)
			} else {
				if result.Id != fullData.Id {
					t.Fail()
					t.Logf("Expected response ID: %v Got: %v", fullData.Id, result.Id)
				}
			}
		}
	}
}

func TestHTTP_Delete(t *testing.T) {
	// This is not really implemented, however expected to return 200 - No matter what
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/assets/"+testContentId, nil)
	httpTestServiceInstance.Router().ServeHTTP(recorder, request)
	if recorder.Code != 200 {
		t.Fail()
		t.Logf("Expected non failure on request: Got Code: %v", recorder.Code)
	}
}
