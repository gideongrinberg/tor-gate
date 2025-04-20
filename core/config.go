package core

import (
	"encoding/json"
	"os"
)

type Config struct {
	Domain             string            `json:"domain,omitempty"`
	Whitelist          []string          `json:"whitelist,omitempty"`
	Blacklist          []string          `json:"blacklist,omitempty"`
	Translations       map[string]string `json:"translations,omitempty"`
	ShowDisclaimer     bool              `json:"showDisclaimer,omitempty"`
	WhitelistOnly      bool              `json:"whitelistOnly,omitempty"`
	EnableTranslations bool              `json:"enableTranslations,omitempty"`
}

func LoadConfig() *Config {
	config := &Config{
		Domain:         "localhost:8080",
		Whitelist:      []string{},
		Blacklist:      []string{},
		Translations:   make(map[string]string),
		ShowDisclaimer: true,
	}

	if _, err := os.Stat("./torgate.json"); err == nil {
		Logger().Println("Loading config from file.")
		file, err := os.ReadFile("./torgate.json")
		if err != nil {
			Logger().Println("Failed to read config file. Using default options.")
		} else {
			json.Unmarshal(file, config)
		}
	} else {
		Logger().Println("Could not find torgate.json, loading default config.")
	}

	confStr, _ := json.Marshal(config)
	Logger().Println("Using config:", string(confStr))

	return config
}
