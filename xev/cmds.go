package xev

import (
	"archive/zip"
	"fmt"
	"fw/cp/module"
	"io/ioutil"
	"path/filepath"
	"ypk/assert"
)

const CODE = "code"

func Load(path, name string) (ret *module.Module) {
	fmt.Println(path + ` ` + name)
	var data []byte
	r, err := zip.OpenReader(filepath.Join(path, CODE, name))
	assert.For(err == nil, 40)
	for _, f := range r.File {
		if f.Name == CODE {
			r, _ := f.Open()
			data, _ = ioutil.ReadAll(r)
		}
	}
	//data, _ = ioutil.ReadFile(filepath.Join(path, CODE, name))
	fmt.Println(len(data))
	if data != nil {
		result := LoadOXF(data)
		ret = DoAST(result)
	}
	return ret
}
