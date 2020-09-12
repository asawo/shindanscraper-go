package slackbot

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"shindanscraper-go/scraper"

	"github.com/kelseyhightower/envconfig"
	"github.com/slack-go/slack"
)

type config struct {
	SigningSecret string `envconfig:"SIGNING_SECRET" required:"true"`
}

// SlashCommandHandler handles the slash command from shindanbot
func SlashCommandHandler(w http.ResponseWriter, r *http.Request) {
	var env config
	if err := envconfig.Process("", &env); err != nil {
		log.Println(err)
	}

	signingSecret := env.SigningSecret

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

		// jsonShindan, _ := json.MarshalIndent(shindans, "", "    ")
		// fmt.Println(string(jsonShindan))

		msgBlock := CreateBlock(shindans)

		w.Header().Set("Content-Type", "application/json")
		w.Write(msgBlock)
		return

	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// CreateBlock creates slack message block from map object
func CreateBlock(shindan map[int]scraper.ShindanObj) []byte {

	divSection := slack.NewDividerBlock()

	headerText := slack.NewTextBlockObject("mrkdwn", "*Header*", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil, nil)

	rankOne := slack.NewTextBlockObject("mrkdwn", "test text", false, false)
	rankOneSection := slack.NewSectionBlock(rankOne, nil, nil, nil)

	// Build Message with blocks created above
	msg := slack.NewBlockMessage(
		headerSection,
		divSection,
		rankOneSection,
		divSection,
	)

	b, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		log.Printf("Error marshalling json: %v", err)
	}

	return b
}
