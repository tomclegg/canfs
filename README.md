# canfs

Package canfs embeds file data in a Go program and exports it as an
http.FileSystem.

To use a FileSystem, import canfs and add a "go generate" step to your project.

    import _ "github.com/tomclegg/canfs" // ensures "go get -d mypkg" also gets canfs, so "go generate mypkg" works
    //go:generate go run $GOPATH/src/github.com/tomclegg/canfs/generate.go -id=assetfs -out=assetfs_generated.go -dir=./assets -pkg=main

When you run "go generate [your-package]", canfs will create a file called
"assetfs_generated.go" with a FileSystem variable called "assetfs" whose root
directory is a mirror of the "assets" directory from your source tree.

For example, http.ListenAndServe(":", assetfs) will respond to "/foo" with
content from "./assets/foo" from your source directory.


### Options

-out=fnm sets the output (generated code) filename.

-dir=path sets the source directory.

-pkg=name sets the package name in the generated code.

-id=name sets the name of the filesystem variable in the generated code.


### Features And Limitations

File metadata (notably modification times) are preserved.

The generated files are verbose. Consider adding "*_generated.go" to .gitignore.

File data is not compressed. Consider using upx.

Symbolic links to files are followed. Symbolic links to directories are not
followed.

Directory listings are not yet supported.
