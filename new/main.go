package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/bentranter/lazyblog/handler"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

var (
	httpAddr  = flag.String("http", "8000", "HTTP service address")
	pprofAddr = flag.String("pprof", "3000", "pprof service address")
)

func main() {
	flag.Parse()

	r := fasthttprouter.New()
	errCh := make(chan error, 10)

	log.Printf("Starting HTTP server on port %s\n", *httpAddr)
	log.Printf("pprof server on port %s\n", *pprofAddr)

	r.GET("/api", handler.Info)
	r.GET("/api/posts", handler.GetAllPosts)
	r.GET("/api/posts/:id", handler.GetPostJSON)
	r.POST("/api/posts", handler.SetPost)
	r.ServeFiles("/static/*filepath", "./static")

	// pprof server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("go to /pprof/debug"))
	})

	go func() {
		errCh <- fasthttp.ListenAndServe(":"+*httpAddr, r.Handler)
	}()
	go func() {
		errCh <- http.ListenAndServe(":"+*pprofAddr, nil)
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errCh:
			if err != nil {
				log.Fatalf("%s\n", err.Error())
			}
		case s := <-signalCh:
			log.Printf("Captured %v. Exiting...", s)
			os.Exit(0)
		}
	}

}
