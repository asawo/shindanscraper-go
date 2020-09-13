package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/slash", SlashCommandHandler)

	fmt.Println("[INFO] Server listening on http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
