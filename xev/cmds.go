package xev

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func Load(path, name string) {
	fmt.Println(path + ` ` + name)
	var data []byte
	data, _ = ioutil.ReadFile(filepath.Join(path, name))
	fmt.Println(len(data))
	if data != nil {
		result := LoadOXF(data)
		DoAST(result)
	}
}
