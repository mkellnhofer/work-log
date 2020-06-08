package main

import (
	"kellnhofer.com/work-log/config"
	"kellnhofer.com/work-log/db"
	"kellnhofer.com/work-log/db/tx"
	"kellnhofer.com/work-log/service"
	"kellnhofer.com/work-log/view/controller"
	"kellnhofer.com/work-log/view/middleware"
)

// Initializer initializes application components.
type Initializer struct {
	conf *config.Config

	db *db.Db

	entryServ *service.EntryService
	sessServ  *service.SessionService
	userServ  *service.UserService
	jobServ   *service.JobService

	errVCtrl   *controller.ErrorController
	authVCtrl  *controller.AuthController
	entryVCtrl *controller.EntryController

	txMidw    *tx.TransactionMiddleware
	errVMidw  *middleware.ErrorMiddleware
	sessVMidw *middleware.SessionMiddleware
	authVMidw *middleware.AuthMiddleware
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

// GetEntryService returns a initialized entry service object.
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

// GetJobService returns a initialized job service object.
func (i *Initializer) GetJobService() *service.JobService {
	if i.jobServ == nil {
		i.jobServ = service.NewJobService(i.GetSessionService())
	}
	return i.jobServ
}

// --- View controller functions ---

// GetErrorViewController returns a initialized error view controller object.
func (i *Initializer) GetErrorViewController() *controller.ErrorController {
	if i.errVCtrl == nil {
		i.errVCtrl = controller.NewErrorController()
	}
	return i.errVCtrl
}

// GetAuthViewController returns a initialized auth view controller object.
func (i *Initializer) GetAuthViewController() *controller.AuthController {
	if i.authVCtrl == nil {
		i.authVCtrl = controller.NewAuthController(i.GetUserService())
	}
	return i.authVCtrl
}

// GetEntryViewController returns a initialized entry view controller object.
func (i *Initializer) GetEntryViewController() *controller.EntryController {
	if i.entryVCtrl == nil {
		i.entryVCtrl = controller.NewEntryController(i.GetUserService(), i.GetEntryService())
	}
	return i.entryVCtrl
}

// --- General middleware functions ---

// GetTransactionMiddleware returns a initialized transaction middleware object.
func (i *Initializer) GetTransactionMiddleware() *tx.TransactionMiddleware {
	if i.txMidw == nil {
		i.txMidw = tx.NewTransactionMiddleware()
	}
	return i.txMidw
}

// --- View middleware functions ---

// GetErrorViewMiddleware returns a initialized error view middleware object.
func (i *Initializer) GetErrorViewMiddleware() *middleware.ErrorMiddleware {
	if i.errVMidw == nil {
		i.errVMidw = middleware.NewErrorMiddleware()
	}
	return i.errVMidw
}

// GetSessionViewMiddleware returns a initialized session view middleware object.
func (i *Initializer) GetSessionViewMiddleware() *middleware.SessionMiddleware {
	if i.sessVMidw == nil {
		i.sessVMidw = middleware.NewSessionMiddleware(i.GetSessionService())
	}
	return i.sessVMidw
}

// GetAuthViewMiddleware returns a initialized auth view middleware object.
func (i *Initializer) GetAuthViewMiddleware() *middleware.AuthMiddleware {
	if i.authVMidw == nil {
		i.authVMidw = middleware.NewAuthMiddleware()
	}
	return i.authVMidw
}
