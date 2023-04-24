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
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	db, err := sqlx.Open("mysql", os.Getenv("ASSETSDBCON"))
	if err != nil {
		fmt.Errorf("ERROR: Unable to establish database connection: %v\n", err.Error())
	}

	var dataStore = flag.String("datastore", "asset/data", "Path to asset data store")
	var spoolStore = flag.String("spoolstore", "asset/tmp", "Path to asset temporary data store")
	var address = flag.String("address", "0.0.0.0:8003", "Address to listen to. Default: 0.0.0.0:8003")
	flag.Parse()

	listener, err := net.Listen("tcp", *address)
	if err != nil {
		fmt.Errorf("Failed to listen to specified address: %v ERROR: %v\n", address, err)
	}

	httpService := CreateHTTPService(db, *dataStore, *spoolStore)
	httpService.Run(listener)
}
