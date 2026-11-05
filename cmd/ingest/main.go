package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tolaniverse/streamforge/internal/router"
)

func main() {
	r := chi.NewRouter()
	router.Routes(r)

	port := ":3000"
	fmt.Printf("Listening at port %s", port)
	http.ListenAndServe(":3000", r)
}
