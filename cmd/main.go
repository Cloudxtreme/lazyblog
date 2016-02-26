package main

import (
	"flag"
	"net/http"

	"github.com/bentranter/lazyblog"
)

func main() {
	username := flag.String("username", "", "the username you'll login with")
	password := flag.String("password", "", "your login password")
	flag.Parse()
	lazyblog.Setup(*username, *password)

	defer lazyblog.DefaultStore.Close()
	http.ListenAndServe(":3000", lazyblog.Router)
}
