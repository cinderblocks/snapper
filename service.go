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
	"io"
)

type service struct {
	model AssetModel
	store AssetStore
}

type Service interface {
	GetFullAssetData(id string) (data FullAssetData, err error)
	GetAssetMetaData(id string) (data AssetBase, err error)
	GetAssetData(id string) (io.ReadCloser, int8, error)
	CreateAsset(data *FullAssetData) error
	AssetExists(id string) bool
	AssetsExist(ids []string) []bool
}

func CreateService(db Database, dataDir, spoolDir string) Service {
	return &service{
		model: CreateAssetModel(db),
		store: CreateAssetStore(dataDir, spoolDir),
	}
}

func (s service) GetFullAssetData(id string) (data FullAssetData, err error) {
	data.AssetBase, err = s.model.Get(id)
	if err == nil {
		data.Data, err = s.store.GetAsBase64(data.Hash)
	}
	return
}

func (s service) GetAssetMetaData(id string) (data AssetBase, err error) {
	data, err = s.model.Get(id)
	return
}

func (s service) GetAssetData(id string) (io.ReadCloser, int8, error) {
	hash, assetType, err := s.model.GetHashAndType(id)
	if err != nil {
		return nil, 0, err
	}
	reader, err := s.store.Load(hash)
	return reader, assetType, err
}

func (s service) CreateAsset(data *FullAssetData) error {
	var err error = nil
	data.Hash, err = s.store.Store(data.Data)
	if err != nil {
		return err
	}
	return s.model.Put(data.AssetBase)
}

func (s service) AssetExists(id string) bool {
	hash, err := s.model.GetHash(id)
	if err != nil {
		return false
	}
	return s.store.Exists(hash)
}

func (s service) AssetsExist(ids []string) []bool {
	result := make([]bool, len(ids))
	for idx, id := range ids {
		result[idx] = s.AssetExists(id)
	}
	return result
}
