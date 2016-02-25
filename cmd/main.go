package main

import (
	"net/http"

	"github.com/bentranter/lazyblog"
)

func main() {
	http.ListenAndServe(":3000", lazyblog.Router)
}
