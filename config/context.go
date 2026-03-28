package config

import (
	"errors"
	"fmt"
)

var (
	ErrContextNotFound      = errors.New("context not found")
	ErrContextAlreadyExists = errors.New("context already exists")
)

func (cfg *Config) ActiveContext() (Context, error) {
	ctx, err := cfg.GetContext(cfg.CurrentContext)
	if err != nil {
		return Context{}, fmt.Errorf("active context: %w", err)
	}

	return ctx, nil
}

func (cfg *Config) GetContext(name string) (Context, error) {
	ctx, ok := cfg.Contexts[name]
	if !ok {
		return Context{}, fmt.Errorf("getting context: %w: context %q", ErrContextNotFound, name)
	}

	return ctx, nil
}

func (cfg *Config) AddContext(name string, ctx Context) error {
	if _, set := cfg.Contexts[name]; set {
		return fmt.Errorf("adding context: %w: context %q", ErrContextAlreadyExists, name)
	}

	cfg.Contexts[name] = ctx
	return nil
}

func (cfg *Config) DeleteContext(name string) error {
	if cfg.CurrentContext == name {
		return fmt.Errorf("deleting context: cannot delete active context %q", name)
	}

	if _, set := cfg.Contexts[name]; !set {
		return fmt.Errorf("deleting context: %w: context %q", ErrContextNotFound, name)
	}

	delete(cfg.Contexts, name)
	return nil
}

func (cfg *Config) UseContext(name string) error {
	if _, set := cfg.Contexts[name]; !set {
		return fmt.Errorf("using context: %w: context %q", ErrContextNotFound, name)
	}

	cfg.CurrentContext = name
	return nil
}

func (cfg *Config) ListContexts() map[string]Context {
	return cfg.Contexts
}
