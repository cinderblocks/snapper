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
	"compress/gzip"
	"compress/zlib"
	"encoding/xml"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type HTTPService struct {
	service Service
	router  *mux.Router
}

func CreateHTTPService(db Database, dataStore, spoolStore string) *HTTPService {
	return &HTTPService{
		service: CreateService(db, dataStore, spoolStore),
	}
}

type responseWriter struct {
	http.ResponseWriter
}
type compressResponseWriter struct {
	io.Writer
	responseWriter
}

func (w compressResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (h HTTPService) xmlResponse(responseData interface{}, resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/xml")
	resp.Write([]byte(xml.Header))
	err := xml.NewEncoder(resp).Encode(responseData)
	if err != nil {
		http.NotFound(resp, req)
		log.Printf("Failed to encode response: %v\n", err)
	}
}

// Export this bit is the only thing worth using really
func Compressor(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			fn(compressResponseWriter{gz, responseWriter{w}}, r)
		} else if strings.Contains(r.Header.Get("Accept-Encoding"), "deflate") {
			w.Header().Set("Content-Encoding", "deflate")
			zl := zlib.NewWriter(w)
			defer zl.Close()
			fn(compressResponseWriter{zl, responseWriter{w}}, r)
		} else {
			fn(w, r)
		}
	}
}

func (h HTTPService) Router() *mux.Router {
	if h.router != nil {
		return h.router
	}
	router := mux.NewRouter()

	router.HandleFunc("/", h.index).Methods("GET")
	router.HandleFunc("/assets", h.index).Methods("GET")
	router.HandleFunc("/assets", h.create).Methods("POST")
	router.HandleFunc("/assets/{asset_id}/metadata", h.getMetadata).Methods("GET")
	router.HandleFunc("/assets/{asset_id}/data", h.getData).Methods("GET")
	router.HandleFunc("/assets/{asset_id}", h.get).Methods("GET")
	router.HandleFunc("/assets/{asset_id}", h.del).Methods("DELETE")
	router.HandleFunc("/get_assets_exist", h.exists).Methods("POST")
	h.router = router
	return router
}

func (h HTTPService) Run(listener net.Listener) {
	http.Serve(listener, handlers.LoggingHandler(os.Stdout, h.Router()))
}

func (h HTTPService) index(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/html")
	resp.Write([]byte("<html><head><title>OpenSimulator Assets Server</title></head><body><h1>OpenSimulator Assets Server</h1></body></html>"))
}

func (h HTTPService) del(resp http.ResponseWriter, req *http.Request) {
	return
}

func (h HTTPService) get(resp http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["asset_id"]
	full, err := h.service.GetFullAssetData(id)
	if err != nil {
		http.NotFound(resp, req)
		log.Printf("Failed to get asset data for %v Err: %v", id, err)
		return
	}
	h.xmlResponse(full, resp, req)
}

func (h HTTPService) exists(resp http.ResponseWriter, req *http.Request) {
	var ids = ArrayOfStrings{}
	defer req.Body.Close()

	err := xml.NewDecoder(req.Body).Decode(&ids)
	if err != nil {
		http.NotFound(resp, req)
		log.Printf("Failed to decode request: %v\n", err)
		return
	}

	var bools = ArrayOfBoolean{}
	bools.Booleans = h.service.AssetsExist(ids.Strings)

	h.xmlResponse(bools, resp, req)
}

func (h HTTPService) getData(resp http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["asset_id"]
	reader, assetType, err := h.service.GetAssetData(id)
	if err != nil {
		http.NotFound(resp, req)
		return
	} else {
		resp.Header().Set("Content-Type", Asset2Mime(assetType))
		defer reader.Close()
		_, err = io.Copy(resp, reader)
		if err != nil {
			http.NotFound(resp, req)
			log.Printf("Copying asset data to output stream failed: %v\n", err)
		}
	}
}

func (h HTTPService) getMetadata(resp http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["asset_id"]
	meta, err := h.service.GetAssetMetaData(id)
	if err != nil {
		http.NotFound(resp, req)
		log.Printf("Failed to get asset meta data for %v Err: %v\n", id, err)
		return
	}
	h.xmlResponse(meta, resp, req)
}

func (h HTTPService) create(resp http.ResponseWriter, req *http.Request) {
	var fullData FullAssetData
	defer req.Body.Close()
	err := xml.NewDecoder(req.Body).Decode(&fullData)
	if err != nil {
		http.NotFound(resp, req)
		log.Printf("AssetCreate: Failed to decode data: %v\n", err)
		return
	}
	err = h.service.CreateAsset(&fullData)
	if err != nil {
		http.NotFound(resp, req)
		log.Printf("Failed to create asset: %v Error: %v\n", fullData.Id, err)
	} else {
		log.Printf("Successfully created asset: %v Hash: %v\n", fullData.Id, fullData.Hash)
		h.xmlResponse(CreateResponseSuccess{Id: fullData.Id}, resp, req)
	}
}
