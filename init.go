package main

import (
	"kellnhofer.com/work-log/config"
	"kellnhofer.com/work-log/db"
	"kellnhofer.com/work-log/middleware"
	"kellnhofer.com/work-log/service"
)

// Initializer initializes application components.
type Initializer struct {
	conf *config.Config

	db *db.Db

	entryServ *service.EntryService
	sessServ  *service.SessionService
	userServ  *service.UserService

	errMidw *middleware.ErrorMiddleware
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

// --- Service functions ---

// GetEntryService returns a initialized work entry service object.
func (i *Initializer) GetEntryService() *service.EntryService {
	if i.entryServ == nil {
		i.entryServ = service.NewEntryService(i.GetDb().GetEntryRepo())
	}
	return i.entryServ
}

// GetSessionService returns a initialized session service object.
func (i *Initializer) GetSessionService() *service.SessionService {
	if i.sessServ == nil {
		i.sessServ = service.NewSessionService(i.GetDb().GetSessionRepo())
	}
	return i.sessServ
}

// GetUserService returns a initialized user service object.
func (i *Initializer) GetUserService() *service.UserService {
	if i.userServ == nil {
		i.userServ = service.NewUserService(i.GetDb().GetUserRepo())
	}
	return i.userServ
}

// --- Middleware functions ---

// GetErrorMiddleware returns a initialized API auth middleware object.
func (i *Initializer) GetErrorMiddleware() *middleware.ErrorMiddleware {
	if i.errMidw == nil {
		i.errMidw = middleware.NewErrorMiddleware()
	}
	return i.errMidw
}
