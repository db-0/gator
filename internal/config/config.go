package config

import (
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	db_url            string `json:"db_url"`
	current_user_name string `json:"current_user_name"`
}

func Read() Config {
	configDir, err := os.UserHomeDir()
	if err != nil {
		println(err)
	}

	configFilePath := configDir + configFileName

}
