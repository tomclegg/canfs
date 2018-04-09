// SPDX-License-Identifier: CC0-1.0

// Package canfs embeds file data in a Go program and exports it as an
// http.FileSystem.
//
// To use a FileSystem, import canfs and add a "go generate" step to
// your project.
//
//   import _ "github.com/tomclegg/canfs" // ensures "go get -d mypkg" also gets canfs, so "go generate mypkg" works
//   //go:generate go run $GOPATH/src/github.com/tomclegg/canfs/generate.go -id=assetfs -out=assetfs_generated.go -dir=./assets -pkg=main
//
// When you run "go generate [your-package]", canfs will create a file
// called "assetfs_generated.go" with a FileSystem variable called
// "assetfs" whose root directory is a mirror of the "assets"
// directory from your source tree.
//
// Then, to serve "/foo" using content from "./assets/foo" in your
// source directory:
//
//   http.ListenAndServe(":", http.FileServer(assetfs))
//
// Options
//
// -out=fnm sets the output (generated code) filename.
//
// -dir=path sets the source directory.
//
// -pkg=name sets the package name in the generated code.
//
// -id=name sets the name of the filesystem variable in the generated
// code.
//
// Features And Limitations
//
// File metadata (notably modification times) are preserved.
//
// The generated files are verbose. Consider adding "*_generated.go"
// to .gitignore.
//
// File data is not compressed. Consider using upx.
//
// Symbolic links to files are followed. Symbolic links to directories
// are not followed.
//
// Directory listings are not yet supported.
package canfs

import (
	"bytes"
	"net/http"
	"os"
	"strings"
	"time"
)

// FileSystem implements http.FileSystem with embedded content.
type FileSystem struct {
	Content map[string]FileInfo
}

type FileData interface {
	Bytes() []byte
}

type ByteData struct {
	Data []byte
}

func (bd ByteData) Bytes() []byte {
	return bd.Data
}

type StringData struct {
	Data string
}

func (sd StringData) Bytes() []byte {
	return []byte(sd.Data)
}

// Open implements http.FileSystem.
func (fs FileSystem) Open(path string) (http.File, error) {
	if strings.HasSuffix(path, "/") {
		return file{FileInfo: FileInfo{M: os.ModeDir}}, nil
	}
	fi, ok := fs.Content[path]
	if !ok {
		return nil, os.ErrNotExist
	}
	return file{
		Reader:   bytes.NewReader(fi.Bytes()),
		FileInfo: fi,
	}, nil
}

type file struct {
	FileInfo
	*bytes.Reader
}

func (file) Close() error {
	return nil
}

func (file) Readdir(int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f file) Stat() (os.FileInfo, error) {
	return f.FileInfo, nil
}

type FileInfo struct {
	N  string      // name
	M  os.FileMode // mode
	S  int64       // size
	MT int64       // modtime
	FileData
}

func (fi FileInfo) Name() string {
	return fi.N
}

func (fi FileInfo) Mode() os.FileMode {
	return fi.M
}

func (fi FileInfo) Size() int64 {
	return fi.S
}

func (fi FileInfo) ModTime() time.Time {
	return time.Unix(0, fi.MT)
}

func (fi FileInfo) IsDir() bool {
	return fi.M&os.ModeDir != 0
}

func (fi FileInfo) Sys() interface{} {
	return nil
}

//go:generate sh -c "godocdown >README.md"
