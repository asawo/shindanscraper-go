package slackbot

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"shindanscraper-go/scraper"

	"github.com/nlopes/slack"
)

// SlashCommandHandler handles the slash command from shindanbot
func SlashCommandHandler(w http.ResponseWriter, r *http.Request) {
	signingSecret := os.Getenv("SIGNING_SECRET")

	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/shindan":
		shindans, err := scraper.GetShindans("https://shindanmaker.com/c/list?mode=hot")
		if err != nil {
			log.Println(err)
		}
		shindanData, err := json.MarshalIndent(shindans, "", "  ")
		if err != nil {
			log.Println(err)
		}
		jsonShindan := string(shindanData)
		fmt.Print(jsonShindan)

		w.Header().Set("Content-Type", "application/json")
		w.Write(shindanData)

	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
