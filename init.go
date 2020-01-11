package main

import (
	"kellnhofer.com/work-log/config"
	"kellnhofer.com/work-log/db"
)

// Initializer initializes application components.
type Initializer struct {
	conf *config.Config

	db *db.Db
}

// NewInitializer creates a new initializer.
func NewInitializer(c *config.Config) *Initializer {
	return &Initializer{conf: c}
}

// --- Database functions ---

// GetDb returns a initialized DB object.
func (i *Initializer) GetDb() *db.Db {
	if i.db == nil {
		i.db = db.NewDb(i.conf)
	}
	return i.db
}
