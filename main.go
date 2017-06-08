package main

import (
	"net/http"

	"fmt"

	"github.com/gorilla/mux"
)

func main() {
	InitRender()

	router := mux.NewRouter()
	router.PathPrefix("/statics/").Handler(http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics"))))
	/*router.Methods("GET").Path("/").Handler(http.HandlerFunc(homePageHandler))*/
	for _, i := range routes {
		router.Methods(i.Method).Path(i.Pattern).Handler(http.HandlerFunc(i.Function))
	}
	fmt.Println("server runnig on port 8080")
	http.ListenAndServe(":8080", router)
}
