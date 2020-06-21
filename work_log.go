package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"kellnhofer.com/work-log/config"
	"kellnhofer.com/work-log/constant"
	"kellnhofer.com/work-log/loc"
	"kellnhofer.com/work-log/log"
)

func main() {
	// Load config
	conf := config.LoadConfig()

	// Set logging level
	log.SetLevel(conf.LogLevel)

	// Load localization
	loc.LoadLocalization(conf.LocLanguage)

	log.Infof("Starting Work Log server %s.", constant.AppVersion)

	// Create initializer
	init := NewInitializer(conf)

	// Open and create/update database
	db := init.GetDb()
	db.OpenDb()
	defer db.CloseDb()
	db.UpdateDb()

	// Schedule jobs
	init.GetJobService().ScheduleJobs()

	// Create router
	router := mux.NewRouter().StrictSlash(true)

	// Configure view routing
	configureViewRouting(init, router)
	// Configure API routing
	configureApiRouting(init, router)

	// Register router
	http.Handle("/", router)

	// Start HTTP server
	log.Infof("Listen on port '%d'.", conf.ServerPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", conf.ServerPort), nil)
	if err != nil {
		log.Fatalf("Could not start server! (Error: %s)", err)
	}
}

func configureViewRouting(init *Initializer, r *mux.Router) {
	// Create public middleware route
	pubRoute := negroni.New()
	pubRoute.Use(init.GetTransactionMiddleware())
	pubRoute.Use(init.GetErrorViewMiddleware())
	pubRoute.Use(init.GetSessionViewMiddleware())
	pubRoute.Use(init.GetSecurityViewMiddleware())
	// Create protected middleware route
	proRoute := negroni.New()
	proRoute.Use(init.GetTransactionMiddleware())
	proRoute.Use(init.GetErrorViewMiddleware())
	proRoute.Use(init.GetSessionViewMiddleware())
	proRoute.Use(init.GetSecurityViewMiddleware())
	proRoute.Use(init.GetAuthCheckViewMiddleware())

	// Get controllers
	errCtrl := init.GetErrorViewController()
	authCtrl := init.GetAuthViewController()
	entryCtrl := init.GetEntryViewController()

	// Add public endpoints
	addEndpoint(r, pubRoute, "GET", "/", getRootHandler())
	addEndpoint(r, pubRoute, "GET", "/error", errCtrl.GetErrorHandler())
	addEndpoint(r, pubRoute, "GET", "/login", authCtrl.GetLoginHandler())
	addEndpoint(r, pubRoute, "POST", "/login", authCtrl.PostLoginHandler())
	// Add protected endpoints
	addEndpoint(r, proRoute, "GET", "/logout", authCtrl.GetLogoutHandler())
	addEndpoint(r, proRoute, "GET", "/list", entryCtrl.GetListHandler())
	addEndpoint(r, proRoute, "GET", "/create", entryCtrl.GetCreateHandler())
	addEndpoint(r, proRoute, "POST", "/create", entryCtrl.PostCreateHandler())
	addEndpoint(r, proRoute, "GET", "/edit/{id}", entryCtrl.GetEditHandler())
	addEndpoint(r, proRoute, "POST", "/edit/{id}", entryCtrl.PostEditHandler())
	addEndpoint(r, proRoute, "GET", "/copy/{id}", entryCtrl.GetCopyHandler())
	addEndpoint(r, proRoute, "POST", "/copy/{id}", entryCtrl.PostCopyHandler())
	addEndpoint(r, proRoute, "POST", "/delete/{id}", entryCtrl.PostDeleteHandler())
	addEndpoint(r, proRoute, "GET", "/search", entryCtrl.GetSearchHandler())
	addEndpoint(r, proRoute, "POST", "/search", entryCtrl.PostSearchHandler())
	addEndpoint(r, proRoute, "GET", "/overview", entryCtrl.GetOverviewHandler())
	addEndpoint(r, proRoute, "POST", "/overview", entryCtrl.PostOverviewHandler())
	addEndpoint(r, proRoute, "GET", "/overview/export", entryCtrl.GetOverviewExportHandler())
	// Add resource endpoints
	fileSrv := http.FileServer(http.Dir("./resources"))
	r.Handle("/css/{name}", fileSrv).Methods("GET")
	r.Handle("/img/{name}", fileSrv).Methods("GET")
	r.Handle("/font/{name}", fileSrv).Methods("GET")
	r.Handle("/js/{name}", fileSrv).Methods("GET")
}

func getRootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, constant.PathDefault, http.StatusFound)
	}
}

