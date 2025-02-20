package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(*c)
}

func getConfigFilePath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(cwd, configFileName)
	return fullPath, nil
}

func write(c Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(c)
	if err != nil {
		return err
	}

	return nil
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}
	defer data.Close()

	decoder := json.NewDecoder(data)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
