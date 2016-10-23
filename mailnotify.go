package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/quotedprintable"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gregdel/pushover"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: %s api_key user_key\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	app := pushover.New(os.Args[1])

	// Create a new recipient
	recipient := pushover.NewRecipient(os.Args[2])

	subject, err := decodeHeader(os.Getenv("SUBJECT"))
	if err != nil {
		log.Panic(err)
	}
	from, err := decodeHeader(os.Getenv("FROM"))
	if err != nil {
		log.Panic(err)
	}

	// Create the message to send
	message := &pushover.Message{
		Message: subject,
		Title:   from,
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

func decoder(encoding string) (*encoding.Decoder, error) {
	if strings.ToUpper(encoding) == "UTF-8" {
		return unicode.UTF8.NewDecoder(), nil
	} else if strings.ToUpper(encoding) == "ISO-8859-1" {
		return charmap.ISO8859_1.NewDecoder(), nil
	} else {
		return nil, fmt.Errorf("Unknown encoding")
	}
}

func decodeHeader(str string) (string, error) {
	re := regexp.MustCompile(`\=\?(?P<charset>.*?)\?(?P<encoding>.*)\?(?P<body>.*?)\?(.*?)\=`)

	matches := re.FindAllStringSubmatch(str, -1)
	if len(matches) == 0 {
		return str, nil
	}

	for _, match := range matches {
		var r io.Reader = strings.NewReader(match[3])

		if match[2] == "Q" {
			r = quotedprintable.NewReader(r)
		} else if match[2] == "B" {
			r = base64.NewDecoder(base64.StdEncoding, r)
		}

		if d, err := decoder(match[1]); err == nil {
			r = d.Reader(r)
		}

		if val, err := ioutil.ReadAll(r); err == nil {
			str = strings.Replace(str, match[0], string(val), -1)
		} else if err != nil {
			fmt.Println(err.Error())
			continue
		}

	}

	return str, nil
}
