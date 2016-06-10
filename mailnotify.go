package main

import (
    "fmt"
    "bitbucket.org/kisom/gopush/pushover"
    "os"
    "path/filepath"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Printf("usage: %s api_key user_key\n", filepath.Base(os.Args[0]))
        os.Exit(1)
    }

    pushover.Verbose = false;
    identity := pushover.Authenticate(
        os.Args[1],
        os.Args[2],
    )

    sent := pushover.Notify_titled(identity, os.Getenv("SUBJECT"), os.Getenv("FROM"))
    if !sent {
        fmt.Println("[!] notification failed.")
        os.Exit(1)
    }
}
