package main

import (
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
}
