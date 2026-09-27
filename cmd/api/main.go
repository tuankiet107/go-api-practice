package main

import (
	"log"
	"net/http"

	"go-api-practice/internal/router"
)

func main() {
	r := router.New()

	log.Println("server is listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
