package main

import (
	"net/http"

	_ "github.com/tomclegg/canfs"
)

//go:generate go run $GOPATH/src/github.com/tomclegg/canfs/generate.go -pkg=main -id=assets -out=assets_generated.go -dir=./assets

func main() {
	http.ListenAndServe(":12345", http.FileServer(assets))
}
