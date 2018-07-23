package main

import (
	"log"
	"net/http"
	"time"
)

var client http.Client

func main() {
	loadConfig()

	client = http.Client{
		Timeout: time.Duration(config.ClientTimeout) * time.Second,
	}

	log.Println("Trying to listen on " + config.ListenAddr)
	srv := &http.Server{
		Addr:           config.ListenAddr,
		Handler:        handleHTTPHandler(handleHTTPRequest),
		ReadTimeout:    time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20, /* ??? FIXME */
	}
	srv.SetKeepAlivesEnabled(true)
	log.Println("Started server on " + srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func handleHTTPHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.RemoteAddr, r.URL)
		f(w, r)
	}
}
