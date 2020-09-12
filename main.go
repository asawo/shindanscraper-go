package main

import (
	"fmt"
	"net/http"
	"shindanscraper-go/slackbot"
)

func main() {
	http.HandleFunc("/slash", slackbot.SlashCommandHandler)

	fmt.Println("[INFO] Server listening on http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
