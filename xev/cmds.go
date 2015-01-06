package xev

import (
	"fmt"
	"fw/cp/module"
	"io/ioutil"
	"path/filepath"
)

const CODE = "code"

func Load(path, name string) (ret *module.Module) {
	fmt.Println(path + ` ` + name)
	var data []byte
	data, _ = ioutil.ReadFile(filepath.Join(path, CODE, name))
	fmt.Println(len(data))
	if data != nil {
		result := LoadOXF(data)
		ret = DoAST(result)
	}
	return ret
}
