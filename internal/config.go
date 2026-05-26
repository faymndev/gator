package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const ConfigFile = ".gatorconfig.json"

func getConfigFile() (string, error) {
	baseDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home dir: %w", err)
	}

	filePath := filepath.Join(baseDir, ConfigFile)
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to create config file: %w", err)
		}
		file.Close()
	} else if err != nil {
		return "", fmt.Errorf("failed to check config file status: %w", err)
	}

	return filePath, nil
}

func Read() (*Config, error) {
	filePath, err := getConfigFile()

	data, err := os.ReadFile(filePath)
	if err != nil {
		return &Config{}, err
	}

	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		// ignore error since config was bad, we'll just return a fresh config obj
		return &Config{}, err
	}
	return config, nil
}

func (c *Config) Write() error {
	configFile, err := getConfigFile()
	if err != nil {
		return err
	}

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFile, data, 0644)
	return err
}
