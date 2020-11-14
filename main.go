package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/hawwwdi/Twitter_clone/api"
	"github.com/hawwwdi/Twitter_clone/db"
)

func main() {
	dbPort := flag.String("dbp", "6037", "redis server port")
	serverPort := flag.String("sp", "8080", "server port")
	flag.Parse()
	db := db.NewDB("localhost:" + *dbPort)
	api.InitHub(db)
	hub, err := api.GetHub()
	if err != nil {
		panic(err)
	}
	log.Fatalln(http.ListenAndServe("localhost:"+*serverPort, hub.Router))
}
