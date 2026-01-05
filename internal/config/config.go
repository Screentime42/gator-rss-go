package config

import (
	"encoding/json"
	"os"
	"fmt"
)

type Config struct {
	URL 					string `json:"db_url"`
	CurrentUserName 	string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	fullPath := home + "/.gatorconfig.json"

	fmt.Println(fullPath)
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(c); err != nil {
		return err
	}
	
	return nil
}
