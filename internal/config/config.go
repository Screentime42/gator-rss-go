package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)
const configFileName = ".gatorconfig.json"


type Config struct {
	DBURL 					string `json:"db_url"`
	CurrentUserName 		string `json:"current_user_name"`
}


func Read() (Config, error) {
	var cfg Config
	
	fullPath, err := getConfigFilePath()
	if err != nil {
		return cfg, err
	}

	f, err := os.Open(fullPath)
	if err != nil {
		return cfg, err
	}

	defer f.Close()

	dec := json.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}


func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(*c)
}


func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configFileName), nil
}


func write(cfg Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(cfg); err != nil {
		return err
	}

	return nil
}
