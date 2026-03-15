package config

import (
	"os"
	"path"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DBPath        string             `yaml:"db_path" mapstructure:"db_path" validate:"required"`
	ActiveContext string             `yaml:"active_context" mapstructure:"active_context" validate:"required"`
	Contexts      map[string]Context `yaml:"contexts" mapstructure:"contexts"`
}

type Context struct {
	Name        string `yaml:"name" mapstructure:"name"`
	Description string `yaml:"description" mapstructure:"description"`
}

var defaultConfig = Config{
	DBPath:        "~/.journl/journl.db",
	ActiveContext: "default",
	Contexts: map[string]Context{
		"default": {Name: "default", Description: "The default context"},
	},
}

var config Config

var (
	configDirPath  string
	configFilePath string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configDirPath = path.Join(home, ".journl")
	configFilePath = path.Join(configDirPath, "config.yaml")
}

func Get() Config {
	return config
}

func InitConfig(cfgFile string) error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		if _, err := os.Stat(configFilePath); err != nil {
			if os.IsNotExist(err) {
				if err := SetConfig(defaultConfig); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		// Search config in home directory with name ".journl" (without extension).
		viper.AddConfigPath(configDirPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	return nil
}

func SetConfig(cfg Config) error {
	bs, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}
	if err := os.WriteFile(configFilePath, bs, 0o644); err != nil {
		return err
	}

	config = cfg

	return nil
}
