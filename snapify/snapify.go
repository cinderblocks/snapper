package main

import (
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/golang/snappy"
)

type context struct {
	toProcess []string
}

func (c *context) visit(p string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		if path.Ext(p) != ".snappy" {
			c.toProcess = append(c.toProcess, p)
		}
	}
	return nil
}

func hash(p string) ([]byte, error) {
	f, e := os.Open(p)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	hasher := sha256.New()
	var reader io.Reader = nil
	if path.Ext(p) == ".gz" {
		gz, e := gzip.NewReader(f)
		if e != nil {
			return nil, e
		}
		defer gz.Close()
		reader = gz
	} else if path.Ext(p) == ".snappy" {
		reader = snappy.NewReader(f)
	} else {
		reader = f
	}

	_, e = io.Copy(hasher, reader)
	if e != nil {
		return nil, e
	}
	return hasher.Sum(nil), nil
}

func validate(old, new string) error {
	oldSum, e := hash(old)
	if e != nil {
		return e
	}
	newSum, e := hash(new)
	if e != nil {
		return e
	}
	if len(oldSum) != len(newSum) {
		return os.ErrInvalid
	}
	for i, _ := range oldSum {
		if oldSum[i] != newSum[i] {
			return os.ErrInvalid
		}
	}
	return nil
}

func convert(p string) (string, error) {
	newpath := p + ".snappy"
	isGzip := false
	if path.Ext(p) == ".gz" {
		isGzip = true
		newpath = p[0:len(p)-3] + ".snappy"
	}
	inFile, e := os.Open(p)
	if e != nil {
		fmt.Print("Failed to open input file: %s Err: %v\n", p, e)
		return "", e
	}
	defer inFile.Close()
	outFile, e := os.Create(newpath)
	if e != nil {
		if os.IsExist(e) {
			return newpath, nil
		}
		fmt.Printf("Failed to create output file: %s Err: %v\n", newpath, e)
		return "", e
	}
	defer outFile.Close()
	writer := snappy.NewWriter(outFile)
	if isGzip {
		g, e := gzip.NewReader(inFile)
		if e != nil {
			fmt.Printf("Failed to open gzipped file for reading: %s Err: %v\n", p, e)
			return "", e
		}
		defer g.Close()
		io.Copy(writer, g)
	} else {
		io.Copy(writer, inFile)
	}
	return outFile.Name(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s path-to-assets-data-storage\n", os.Args[0])
		os.Exit(1)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())

	c := context{
		toProcess: []string{},
	}
	filepath.Walk(os.Args[1], c.visit)

	count := len(c.toProcess)
	for i, p := range c.toProcess {
		fmt.Printf("\rProcessing file %d out of %d", i+1, count)
		n, e := convert(p)
		if e != nil {
			fmt.Printf("Failed to convert: %v\n", e)
		} else {
			e = validate(p, n)
			if e != nil {
				fmt.Printf("Validation failed: %v\n", e)
			}
		}
	}
}
