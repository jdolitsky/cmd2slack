package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

type Payload struct {
	Text      string `json:"text"`
	Username  string `json:"username,omitempty"`
	Channel   string `json:"channel,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
	IconURL   string `json:"icon_url,omitempty"`
}

func main() {
	payload := &Payload{}

	var (
		hook     = flag.String("hook", "", "Slack Incoming Webhook URL")
		channel  = flag.String("channel", "", "Channel where to post the output")
		emoji    = flag.String("emoji", "", "Emoji to use")
		username = flag.String("username", "", "Username")
		icon     = flag.String("icon", "", "URL of icon to use")
	)

	flag.Parse()

	if *hook == "" {
		fmt.Fprintln(os.Stderr, "-hook is required")
		os.Exit(1)
	}

	if *channel != "" {
		payload.Channel = *channel
	}

	if *username != "" {
		payload.Username = *username
	}

	if *emoji != "" {
		payload.IconEmoji = *emoji
	}

	if *icon != "" {
		payload.IconURL = *icon
	}

	args := flag.Args()

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "expected a command")
		os.Exit(2)
	}

	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}

	payload.Text = "```\n" + string(out) + "\n```"

	body := new(bytes.Buffer)

	json.NewEncoder(body).Encode(payload)

	res, err := http.Post(*hook, "application/json", body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
	}

	if res.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "posting to slack failed with", res.Status)
		// TODO this should be more expressive
		os.Exit(5)
	}
}
