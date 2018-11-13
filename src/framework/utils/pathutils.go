package utils

//golibs:https://github.com/SimonWaldherr/golibs

import (
	"os"
)

//CreateDirIfNotExist
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}
