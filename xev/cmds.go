package xev

import (
	"cp/node"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func Load(path, name string) (ret node.Node) {
	fmt.Println(path + ` ` + name)
	var data []byte
	data, _ = ioutil.ReadFile(filepath.Join(path, name))
	fmt.Println(len(data))
	if data != nil {
		result := LoadOXF(data)
		ret = DoAST(result)
	}
	return ret
}
