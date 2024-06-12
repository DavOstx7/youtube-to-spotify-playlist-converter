package config

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"sync"
)

var (
	configDirPath string
	configDirOnce sync.Once
	configDirErr  error
)

type Config interface{}

func getConfigFilePath(filename string) (string, error) {
	configDirOnce.Do(func() {
		cwd, err := os.Getwd()
		if err != nil {
			configDirErr = err
			return
		}

		configDirPath = filepath.Join(cwd, "../config/")
	})

	return path.Join(configDirPath, filename), configDirErr
}

func loadConfig(filename string, config Config) error {
	configFilePath, err := getConfigFilePath(filename)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, config); err != nil {
		return err
	}

	return nil
}
