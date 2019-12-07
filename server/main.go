package main

import (
	"db-forum/database"
	"db-forum/router"
	"flag"
	"log"

	"github.com/valyala/fasthttp"
)

func main() {
	flag.Parse()
	if err := database.InitDB(config.DB); err != nil {
		log.Println("can't init DB", err.Error())
		return
	}
	r := router.CreateRouter()
	log.Println("starting server on " + config.Port)
	log.Fatal(fasthttp.ListenAndServe(config.Port, r.Handler))
}
