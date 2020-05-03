package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/denis-onder/docuhub/app/pkg/handlers"
	"github.com/gorilla/mux"
)

// Router
func createRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/all", handlers.GetDocumentaries)

	return router
}

// Start the web server
func Start(port int) {
	router := createRouter()

	// Server init
	p := ":" + strconv.Itoa(port)
	fmt.Printf("Server running!\nhttp://localhost%s/\n", p)
	http.ListenAndServe(p, router)
}
