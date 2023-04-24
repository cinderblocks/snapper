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
	"database/sql"
	"strings"
	"testing"
)

func testModelAssetInstance(forGet bool) AssetBase {
	dbFlags := int64(0)
	flags := "Rewritable|Collectible"
	if forGet {
		dbFlags = AssetFlagsFromString(flags)
		flags = ""
	}
	return AssetBase{
		FullId:      testContentId,
		Id:          testContentId,
		Name:        "TestAsset",
		Description: "Test Asset Description",
		Flags:       flags,
		Type:        7,
		DBFlags:     dbFlags,
		Hash:        testFileDataContentHash,
	}
}

type mockDatabase struct {
	t    *testing.T
	data AssetBase
}

func (m *mockDatabase) Get(dest interface{}, query string, args ...interface{}) error {
	if len(args) != 1 {
		m.t.Fail()
		m.t.Logf("Expected 1 argument got: %v", len(args))
	} else if m.data.Id != args[0].(string) {
		m.t.Fail()
		m.t.Logf("Expected Id: '%s' as argument got: %v", args[0])
	}
	switch dest.(type) {
	case *AssetBase:
		*dest.(*AssetBase) = m.data
	case *string:
		*dest.(*string) = m.data.Hash
	}

	return nil
}

func (m *mockDatabase) Exec(query string, args ...interface{}) (sql.Result, error) {
	placeholders := strings.Count(query, "?")
	if len(args) != placeholders {
		m.t.Fail()
		m.t.Logf("Upsert requires %d parameters. Got: %d arguments", placeholders, len(args))
	} else if placeholders != 11 {
		m.t.Fail()
		m.t.Log("This test was designed for 11 parameters and needs an update!")
	} else {
		expectedFlags := AssetFlagsFromString(m.data.Flags)
		if expectedFlags != args[5].(int64) || expectedFlags != args[10] {
			m.t.Fail()
			m.t.Logf("Expected flags to be converted")
		}
	}
	var ok bool = false
	m.data.Id, ok = args[0].(string)
	if !ok {
		m.t.Fail()
		m.t.Log("Expected first value to be a string with the id")
	}
	m.data.Type, ok = args[1].(int8)
	if !ok {
		m.t.Fail()
		m.t.Log("Expected second value to be an int8 with the asset type")
	}
	m.data.Hash, ok = args[2].(string)
	if !ok {
		m.t.Fail()
		m.t.Log("Expected third value to be a string with the hash")
	}
	m.data.Name, ok = args[3].(string)
	if !ok {
		m.t.Fail()
		m.t.Log("Expected fourth value to be a string with the name")
	}
	m.data.Description, ok = args[4].(string)
	if !ok {
		m.t.Fail()
		m.t.Log("Expected fifth value to be a string with the description")
	}
	m.data.DBFlags, ok = args[5].(int64)
	if !ok {
		m.t.Fail()
		m.t.Log("Expected sixth value to be a int64 with the asset flags")
	}

	if t, ok := args[6].(int8); !ok || t != m.data.Type {
		m.t.Fail()
		m.t.Log("Asset type or value mismatch on seventh parameter on update")
	}
	if t, ok := args[7].(string); !ok || t != m.data.Hash {
		m.t.Fail()
		m.t.Log("Asset type or value mismatch on eigth parameter on update")
	}
	if t, ok := args[8].(string); !ok || t != m.data.Name {
		m.t.Fail()
		m.t.Log("Asset type or value mismatch on ninth parameter on update")
	}
	if t, ok := args[9].(string); !ok || t != m.data.Description {
		m.t.Fail()
		m.t.Log("Asset type or value mismatch on tenth parameter on update")
	}
	if t, ok := args[10].(int64); !ok || t != m.data.DBFlags {
		m.t.Fail()
		m.t.Log("Asset type or value mismatch on eleventh parameter on update")
	}

	return nil, nil
}

func TestAssetModel_Get(t *testing.T) {
	m := &mockDatabase{t: t, data: testModelAssetInstance(true)}
	a, _ := CreateAssetModel(m).Get(m.data.Id)
	if a.FullId != m.data.Id {
		t.Fail()
		t.Logf("Failure: Model must set FullId -> %s Got: %s", m.data.Id, a.FullId)
	}
	if a.Flags != AssetFlagsToString(m.data.DBFlags) {
		t.Fail()
		t.Logf("Failure: Model must convert dbflags to string flags")
	}
}

func TestAssetModel_GetHash(t *testing.T) {
	m := &mockDatabase{t: t, data: testModelAssetInstance(true)}
	hash, _ := CreateAssetModel(m).GetHash(m.data.Id)
	if hash != m.data.Hash {
		t.Fail()
		t.Logf("Failure: expected returned hash: %s Got: %s", m.data.Hash, hash)
	}
}

func TestAssetModel_GetHashAndType(t *testing.T) {
	m := &mockDatabase{t: t, data: testModelAssetInstance(true)}
	hash, assetType, _ := CreateAssetModel(m).GetHashAndType(m.data.Id)
	if hash != m.data.Hash {
		t.Fail()
		t.Logf("Failure: expected returned hash: %s Got: %s", m.data.Hash, hash)
	}
	if assetType != m.data.Type {
		t.Fail()
		t.Logf("Failure: expected returned typed: %d Got: %d", m.data.Type, assetType)
	}
}

func TestAssetModel_Put(t *testing.T) {
	m := &mockDatabase{t: t, data: AssetBase{}}
	data := testModelAssetInstance(false)
	CreateAssetModel(m).Put(data)
	if m.data.Id != data.Id {
		t.Fail()
		t.Logf("Expected Id to be: %s Got: %s", m.data.Id, data.Id)
	}
	if m.data.Type != data.Type {
		t.Fail()
		t.Logf("Expected Type to be: %d Got: %d", m.data.Type, data.Type)
	}
	if m.data.Hash != data.Hash {
		t.Fail()
		t.Logf("Expected Hash to be: %s Got: %s", m.data.Hash, data.Hash)
	}
	if m.data.Name != data.Name {
		t.Fail()
		t.Logf("Expected Name to be: %s Got: %s", m.data.Name, data.Name)
	}
	if m.data.Description != data.Description {
		t.Fail()
		t.Logf("Expected Name to be: %s Got: %s", m.data.Name, data.Name)
	}
	if m.data.DBFlags != data.DBFlags {
		t.Fail()
		t.Logf("Expected DBFlags to be: %d Got: %d", m.data.DBFlags, data.DBFlags)
	}
}
