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

var mime2asset = map[string]int8{
	"image/jp2":                           0,
	"application/ogg":                     1,
	"application/x-metaverse-callingcard": 2,
	"application/x-metaverse-landmark":    3,
	"application/x-metaverse-clothing":    5,
	"application/x-metaverse-primitive":   6,
	"application/x-metaverse-notecard":    7,
	"application/x-metaverse-folder":      8,
	"application/x-metaverse-lsl":         10,
	"application/x-metaverse-lso":         11,
	"image/tga":                           12,
	"application/x-metaverse-bodypart":    13,
	"audio/x-wav":                         17,
	"image/jpeg":                          19,
	"application/x-metaverse-animation":   20,
	"application/x-metaverse-gesture":     21,
	"application/x-metaverse-simstate":    22,
}
var asset2mime = map[int8]string{
	0:  "image/jp2",
	1:  "application/ogg",
	2:  "application/x-metaverse-callingcard",
	3:  "application/x-metaverse-landmark",
	5:  "application/x-metaverse-clothing",
	6:  "application/x-metaverse-primitive",
	7:  "application/x-metaverse-notecard",
	8:  "application/x-metaverse-folder",
	10: "application/x-metaverse-lsl",
	11: "application/x-metaverse-lso",
	12: "image/tga",
	13: "application/x-metaverse-bodypart",
	17: "audio/x-wav",
	19: "image/jpeg",
	20: "application/x-metaverse-animation",
	21: "application/x-metaverse-gesture",
	22: "application/x-metaverse-simstate",
}

func Mime2Asset(mime string) int8 {
	t, ok := mime2asset[mime]
	if !ok {
		return -1
	}
	return t
}

func Asset2Mime(asset int8) string {
	t, ok := asset2mime[asset]
	if !ok {
		return "application/octet-stream"
	}
	return t
}
