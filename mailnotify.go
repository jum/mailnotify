package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gregdel/pushover"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: %s api_key user_key\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	app := pushover.New(os.Args[1])

	// Create a new recipient
	recipient := pushover.NewRecipient(os.Args[2])

	// Create the message to send
	message := &pushover.Message{
		Message: os.Getenv("SUBJECT"),
		Title:   os.Getenv("FROM"),
	}
	msgid := os.Getenv("MSGID")
	if len(msgid) > 0 {
		url := (&url.URL{
			Scheme: "message",
			Path:   msgid,
		}).String()
		message.URL = url
		message.URLTitle = "Apple Mail"
	}

	// Send the message to the recipient
	response, err := app.SendMessage(message, recipient)
	if err != nil {
		log.Panic(err)
	}
	if false {
		log.Printf("pushover resp %v", response)
	}
}
