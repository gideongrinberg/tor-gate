package core

type Config struct {
	Mode         string            `json:"mode,omitempty"`
	Whitelist    []string          `json:"whitelist,omitempty"`
	Blacklist    []string          `json:"blacklist,omitempty"`
	Translations map[string]string `json:"translations,omitempty"`
}

// func LoadConfig() {
// 	defaultConfig := &Config{
// 		Mode:         "Blacklist",
// 		Whitelist:    []string{},
// 		Blacklist:    []string{},
// 		Translations: make(map[string]string),
// 	}

// 	if _, err := os.Stat("./torgate.json"); err == nil {
// 		log
// 	}
// }
