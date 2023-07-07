package main

import (
	"flag"
	"net/http"

	"github.com/LuisFlahan4051/basic-counter-app/database"
	"github.com/gorilla/mux"
)

var (
	port *string
	URLs []string
)

func main() {
	initFlags()
	database.InitDatabaseIfNotExists()

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	// RUN SERVER
	http.ListenAndServe(":"+*port, router)
}

func initFlags() {
	port = flag.String("port", "8080", "Port to use")

	flag.Parse()
}
