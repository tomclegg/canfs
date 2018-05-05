// SPDX-License-Identifier: CC0-1.0

// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tomclegg/canfs"
)

func main() {
	id := flag.String("id", "canfs", "`identifier` of the generated http.Filesystem")
	out := flag.String("out", "canfs_generated.go", "write generated code to `filename`")
	dir := flag.String("dir", "canfs_data", "use local `directory` as filesystem root")
	pkg := flag.String("pkg", "main", "package `name`")
	flag.Parse()
	err := generate(*id, *out, *dir, *pkg)
	if err != nil {
		log.Fatal(err)
	}
}

func generate(id, out, dir, pkg string) error {
	outFile, err := os.OpenFile(out+".tmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		os.Remove(out + ".tmp")
	}()

	content := map[string]canfs.FileInfo{}
	root := strings.TrimSuffix(path.Clean(dir), "/")
	err = filepath.Walk(root, builder{strip: root, content: content}.walkFunc)
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "package %s\n", pkg)
	fmt.Fprintf(buf, "import \"github.com/tomclegg/canfs\"\n")
	fmt.Fprintf(buf, "var %s = canfs.FileSystem{Content: map[string]canfs.FileInfo{\n", id)
	var paths []string
	for path := range content {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Fprintf(buf, "\t%q: %#v,\n", path, content[path])
	}
	fmt.Fprintf(buf, "}}\n")

	gofmt := exec.Command("gofmt", "-s")
	gofmt.Stdin = buf
	gofmt.Stdout = outFile
	gofmt.Stderr = os.Stderr
	err = gofmt.Run()
	if err != nil {
		return err
	}

	err = outFile.Close()
	if err != nil {
		return err
	}
	return os.Rename(out+".tmp", out)
}

type builder struct {
	strip   string
	content map[string]canfs.FileInfo
}

func (bldr builder) walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil || info.IsDir() {
		return err
	}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	path = strings.TrimPrefix(path, bldr.strip)

	var data canfs.FileData
	if s := string(buf); bytes.Compare([]byte(s), buf) == 0 {
		data = canfs.StringData{s}
	} else {
		data = canfs.ByteData{buf}
	}
	bldr.content[path] = canfs.FileInfo{
		N:        info.Name(),
		M:        info.Mode(),
		S:        info.Size(),
		MT:       info.ModTime().UnixNano(),
		FileData: data,
	}
	return nil
}
