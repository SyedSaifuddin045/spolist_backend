package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SyedSaifuddin045/Spolist_Backend/song"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	godotenv.Load()

	Port := os.Getenv("PORT")
	if Port == "" {
		log.Fatal("PORT is not defined in the environment file")
	}
	fmt.Println("Starting server on :" + Port)
	// Define CORS options
	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // You can specify specific origins instead of allowing all with "*"
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})
	http.HandleFunc("/", handler)
	http.Handle("/download_song", corsOptions.Handler(http.HandlerFunc(song.HandleSongDownload)))
	http.ListenAndServe(":"+Port, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
