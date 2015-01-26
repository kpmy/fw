package xev

import (
	"archive/zip"
	"fmt"
	"fw/cp/module"
	"io"
	"path/filepath"
	"ypk/assert"
)

const CODE = "code"

func Load(path, name string) (ret *module.Module) {
	fmt.Println(path + ` ` + name)
	//var data []byte
	var rd io.Reader
	r, err := zip.OpenReader(filepath.Join(path, CODE, name))
	assert.For(err == nil, 40)
	for _, f := range r.File {
		if f.Name == CODE {
			rd, _ = f.Open()
			//data, _ = ioutil.ReadAll(r)
		}
	}
	//data, _ = ioutil.ReadFile(filepath.Join(path, CODE, name))
	if r != nil {
		result := LoadOXF(rd)
		ret = DoAST(result)
		fmt.Println("load", len(ret.Nodes), "nodes")
	}
	return ret
}
