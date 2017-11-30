# canfs

    import "github.com/tomclegg/canfs"

Package canfs embeds file data in a Go program and exports it as an
http.FileSystem.

To use a FileSystem, import canfs and add a "go generate" step to your project:

    import _ "github.com/tomclegg/canfs"
    //go:generate go run $GOPATH/src/github.com/tomclegg/canfs/generate.go -id=assetfs -out=assetfs_generated.go -dir=./assets

When you run "go generate [your-package]", canfs will create a file called
"assetfs_generated.go" with a FileSystem variable called "assetfs" whose root
directory is a mirror of the "assets" directory from your source tree.

For example, http.ListenAndServe(":", assetfs) will respond to "/foo" with
content from "./assets/foo" from your source directory.


Features and limitations

File metadata (notably modification times) are preserved.

The generated files are verbose. Consider adding "*_generated.go" to .gitignore.

File data is not compressed. Consider using upx.

Symbolic links to files are followed. Symbolic links to directories are not
followed.

Directory listings are not yet supported.

## Usage

#### type FileData

```go
type FileData struct {
	Name    string
	Mode    os.FileMode
	Size    int64
	ModTime time.Time
	Bytes   []byte
}
```

FileData is embedded file content and metadata.

#### type FileSystem

```go
type FileSystem struct {
	Content map[string]FileData
}
```

FileSystem implements http.FileSystem with embedded content.

#### func (FileSystem) Open

```go
func (fs FileSystem) Open(path string) (http.File, error)
```
Open implements http.FileSystem.
