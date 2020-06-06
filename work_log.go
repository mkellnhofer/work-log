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

	// Create public view middleware route
	pubVRoute := negroni.New()
	pubVRoute.Use(init.GetErrorViewMiddleware())
	pubVRoute.Use(init.GetSessionViewMiddleware())
	// Create protected view middleware route
	proVRoute := negroni.New()
	proVRoute.Use(init.GetErrorViewMiddleware())
	proVRoute.Use(init.GetSessionViewMiddleware())
	proVRoute.Use(init.GetAuthViewMiddleware())

	// Get view controllers
	errViewCtrl := init.GetErrorViewController()
	authViewCtrl := init.GetAuthViewController()
	entryViewCtrl := init.GetEntryViewController()

	// Add public view endpoints
	addEndpoint(router, pubVRoute, "GET", "/", getRootHandler())
	addEndpoint(router, pubVRoute, "GET", "/error", errViewCtrl.GetErrorHandler())
	addEndpoint(router, pubVRoute, "GET", "/login", authViewCtrl.GetLoginHandler())
	addEndpoint(router, pubVRoute, "POST", "/login", authViewCtrl.PostLoginHandler())
	// Add protected view endpoints
	addEndpoint(router, proVRoute, "GET", "/logout", authViewCtrl.GetLogoutHandler())
	addEndpoint(router, proVRoute, "GET", "/list", entryViewCtrl.GetListHandler())
	addEndpoint(router, proVRoute, "GET", "/create", entryViewCtrl.GetCreateHandler())
	addEndpoint(router, proVRoute, "POST", "/create", entryViewCtrl.PostCreateHandler())
	addEndpoint(router, proVRoute, "GET", "/edit/{id}", entryViewCtrl.GetEditHandler())
	addEndpoint(router, proVRoute, "POST", "/edit/{id}", entryViewCtrl.PostEditHandler())
	addEndpoint(router, proVRoute, "GET", "/copy/{id}", entryViewCtrl.GetCopyHandler())
	addEndpoint(router, proVRoute, "POST", "/copy/{id}", entryViewCtrl.PostCopyHandler())
	addEndpoint(router, proVRoute, "POST", "/delete/{id}", entryViewCtrl.PostDeleteHandler())
	addEndpoint(router, proVRoute, "GET", "/search", entryViewCtrl.GetSearchHandler())
	addEndpoint(router, proVRoute, "POST", "/search", entryViewCtrl.PostSearchHandler())
	addEndpoint(router, proVRoute, "GET", "/overview", entryViewCtrl.GetOverviewHandler())
	addEndpoint(router, proVRoute, "POST", "/overview", entryViewCtrl.PostOverviewHandler())
	addEndpoint(router, proVRoute, "GET", "/overview/export", entryViewCtrl.GetOverviewExportHandler())
	// Add view resource endpoints
	fileSrv := http.FileServer(http.Dir("./resources"))
	router.Handle("/css/{name}", fileSrv).Methods("GET")
	router.Handle("/img/{name}", fileSrv).Methods("GET")
	router.Handle("/font/{name}", fileSrv).Methods("GET")
	router.Handle("/js/{name}", fileSrv).Methods("GET")

	// Register router
	http.Handle("/", router)

	// Start HTTP server
	log.Infof("Listen on port '%d'.", conf.ServerPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", conf.ServerPort), nil)
	if err != nil {
		log.Fatalf("Could not start server! (Error: %s)", err)
	}
}

func addEndpoint(r *mux.Router, m *negroni.Negroni, method string, path string, hf http.HandlerFunc) {
	r.Handle(path, createHandler(m, hf)).Methods(method)
}

func createHandler(r *negroni.Negroni, h http.HandlerFunc) http.Handler {
	nr := r.With()
	nr.UseHandlerFunc(h)
	return nr
}

func getRootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, constant.PathDefault, http.StatusFound)
	}
}
