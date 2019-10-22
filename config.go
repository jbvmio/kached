package kached

import (
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/ristretto"
)

// Config holds the configuration options for a cached database.
type Config struct {
	Cache    *ristretto.Config
	Database badger.Options
}

// NewConfig returns a new Config using defaults.
func NewConfig(dbPath string) *Config {
	c := &ristretto.Config{
		NumCounters: 10000,
		MaxCost:     1000,
		BufferItems: 64,
	}
	d := badger.DefaultOptions(dbPath)
	return &Config{
		Cache:    c,
		Database: d,
	}
}
