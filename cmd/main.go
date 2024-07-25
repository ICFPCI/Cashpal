package main

import (
	"cashpal/api"
	"cashpal/middleware"
	"fmt"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	api.SetupURLs(router)

	middlewares := middleware.ChainMiddleware(
		middleware.Logger,
		middleware.JSON,
	)

	fmt.Println("Backend listening on http://localhost:8000/")

	http.ListenAndServe(":8000", middlewares(router))
}
