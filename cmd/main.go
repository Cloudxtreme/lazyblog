package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/bentranter/lazyblog"
)

func main() {
	username := flag.String("username", "", "the username you'll login with")
	password := flag.String("password", "", "your login password")
	flag.Parse()
	lazyblog.Setup(*username, *password)

	defer lazyblog.DefaultStore.Close()

	if os.Getenv("LAZYBLOG_ENV") == "dev" {
		log.Fatalln(http.ListenAndServe(":3000", lazyblog.Router))
	}
	log.Fatalln(http.ListenAndServe(":80", lazyblog.Router))
}
