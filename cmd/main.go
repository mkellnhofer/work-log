package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"kellnhofer.com/work-log/pkg/config"
	"kellnhofer.com/work-log/pkg/constant"
	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/pkg/log"
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
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover(), createLoggerMiddleware())

	// Add view handlers
	addViewHandlers(init, e)
	// Add API handlers
	addApiHandlers(init, e)
	// Add Swagger UI handlers
	addSwaggerUiHandlers(e)

	// Start HTTP server
	log.Infof("Listen on port '%d'.", conf.ServerPort)
	err := e.Start(fmt.Sprintf(":%d", conf.ServerPort))
	if err != nil {
		log.Fatalf("Could not start server! (Error: %s)", err)
	}
}

func addViewHandlers(init *Initializer, e *echo.Echo) {
	// Create public middleware
	pubRoute := []echo.MiddlewareFunc{
		init.GetTransactionMiddleware().CreateHandler,
		init.GetErrorViewMiddleware().CreateHandler,
		init.GetSessionViewMiddleware().CreateHandler,
		init.GetSecurityViewMiddleware().CreateHandler}
	// Create protected middleware
	proRoute := []echo.MiddlewareFunc{
		init.GetTransactionMiddleware().CreateHandler,
		init.GetErrorViewMiddleware().CreateHandler,
		init.GetSessionViewMiddleware().CreateHandler,
		init.GetSecurityViewMiddleware().CreateHandler,
		init.GetAuthCheckViewMiddleware().CreateHandler}

	// Get controllers
	errCtrl := init.GetErrorViewController()
	authCtrl := init.GetAuthViewController()
	entryCtrl := init.GetEntryViewController()
	logCtrl := init.GetLogViewController()
	overviewCtrl := init.GetOverviewViewController()
	searchCtrl := init.GetSearchViewController()

	// Register public handlers
	e.GET("/", getRootHandler())
	e.GET("/error", errCtrl.GetErrorHandler(), pubRoute...)
	e.GET("/login", authCtrl.GetLoginHandler(), pubRoute...)
	e.POST("/login", authCtrl.PostLoginHandler(), pubRoute...)
	// Register protected handlers
	e.GET("/logout", authCtrl.GetLogoutHandler(), proRoute...)
	e.GET("/log", logCtrl.GetLogHandler(), proRoute...)
	e.GET("/search", searchCtrl.GetSearchHandler(), proRoute...)
	e.POST("/search", searchCtrl.PostSearchHandler(), proRoute...)
	e.GET("/overview", overviewCtrl.GetOverviewHandler(), proRoute...)
	e.GET("/overview/export", overviewCtrl.GetOverviewExportHandler(), proRoute...)
	e.GET("/create", entryCtrl.GetCreateHandler(), proRoute...)
	e.POST("/create", entryCtrl.PostCreateHandler(), proRoute...)
	e.GET("/copy/:id", entryCtrl.GetCopyHandler(), proRoute...)
	e.GET("/edit/:id", entryCtrl.GetEditHandler(), proRoute...)
	e.POST("/edit/:id", entryCtrl.PostEditHandler(), proRoute...)
	e.GET("/delete/:id", entryCtrl.GetDeleteHandler(), proRoute...)
	e.POST("/delete/:id", entryCtrl.PostDeleteHandler(), proRoute...)
	e.POST("/cancel", entryCtrl.PostCancelHandler(), proRoute...)
	// Register resource handlers
	e.Static("/css/", "static/resources/css")
	e.Static("/img/", "static/resources/img")
	e.Static("/font/", "static/resources/font")
	e.Static("/js/", "static/resources/js")
}

func getRootHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.Redirect(http.StatusFound, constant.ViewPathDefault)
	}
}

func addApiHandlers(init *Initializer, e *echo.Echo) {
	// Create API group
	g := e.Group(constant.ApiPath)

	// Add protected middleware
	g.Use(createCorsMiddleware(),
		init.GetTransactionMiddleware().CreateHandler,
		init.GetErrorApiMiddleware().CreateHandler,
		init.GetSecurityApiMiddleware().CreateHandler,
		init.GetAuthCheckApiMiddleware().CreateHandler)

	// Get controllers
	entryCtrl := init.GetEntryApiController()
	userCtrl := init.GetUserApiController()

	// Register protected handlers
	g.GET("/entries", entryCtrl.GetEntriesHandler())
	g.POST("/entries", entryCtrl.CreateEntryHandler())
	g.GET("/entries/:id", entryCtrl.GetEntryHandler())
	g.PUT("/entries/:id", entryCtrl.UpdateEntryHandler())
	g.DELETE("/entries/:id", entryCtrl.DeleteEntryHandler())
	g.GET("/entry_types", entryCtrl.GetEntryTypesHandler())
	g.GET("/entry_activities", entryCtrl.GetEntryActivitiesHandler())
	g.POST("/entry_activities", entryCtrl.CreateEntryActivityHandler())
	g.PUT("/entry_activities/:id", entryCtrl.UpdateEntryActivityHandler())
	g.DELETE("/entry_activities/:id", entryCtrl.DeleteEntryActivityHandler())
	g.GET("/user", userCtrl.GetCurrentUserHandler())
	g.PUT("/user/password", userCtrl.UpdateCurrentUserPasswordHandler())
	g.GET("/user/roles", userCtrl.GetCurrentUserRolesHandler())
	g.GET("/users", userCtrl.GetUsersHandler())
	g.POST("/users", userCtrl.CreateUserHandler())
	g.GET("/users/:id", userCtrl.GetUserHandler())
	g.PUT("/users/:id", userCtrl.UpdateUserHandler())
	g.DELETE("/users/:id", userCtrl.DeleteUserHandler())
	g.PUT("/users/:id/password", userCtrl.UpdateUserPasswordHandler())
	g.GET("/users/:id/roles", userCtrl.GetUserRolesHandler())
	g.PUT("/users/:id/roles", userCtrl.UpdateUserRolesHandler())
}

func addSwaggerUiHandlers(e *echo.Echo) {
	e.GET("/api", getSwaggerRootHandler())
	e.Static("/api/", "static/swagger-ui/")
}

func getSwaggerRootHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.Redirect(http.StatusMovedPermanently, "/api/swagger-ui.html")
	}
}

func createLoggerMiddleware() echo.MiddlewareFunc {
	return log.NewLoggerMiddleware().CreateHandler
}

func createCorsMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
}
