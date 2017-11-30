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
	"regexp"
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
	outFile, err := os.OpenFile(out+".tmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func() {
		os.Remove(out + ".tmp")
	}()

	content := map[string]canfs.FileData{}
	root := strings.TrimSuffix(path.Clean(dir), "/")
	err = filepath.Walk(root, builder{strip: root, content: content}.walkFunc)
	if err != nil {
		return err
	}
	data := canfs.FileSystem{Content: content}
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "package %s\n", pkg)
	if len(content) > 0 {
		fmt.Fprintf(buf, "import \"time\"\n")
	}
	fmt.Fprintf(buf, "import \"github.com/tomclegg/canfs\"\n")
	fmt.Fprintf(buf, "var %s = %#v\n", id, data)

	re, err := regexp.Compile(`time.Time{sec:(\d+), nsec:(\d+).*?}`)
	if err != nil {
		return err
	}
	munged := re.ReplaceAll(buf.Bytes(), []byte(`time.Unix($1-62135596800,$2)`))

	gofmt := exec.Command("gofmt", "-s")
	gofmt.Stdin = bytes.NewReader(munged)
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
	content map[string]canfs.FileData
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
	bldr.content[path] = canfs.FileData{
		Name:    info.Name(),
		Mode:    info.Mode(),
		Size:    info.Size(),
		ModTime: info.ModTime(),
		Bytes:   buf,
	}
	return nil
}
