package config

import "path/filepath"

func GetRootDir() string {
	dir, _ := filepath.Abs(".")
	return dir
}
