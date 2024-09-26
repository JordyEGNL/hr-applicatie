package main

import (
	"os"
	"strconv"
)

// API configuration
var api = struct {
	port            string
	debug           bool
	defaultuser     string
	defaultpassword string
}{
	port:            "5000",
	debug:           true,
	defaultuser:     "admin@example.com",
	defaultpassword: "Admin01!",
}

var dbcred = struct {
	host     string
	dbname   string
	port     string
	user     string
	password string
}{
	host:     "db",
	dbname:   "hrsysteem",
	port:     "3306",
	user:     "dbuser",
	password: "Admin01!",
}

// Get env from docker and set it to the config
func getEnv() {
	if os.Getenv("API_PORT") != "" {
		api.port = os.Getenv("API_PORT")
	}
	if os.Getenv("API_DEBUG") != "" {
		api.debug, _ = strconv.ParseBool(os.Getenv("API_DEBUG"))
	}
	if os.Getenv("DB_HOST") != "" {
		dbcred.host = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_NAME") != "" {
		dbcred.dbname = os.Getenv("DB_NAME")
	}
	if os.Getenv("DB_PORT") != "" {
		dbcred.port = os.Getenv("DB_PORT")
	}
	if os.Getenv("DB_USER") != "" {
		dbcred.user = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_PASSWORD") != "" {
		dbcred.password = os.Getenv("DB_PASSWORD")
	}
	if os.Getenv("API_DEFAULT_USER") != "" {
		api.defaultuser = os.Getenv("API_DEFAULT_USER")
	}
	if os.Getenv("API_DEFAULT_PASSWORD") != "" {
		api.defaultpassword = os.Getenv("API_DEFAULT_PASSWORD")
	}
}
