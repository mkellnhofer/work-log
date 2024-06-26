package main

import (
	ac "kellnhofer.com/work-log/api/controller"
	am "kellnhofer.com/work-log/api/middleware"
	"kellnhofer.com/work-log/pkg/config"
	"kellnhofer.com/work-log/pkg/db"
	"kellnhofer.com/work-log/pkg/db/tx"
	"kellnhofer.com/work-log/pkg/service"
	vc "kellnhofer.com/work-log/web/controller"
	vm "kellnhofer.com/work-log/web/middleware"
)

// Initializer initializes application components.
type Initializer struct {
	conf *config.Config

	db *db.Db

	entryServ *service.EntryService
	sessServ  *service.SessionService
	userServ  *service.UserService
	jobServ   *service.JobService

	errVCtrl      *vc.ErrorController
	authVCtrl     *vc.AuthController
	entryVCtrl    *vc.EntryController
	logVCtrl      *vc.LogController
	overviewVCtrl *vc.OverviewController
	searchVCtrl   *vc.SearchController
	entryACtrl    *ac.EntryController
	userACtrl     *ac.UserController

	txMidw    *tx.TransactionMiddleware
	errVMidw  *vm.ErrorMiddleware
	sessVMidw *vm.SessionMiddleware
	secVMidw  *vm.SecurityMiddleware
	authVMidw *vm.AuthCheckMiddleware
	errAMidw  *am.ErrorMiddleware
	secAMidw  *am.SecurityMiddleware
	authAMidw *am.AuthCheckMiddleware
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
		i.entryServ = service.NewEntryService(i.GetDb().GetTransactionManager(),
			i.GetDb().GetEntryRepo())
	}
	return i.entryServ
}

// GetSessionService returns a initialized session service object.
func (i *Initializer) GetSessionService() *service.SessionService {
	if i.sessServ == nil {
		i.sessServ = service.NewSessionService(i.GetDb().GetTransactionManager(),
			i.GetDb().GetSessionRepo())
	}
	return i.sessServ
}

// GetUserService returns a initialized user service object.
func (i *Initializer) GetUserService() *service.UserService {
	if i.userServ == nil {
		i.userServ = service.NewUserService(i.GetDb().GetTransactionManager(),
			i.GetDb().GetUserRepo(), i.GetDb().GetContractRepo())
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
func (i *Initializer) GetErrorViewController() *vc.ErrorController {
	if i.errVCtrl == nil {
		i.errVCtrl = vc.NewErrorController()
	}
	return i.errVCtrl
}

// GetAuthViewController returns a initialized auth view controller object.
func (i *Initializer) GetAuthViewController() *vc.AuthController {
	if i.authVCtrl == nil {
		i.authVCtrl = vc.NewAuthController(i.GetUserService())
	}
	return i.authVCtrl
}

// GetEntryViewController returns a initialized entry view controller object.
func (i *Initializer) GetEntryViewController() *vc.EntryController {
	if i.entryVCtrl == nil {
		i.entryVCtrl = vc.NewEntryController(i.GetUserService(), i.GetEntryService())
	}
	return i.entryVCtrl
}

// GetLogViewController returns a initialized log view controller object.
func (i *Initializer) GetLogViewController() *vc.LogController {
	if i.logVCtrl == nil {
		i.logVCtrl = vc.NewLogController(i.GetUserService(), i.GetEntryService())
	}
	return i.logVCtrl
}

// GetOverviewViewController returns a initialized overview view controller object.
func (i *Initializer) GetOverviewViewController() *vc.OverviewController {
	if i.overviewVCtrl == nil {
		i.overviewVCtrl = vc.NewOverviewController(i.GetUserService(), i.GetEntryService())
	}
	return i.overviewVCtrl
}

// GetSearchViewController returns a initialized search view controller object.
func (i *Initializer) GetSearchViewController() *vc.SearchController {
	if i.searchVCtrl == nil {
		i.searchVCtrl = vc.NewSearchController(i.GetUserService(), i.GetEntryService())
	}
	return i.searchVCtrl
}

// --- API controller functions ---

// GetEntryApiController returns a initialized entry API controller object.
func (i *Initializer) GetEntryApiController() *ac.EntryController {
	if i.entryACtrl == nil {
		i.entryACtrl = ac.NewEntryController(i.GetEntryService())
	}
	return i.entryACtrl
}

// GetUserApiController returns a initialized user API controller object.
func (i *Initializer) GetUserApiController() *ac.UserController {
	if i.userACtrl == nil {
		i.userACtrl = ac.NewUserController(i.GetUserService())
	}
	return i.userACtrl
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
func (i *Initializer) GetErrorViewMiddleware() *vm.ErrorMiddleware {
	if i.errVMidw == nil {
		i.errVMidw = vm.NewErrorMiddleware()
	}
	return i.errVMidw
}

// GetSessionViewMiddleware returns a initialized session view middleware object.
func (i *Initializer) GetSessionViewMiddleware() *vm.SessionMiddleware {
	if i.sessVMidw == nil {
		i.sessVMidw = vm.NewSessionMiddleware(i.GetSessionService())
	}
	return i.sessVMidw
}

// GetSecurityViewMiddleware returns a initialized security view middleware object.
func (i *Initializer) GetSecurityViewMiddleware() *vm.SecurityMiddleware {
	if i.secVMidw == nil {
		i.secVMidw = vm.NewSecurityMiddleware(i.GetUserService())
	}
	return i.secVMidw
}

// GetAuthCheckViewMiddleware returns a initialized auth check view middleware object.
func (i *Initializer) GetAuthCheckViewMiddleware() *vm.AuthCheckMiddleware {
	if i.authVMidw == nil {
		i.authVMidw = vm.NewAuthCheckMiddleware(i.GetUserService())
	}
	return i.authVMidw
}

// --- API middleware functions ---

// GetErrorApiMiddleware returns a initialized error API middleware object.
func (i *Initializer) GetErrorApiMiddleware() *am.ErrorMiddleware {
	if i.errAMidw == nil {
		i.errAMidw = am.NewErrorMiddleware()
	}
	return i.errAMidw
}

// GetSecurityApiMiddleware returns a initialized security API middleware object.
func (i *Initializer) GetSecurityApiMiddleware() *am.SecurityMiddleware {
	if i.secAMidw == nil {
		i.secAMidw = am.NewSecurityMiddleware(i.GetUserService())
	}
	return i.secAMidw
}

// GetAuthCheckApiMiddleware returns a initialized auth check API middleware object.
func (i *Initializer) GetAuthCheckApiMiddleware() *am.AuthCheckMiddleware {
	if i.authAMidw == nil {
		i.authAMidw = am.NewAuthCheckMiddleware()
	}
	return i.authAMidw
}
