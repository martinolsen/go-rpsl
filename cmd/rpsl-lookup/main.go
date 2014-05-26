package main

import (
	"flag"
	"io"
	"compress/gzip"
	"os"

	"github.com/martinolsen/go-rpsl"
)

func main() {
	db := flag.String("db", "", "RPSL database")

	flag.Parse()

	if *db == "" {
		flag.Usage()
		os.Exit(1)
	}

	file, err := os.Open(*db)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var reader io.Reader = file

	if r, err := gzip.NewReader(file); err == nil {
		reader = r
	}

	for _, object := range rpsl.Lookup(rpsl.NewReader(reader), flag.Arg(0)) {
		print(object.String())
	}
}
