package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	Port := os.Getenv("PORT")
	if Port == "" {
		log.Fatal("PORT is not defined in the environment file")
	}
	fmt.Println("Starting server on :" + Port)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+Port, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
