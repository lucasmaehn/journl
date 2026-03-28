package logstore

import (
	"fmt"

	"github.com/lucasmaehn/journl/config"
)

func New(ctxName string, cfg config.StoreConfig) (LogStore, error) {
	switch cfg.Format {
	case config.StoreFormatJSONL:
		return NewJSONL(ctxName, cfg)
	case config.StoreFormatSQLite:
		return NewSQLite(ctxName, cfg)
	case config.StoreFormatCustom:
		return NewCustom(ctxName, cfg)
	default:
		return nil, fmt.Errorf("invalid StoreFormat: %q", cfg.Format)
	}
}
