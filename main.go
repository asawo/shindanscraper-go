package main

import (
	"fmt"
	"log"
	"net/http"
	"shindanscraper-go/slackbot"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/slash", slackbot.SlashCommandHandler)

	fmt.Println("[INFO] Server listening on http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
