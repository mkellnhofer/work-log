package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"kellnhofer.com/work-log/config"
	"kellnhofer.com/work-log/constant"
	"kellnhofer.com/work-log/log"
)

func main() {
	// Load config
	conf := config.LoadConfig()

	// Set logging level
	log.SetLevel(conf.LogLevel)

	log.Infof("Starting Work Log server %s.", constant.AppVersion)

	// Create initializer
	init := NewInitializer(conf)

	// Open and create/update database
	db := init.GetDb()
	db.OpenDb()
	defer db.CloseDb()
	db.UpdateDb()

	// Create router
	router := mux.NewRouter().StrictSlash(true)

	// Create public middleware route
	pubRoute := negroni.New()
	pubRoute.Use(init.GetErrorMiddleware())
	pubRoute.Use(init.GetSessionMiddleware())
	// Create protected middleware route
	proRoute := negroni.New()
	proRoute.Use(init.GetErrorMiddleware())
	proRoute.Use(init.GetSessionMiddleware())
	proRoute.Use(init.GetAuthMiddleware())

	// Add public endpoints
	addEndpoint(router, pubRoute, "GET", "/error", init.GetErrorController().GetErrorHandler())
	addEndpoint(router, pubRoute, "GET", "/login", init.GetAuthController().GetLoginHandler())
	addEndpoint(router, pubRoute, "POST", "/login", init.GetAuthController().PostLoginHandler())
	// TODO!!!
	// Add protected endpoints
	addEndpoint(router, proRoute, "GET", "/logout", init.GetAuthController().GetLogoutHandler())
	// TODO!!!
	// Add resource endpoints
	// TODO!!!

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
