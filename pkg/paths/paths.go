package paths

import (
	"log"
	"os"
)

func CreatePaths(path string) {
	if NotExist(path) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Fatalf("create path: can not create %s - %v\n", path, err)
		}
	}
}

func NotExist(name string) bool {
	_, err := os.Stat(name)
	return os.IsNotExist(err)
}
