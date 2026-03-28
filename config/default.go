package config

import "path"

func Default(homeDir string) *Config {
	return &Config{
		CurrentContext: "default",
		Contexts: map[string]Context{
			"default": {
				Name:        "default",
				Description: "The default context for journl",
				Store: StoreConfig{
					Format: StoreFormatSQLite,
					Path:   path.Join(homeDir, ".journl", "db.sqlite"),
					SQLite: &SQLiteConfig{},
				},
			},
		},
	}
}
