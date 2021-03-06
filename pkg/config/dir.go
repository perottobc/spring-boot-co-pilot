package config

import (
	"errors"
	"fmt"
	"github.com/co-pilot-cli/co-pilot/pkg/file"
)

func (dirCfg DirConfig) Dir() string {
	return dirCfg.Path
}

func (dirCfg DirConfig) FilePath(fileName string) (string, error) {
	path := file.Path("%s/%s", dirCfg.Dir(), fileName)
	if !file.Exists(path) {
		return "", errors.New(fmt.Sprintf("could not find %s in cloud config", fileName))
	}

	return path, nil
}
