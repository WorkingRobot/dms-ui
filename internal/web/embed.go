// Package web embeds the built SvelteKit frontend into the Go binary.
//
// The SvelteKit static build is copied into ./dist at image build time
// (see Dockerfile). The dist directory is committed with a placeholder so
// the module always compiles; `make` / Docker overwrites it with the real
// build before `go build`.
package web

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var embedded embed.FS

// FS returns the frontend filesystem rooted at the build output.
func FS() fs.FS {
	sub, err := fs.Sub(embedded, "dist")
	if err != nil {
		panic(err)
	}
	return sub
}
