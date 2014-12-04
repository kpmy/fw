package main

import (
	"os"
	"xev"
)

func main() {
	path, _ := os.Getwd()
	xev.Load(path, "PrivDemo1.oxf")
}
