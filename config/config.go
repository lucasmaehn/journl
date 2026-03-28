package config

type Config struct {
	CurrentContext string             `yaml:"current_context" mapstructure:"current_context"`
	Contexts       map[string]Context `yaml:"contexts" mapstructure:"contexts"`
}

type Context struct {
	Name        string      `yaml:"name" mapstructure:"name"`
	Description string      `yaml:"description" mapstructure:"description"`
	Store       StoreConfig `yaml:"store" mapstructure:"store"`
}

type StoreConfig struct {
	Format StoreFormat   `yaml:"format" mapstructure:"format"`
	Path   string        `yaml:"path" mapstructure:"path"`
	JSONL  *JSONLConfig  `yaml:"jsonl,omitempty" mapstructure:"jsonl"`
	SQLite *SQLiteConfig `yaml:"sqlite,omitempty" mapstructure:"sqlite"`
	Custom *CustomConfig `yaml:"custom,omitempty" mapstructure:"custom"`
}

type StoreFormat string

const (
	StoreFormatJSONL  StoreFormat = "jsonl"
	StoreFormatSQLite StoreFormat = "sqlite"
	StoreFormatCustom StoreFormat = "custom"
)

type (
	JSONLConfig  struct{}
	SQLiteConfig struct{}
	CustomConfig struct {
		Template string `yaml:"template" mapstructure:"template"`
	}
)