func configureApiRouting(init *Initializer, r *mux.Router) {
	// Create API sub route
	ar := r.PathPrefix("/api/v1").Subrouter()

	// Create protected middleware route
	proRoute := negroni.New()
	proRoute.Use(init.GetTransactionMiddleware())
	proRoute.Use(init.GetErrorApiMiddleware())
	proRoute.Use(init.GetSecurityApiMiddleware())
	proRoute.Use(init.GetAuthCheckApiMiddleware())

	// Get controllers
	entryCtrl := init.GetEntryApiController()
	userCtrl := init.GetUserApiController()

	// Add protected endpoints
	addEndpoint(ar, proRoute, "GET", "/entries", entryCtrl.GetEntriesHandler())
	addEndpoint(ar, proRoute, "POST", "/entries", entryCtrl.CreateEntryHandler())
	addEndpoint(ar, proRoute, "GET", "/entries/{id}", entryCtrl.GetEntryHandler())
	addEndpoint(ar, proRoute, "PUT", "/entries/{id}", entryCtrl.UpdateEntryHandler())
	addEndpoint(ar, proRoute, "DELETE", "/entries/{id}", entryCtrl.DeleteEntryHandler())
	addEndpoint(ar, proRoute, "GET", "/entry_types", entryCtrl.GetEntryTypesHandler())
	addEndpoint(ar, proRoute, "GET", "/entry_activities", entryCtrl.GetEntryActivitiesHandler())
	addEndpoint(ar, proRoute, "POST", "/entry_activities", entryCtrl.CreateEntryActivityHandler())
	addEndpoint(ar, proRoute, "PUT", "/entry_activities/{id}",
		entryCtrl.UpdateEntryActivityHandler())
	addEndpoint(ar, proRoute, "DELETE", "/entry_activities/{id}",
		entryCtrl.DeleteEntryActivityHandler())
	addEndpoint(ar, proRoute, "GET", "/user", userCtrl.GetCurrentUserHandler())
	addEndpoint(ar, proRoute, "GET", "/user/roles", userCtrl.GetCurrentUserRolesHandler())
	addEndpoint(ar, proRoute, "GET", "/users", userCtrl.GetUsersHandler())
	addEndpoint(ar, proRoute, "POST", "/users", userCtrl.CreateUserHandler())
	addEndpoint(ar, proRoute, "GET", "/users/{id}", userCtrl.GetUserHandler())
	addEndpoint(ar, proRoute, "PUT", "/users/{id}", userCtrl.UpdateUserHandler())
	addEndpoint(ar, proRoute, "DELETE", "/users/{id}", userCtrl.DeleteUserHandler())
	addEndpoint(ar, proRoute, "PUT", "/users/{id}/password", userCtrl.UpdateUserPasswordHandler())
	addEndpoint(ar, proRoute, "GET", "/users/{id}/roles", userCtrl.GetUserRolesHandler())
	addEndpoint(ar, proRoute, "PUT", "/users/{id}/roles", userCtrl.UpdateUserRolesHandler())
}

// --- Helper functions ---

func addEndpoint(r *mux.Router, m *negroni.Negroni, method string, path string, h http.HandlerFunc) {
	r.Handle(path, createHandler(m, h)).Methods(method)
}

func createHandler(r *negroni.Negroni, h http.HandlerFunc) http.Handler {
	nr := r.With()
	nr.UseHandlerFunc(h)
	return nr
}
