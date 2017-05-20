package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	InitRender()
	updateUsers()
	updateTravels()

	router := mux.NewRouter()
	router.PathPrefix("/statics/").Handler(http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics"))))
	/*router.Methods("GET").Path("/").Handler(http.HandlerFunc(homePageHandler))*/
	for _, i := range routes {
		router.Methods(i.Method).Path(i.Pattern).Handler(http.HandlerFunc(i.Function))
	}
	http.ListenAndServe(":8080", router)
}
