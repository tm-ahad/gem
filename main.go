package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var zmap map[string]string = make(map[string]string)
var zippedb bytes.Buffer

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func zip(path string, dir fs.DirEntry, e error) error {
	p := strings.Join(strings.Split(path, "/")[1:], "/")

	if !dir.IsDir() {
		cont, err := os.ReadFile(path)
		handleErr(err)

		zmap[p] = string(cont)
	} else {
		return nil
	}

	return nil
}

func unzip(path string, dst string) {
	var umap map[string]string

	os.Mkdir(dst, 0755)
	zipped_f, err := os.Open(path)

	handleErr(err)

	gr, err := gzip.NewReader(zipped_f)
	handleErr(err)

	gr.Close()
	handleErr(err)

	unzipped, err := io.ReadAll(gr)

	handleErr(err)
	handleErr(yaml.Unmarshal(unzipped, &umap));

	i := 0

	for p, cont := range umap {
		path := fmt.Sprintf("%s/%s", dst, p)

		spl_path := strings.Split(path, "/")
		spl_path = spl_path[:len(spl_path)-1]

		os.MkdirAll(strings.Join(spl_path, "/"), 0755)

		file, err := os.Create(path)
		handleErr(err)

		file.Write([]byte(cont))
		file.Close()

		i += 2
	}
}

func main() {
	com := os.Args[1]
	path := os.Args[2]

	switch com {
	case "zip":
		handleErr(filepath.WalkDir(path, zip))
		file, err := os.Create(fmt.Sprintf("%s.gf.zip", path))

		gw := gzip.NewWriter(&zippedb)

		handleErr(err)

		b, err := yaml.Marshal(zmap)
		handleErr(err)

		gw.Write(b)
		gw.Close()

		file.Write(zippedb.Bytes())
	case "unzip":
		unzip(path, os.Args[3])
	}
}
