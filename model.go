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
)

type AssetModel interface {
	Get(id string) (asset AssetBase, err error)
	GetHash(id string) (hash string, err error)
	GetHashAndType(id string) (hash string, assetType int8, err error)
	Put(asset AssetBase) error
}

type Database interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type assetModel struct {
	db Database
}

func CreateAssetModel(db Database) AssetModel {
	return &assetModel{
		db: db,
	}
}

func (a *assetModel) Get(id string) (asset AssetBase, err error) {
	err = a.db.Get(&asset, "SELECT * FROM `fsassets` WHERE `id` = ? LIMIT 1", id)
	asset.Flags = AssetFlagsToString(asset.DBFlags)
	asset.FullId = asset.Id
	return
}

func (a *assetModel) GetHash(id string) (hash string, err error) {
	err = a.db.Get(&hash, "SELECT `hash` FROM `fsassets` WHERE `id` = ? LIMIT 1", id)
	return
}

func (a *assetModel) GetHashAndType(id string) (hash string, assetType int8, err error) {
	asset, err := a.Get(id)
	if err == nil {
		hash = asset.Hash
		assetType = asset.Type
	}
	return
}

func (a *assetModel) Put(asset AssetBase) error {
	asset.DBFlags = AssetFlagsFromString(asset.Flags)
	_, err := a.db.Exec("INSERT INTO `fsassets` (id, type, hash, name, description, asset_flags, create_time, access_time)"+
		"VALUES(?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(NOW())) ON DUPLICATE KEY UPDATE type = ?, hash = ?, name = ?, description = ?, access_time = UNIX_TIMESTAMP(NOW()), asset_flags = ?",
		asset.Id, asset.Type, asset.Hash, asset.Name, asset.Description, asset.DBFlags,
		asset.Type, asset.Hash, asset.Name, asset.Description, asset.DBFlags)
	return err
}
