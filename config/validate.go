package config

import (
	"errors"
	"fmt"
)

func (cfg *Config) Validate() error {
	if len(cfg.Contexts) == 0 {
		return errors.New("no contexts defined")
	}

	if cfg.CurrentContext == "" {
		return errors.New("current context is not set")
	}

	for name, ctx := range cfg.Contexts {
		if err := ctx.Validate(); err != nil {
			return fmt.Errorf("context %q: %w", name, err)
		}
	}

	return nil
}

func (c *Context) Validate() error {
	if len(c.Name) == 0 {
		return errors.New("name must not be empty")
	}

	if err := c.Store.Validate(); err != nil {
		return fmt.Errorf("store: %w", err)
	}

	return nil
}

func (s *StoreConfig) Validate() error {
	switch s.Format {
	case StoreFormatJSONL:
		if s.Path == "" {
			return errors.New("store.path is required for jsonl")
		}
	case StoreFormatSQLite:
		if s.Path == "" {
			return errors.New("store.path is required for sqlite")
		}
	case StoreFormatCustom:
		if s.Path == "" {
			return errors.New("store.path is required for custom")
		}
		if s.Custom == nil {
			return errors.New("store.custom is required for custom")
		}

		if s.Custom.Template == "" {
			return errors.New("store.custom.template is required for custom")
		}

	default:
		return fmt.Errorf("invalid store format: %q", s.Format)
	}
	return nil
}
