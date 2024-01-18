package config

import (
	"github.com/go-ini/ini"

	"kellnhofer.com/work-log/pkg/log"
)

// Config stores the application's configuration.
type Config struct {
	ServerPort  int
	LogLevel    string
	DbHost      string
	DbPort      int
	DbScheme    string
	DbUsername  string
	DbPassword  string
	LocLanguage string
}

// LoadConfig loads the configuration from "/config/config.ini".
func LoadConfig() *Config {
	return loadConfig("config/config.ini")
}

// LoadTestConfig loads the test configuration from "/config/config_test.ini".
func LoadTestConfig() *Config {
	return loadConfig("config/config_test.ini")
}

func loadConfig(name string) *Config {
	cfg, err := ini.Load(name)
	if err != nil {
		log.Fatalf("Config file missing! %s", err)
	}

	serverPort := getIntValue(cfg, "server", "port")

	logLevel := getStringValue(cfg, "log", "level")

	dbHost := getStringValue(cfg, "database", "host")
	dbPort := getIntValue(cfg, "database", "port")
	dbScheme := getStringValue(cfg, "database", "scheme")
	dbUsername := getStringValue(cfg, "database", "username")
	dbPassword := getStringValue(cfg, "database", "password")

	locLanguage := getStringValue(cfg, "localization", "language")

	return &Config{serverPort, logLevel, dbHost, dbPort, dbScheme, dbUsername, dbPassword,
		locLanguage}
}

func getStringValue(file *ini.File, secName string, keyName string) string {
	return getKey(file, secName, keyName).String()
}

func getIntValue(file *ini.File, secName string, keyName string) int {
	val, err := getKey(file, secName, keyName).Int()
	if err != nil {
		log.Fatalf("Config file has invalid value for key '%s'!", keyName)
	}
	return val
}

func getKey(file *ini.File, secName string, keyName string) *ini.Key {
	sec, err := file.GetSection(secName)
	if err != nil {
		log.Fatalf("Config file missing section '%s'!", secName)
	}

	if !sec.HasKey(keyName) {
		log.Fatalf("Config file missing key '%s'!", keyName)
	}

	return sec.Key(keyName)
}
