package main

import (
	"fmt"
	// "io"
	// "log"
	"strings"

	filetool "github.com/takoyaki-3/go-file-tool"
	// json "github.com/takoyaki-3/go-json"
	gos3 "github.com/takoyaki-3/go-s3"
	"github.com/takoyaki-3/goc"
)

func main() {

	// データ本体のアップロード
	err, files := filetool.DirWalk("./dist", filetool.DirWalkOption{})
	if err != nil {
		return
	}
	goc.Parallel(16, len(files), func(i, rank int) {
		s3, err := gos3.NewSession("s3-conf.json")
		if err != nil {
			fmt.Println(err)
			return
		}

		v := files[i]

		if v.IsDir {
			fmt.Println(err)
			return
		}

		if v.Name == "dist" {
			fmt.Println(err)
			return
		}

		exs := strings.Split(v.Name, ".")
		if len(exs) == 1 {
			fmt.Println(err)
			return
		}
		key := "v0.0.0/" + strings.ReplaceAll(v.Path[len("dist\\"):], "\\", "/")
		fmt.Println(i, len(files), float32(i)/float32(len(files)), key)

		err = s3.UploadFromPath(v.Path, key)
		if err != nil {
			fmt.Println(err)
			return
		}

	})

	return
}
