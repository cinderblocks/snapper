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
	"strings"
)

const (
	Normal      = 0
	Maptile     = 1
	Rewritable  = 2
	Collectable = 4
)

func AssetFlagsFromString(flags string) int64 {
	parts := strings.Split(strings.ToLower(flags), ",")
	var result int64 = Normal
	for _, p := range parts {
		switch strings.Trim(p, " \t\r\n") {
		case "maptile":
			result |= Maptile
		case "rewritable":
			result |= Rewritable
		case "collectable":
			result |= Collectable
		default:
			//
		}
	}
	return result
}

func AssetFlagsToString(flags int64) string {
	result := []string{}
	if (flags & Maptile) == Maptile {
		result = append(result, "Maptile")
	}
	if (flags & Rewritable) == Rewritable {
		result = append(result, "Rewritable")
	}
	if (flags & Collectable) == Collectable {
		result = append(result, "Collectable")
	}
	return strings.Join(result, ",")
}
