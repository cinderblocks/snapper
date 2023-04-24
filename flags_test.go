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
	"testing"
)

func TestAssetFlags_FromString(t *testing.T) {
	if AssetFlagsFromString("Normal,ReWrItabLe,maptile") != Normal|Maptile|Rewritable {
		t.Fail()
		t.Logf("Mismatch on expected output for: 'Normal,ReWrItabLe,maptile' Expected: %d Got: %d", Normal|Maptile|Rewritable, AssetFlagsFromString("Normal,ReWrItabLe,maptile"))
	}
	if AssetFlagsFromString("Normal,collectable,maptile") != Normal|Collectable|Maptile {
		t.Fail()
		t.Logf("Mismatch on expected output for: 'Normal,collectable,maptile' Expected: %d Got: %d", Normal|Collectable|Maptile, AssetFlagsFromString("Normal,collectable,maptile"))
	}
	if AssetFlagsFromString("Collectable") != Collectable {
		t.Fail()
		t.Logf("Mismatch on expected output for: 'Collectable' Expected: %d Got: %d", Collectable, AssetFlagsFromString("Collectable"))
	}
	if AssetFlagsFromString("Maptile") != Maptile {
		t.Fail()
		t.Logf("Mismatch on expected output for: 'Maptile' Expected: %d Got: %d", Maptile, AssetFlagsFromString("Maptile"))
	}
	if AssetFlagsFromString("Rewritable") != Rewritable {
		t.Fail()
		t.Logf("Mismatch on expected output for: 'Rewritable' Expected: %d Got: %d", Rewritable, AssetFlagsFromString("Rewritable"))
	}
}

func TestAssetFlags_ToString(t *testing.T) {
	if AssetFlagsToString(Normal|Maptile|Rewritable) != "Maptile,Rewritable" {
		t.Fail()
		t.Logf("Mismatch on expected output for: 'Normal|Maptile|Rewritable' Expected: %s Got: %s", "Maptile,Rewritable", AssetFlagsToString(Normal|Rewritable|Maptile))
	}
	if AssetFlagsToString(Normal|Collectable|Rewritable) != "Rewritable,Collectable" {
		t.Fail()
		t.Logf("Mismatch on expected output for: 'Normal|Maptile|Rewritable' Expected: %s Got: %s", "Rewritable,Collectable", AssetFlagsToString(Collectable|Maptile))
	}
}

func TestAssetFlags_ToStringCircle(t *testing.T) {
	pingPong := func(flags int64) bool {
		return AssetFlagsFromString(AssetFlagsToString(AssetFlagsFromString(AssetFlagsToString(flags)))) == flags
	}
	if !(pingPong(Maptile|Collectable) && pingPong(Collectable) && pingPong(Rewritable|Maptile) && pingPong(Normal) && pingPong(Normal|Collectable)) {
		t.Fail()
		t.Log("Some ping pong call failed")
	}
}
