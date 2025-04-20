package core

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// default baseline
	defaultCfg := Config{
		Domain:             "localhost:8080",
		Port:               ":8080",
		Whitelist:          []string{},
		Blacklist:          []string{},
		Translations:       map[string]string{},
		ShowDisclaimer:     true,
		WhitelistOnly:      false,
		EnableTranslations: false,
	}

	tests := []struct {
		name        string
		fileContent string // if empty, no file is written
		want        Config
	}{
		{
			name:        "no file → defaults",
			fileContent: "",
			want:        defaultCfg,
		},
		{
			name: "valid JSON overrides everything",
			fileContent: `{
				"domain":"example.com",
				"port":":1234",
				"whitelist":["1.1.1.1","2.2.2.2"],
				"blacklist":["9.9.9.9"],
				"translations":{"en":"Hello","es":"Hola"},
				"showDisclaimer":false,
				"whitelistOnly":true,
				"enableTranslations":true
			}`,
			want: Config{
				Domain:             "example.com",
				Port:               ":1234",
				Whitelist:          []string{"1.1.1.1", "2.2.2.2"},
				Blacklist:          []string{"9.9.9.9"},
				Translations:       map[string]string{"en": "Hello", "es": "Hola"},
				ShowDisclaimer:     false,
				WhitelistOnly:      true,
				EnableTranslations: true,
			},
		},
		{
			name:        "partial JSON → single field",
			fileContent: `{"domain":"partial.local"}`,
			want: func() Config {
				c := defaultCfg
				c.Domain = "partial.local"
				return c
			}(),
		},
		{
			name:        "invalid JSON → defaults",
			fileContent: `{ this is not valid JSON `,
			want:        defaultCfg,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// isolate in temp dir
			origDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("could not Getwd: %v", err)
			}
			defer func() {
				_ = os.Chdir(origDir)
			}()
			tmp := t.TempDir()
			if err := os.Chdir(tmp); err != nil {
				t.Fatalf("could not chdir to temp: %v", err)
			}

			// write file if requested
			if tt.fileContent != "" {
				fpath := filepath.Join(tmp, "torgate.json")
				if err := os.WriteFile(fpath, []byte(tt.fileContent), 0644); err != nil {
					t.Fatalf("could not write config file: %v", err)
				}
			}

			gotPtr := LoadConfig()
			got := *gotPtr

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfig() =\n%#v\nwant\n%#v", got, tt.want)
			}
		})
	}
}
