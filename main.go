package main

import (
	"net/http"

	"github.com/hawwwdi/Twitter_clone/api"
	_ "github.com/hawwwdi/Twitter_clone/db"
)

func main() {
	mux := api.NewRouter()
	panic(http.ListenAndServe(":8080", mux))
}
