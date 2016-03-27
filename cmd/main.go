package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bentranter/lazyblog"
)

const usage = `
NAME:
   lazyblog - viral proof personal blogging platform

USAGE:
   lazyblog <command>

VERSION:
   0.1.0

COMMANDS:
   setup    Create a new login and password combination
   start    Start the server
   help     Show this help information

`

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "setup":
			defer lazyblog.DefaultStore.Close()
			lazyblog.Setup()
		case "start":
			defer lazyblog.DefaultStore.Close()
			numUsers, err := lazyblog.NumUsers()
			if err != nil {
				panic(err)
			}
			if numUsers < 1 {
				log.Fatalln("Please run setup before running start")
				return
			}

			if os.Getenv("LAZYBLOG_ENV") == "dev" {
				log.Fatalln(http.ListenAndServe(":3000", lazyblog.Router))
			}
			log.Fatalln(http.ListenAndServe(":80", lazyblog.Router))
		default:
			fmt.Println(usage)
		}
	} else {
		fmt.Println(usage)
	}
}
