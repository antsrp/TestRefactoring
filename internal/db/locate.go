package db

import (
	"os"
	"path/filepath"
)

const (
	DB_RELATIVE_PATH = `db`
)

func getPathToStorageFolder() string {
	curPath, _ := os.Getwd()
	return filepath.Join(curPath, DB_RELATIVE_PATH)
}
