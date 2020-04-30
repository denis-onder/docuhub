package server

import (
	"fmt"
	"net/http"
	"strconv"
)

// Start the web server
func Start(port int) {
	p := ":" + strconv.Itoa(port)
	fmt.Printf("Server running!\nhttp://localhost%s/\n", p)
	http.ListenAndServe(p, nil)
}
