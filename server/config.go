package main

import (
	"flag"
)

type flags struct {
	Name string
	Port string
	DB   string
}

var config flags

func init() {
	flag.StringVar(&config.Name, "project name", "db-forum", "set name of project")
	flag.StringVar(&config.Port, "port", ":5000", "service port")
	flag.StringVar(&config.DB, "database DSN", "user=docker password=docker dbname=docker sslmode=disable", "DSN for database")
}
