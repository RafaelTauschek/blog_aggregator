package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func (cfg *Config) SetUser(username string) error {
	if username == "" {
		return nil
	}

	cfg.CurrentUser = username

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0600)
	if err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("%s/%s", home, configFileName)
	return path, nil
}
