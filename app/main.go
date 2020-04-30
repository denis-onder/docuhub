package main

import (
	"os"
	"strconv"

	"github.com/denis-onder/docuhub/app/pkg/server"
)

func main() {
	var port int
	p, exists := os.LookupEnv("PORT")
	if !exists {
		port = 9000
	} else {
		parsed, _ := strconv.Atoi(p)
		port = parsed
	}
	server.Start(port)
}
