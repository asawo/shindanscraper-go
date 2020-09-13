package slackbot

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"shindanscraper-go/scraper"
	"strings"

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
			log.Printf("[ERROR] Error getting Shindans: %v", err)
		}

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

	// Create listSection from ShindanObj
	var t []string
	t = append(t, "🔥 *Top 10 Hottest Shindans* 🔥\n")
	for i := 1; i <= 10; i++ {
		t = append(t, fmt.Sprintf("\n%v. <%v|%v>", i, shindan[i].URL, shindan[i].Title))
	}
	sh := strings.Join(t, "")

	shindanList := slack.NewTextBlockObject("mrkdwn", sh, false, false)
	listSection := slack.NewSectionBlock(shindanList, nil, nil)

	msg := slack.NewBlockMessage(
		divSection,
		listSection,
		divSection,
	)

	b, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		fmt.Println(err)
	}

	return b
}
