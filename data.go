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
	"encoding/xml"
)

type CreateResponseSuccess struct {
	Id string `xml:"string"`
}

type ArrayOfStrings struct {
	XMLName xml.Name `xml:"ArrayOfStrings"`
	Strings []string `xml:"string"`
}

type ArrayOfBoolean struct {
	XMLName  xml.Name `xml:"ArrayOfBoolean"`
	Booleans []bool   `xml:"boolean"`
}

type AssetBase struct {
	XMLName     xml.Name `xml:"AssetBase" db:"-"`
	FullId      string   `xml:"FullID>Guid,omitempty" db:"-"`
	Id          string   `xml:"ID" db:"id"`
	Name        string   `xml:"Name" db:"name"`
	Description string   `xml:"Description" db:"description"`
	Flags       string   `xml:"Flags" db:"-"`
	DBFlags     int64    `xml:"-" db:"asset_flags"`
	Type        int8     `xml:"Type" db:"type"`
	CreatorID   string   `xml:"CreatorID,omitempty" db:"-"`
	Temporary   bool     `xml:"Temporary,omitempty" db:"-"`
	Local       bool     `xml:"Local,omitempty" db:"-"`
	AccessTime  int64    `xml:"-" db:"create_time"`
	CreateTime  int64    `xml:"-" db:"access_time"`
	Hash        string   `xml:"-" db:"hash"`
}

type FullAssetData struct {
	AssetBase
	Data string `xml:"Data,omitempty" db:"-"`
}
