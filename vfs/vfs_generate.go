// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	var VFS http.FileSystem = http.Dir("../assets")
	err := vfsgen.Generate(VFS, vfsgen.Options{
		PackageName:  "vfs",
		BuildTags:    "!dev",
		VariableName: "VFS",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
