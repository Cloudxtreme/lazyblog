package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bentranter/lazyblog"
)

func main() {
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
		log.Fatalln("Please choose either setup or serve")
	}
}
