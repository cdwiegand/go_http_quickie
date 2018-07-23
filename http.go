package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func handleHttpRequest(w http.ResponseWriter, r *http.Request) {
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	start := time.Now()
	bodyIn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	log.Printf("%s %s %s %s", r.Proto, r.Method, r.RemoteAddr, r.URL)

	// handle POST/other verbs, please? FIXME
	req, err := http.NewRequest(r.Method, "http://localhost"+r.URL.Path, bytes.NewBuffer(bodyIn))
	if err != nil {
		log.Println("ERROR: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// pass thru orig request headers to backend
	for k, v := range r.Header {
		for _, v2 := range v {
			req.Header.Add(k, v2)
			log.Printf(" <- %q: %q\n", k, v2)
		}
	}

	resp, err := client.Do(req) // once you pass this line you MUST defer a Body.Close()!
	if err != nil {
		log.Println("ERROR: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		// MUST defer body close or we leak resources!
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ERROR: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get response headers and return...
	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
		log.Printf(" -> %q: %q\n", k, v)
	}
	// now set some of our own... these are SET as they replace any existing value
	w.Header().Set("X-Hello", "Darkness, my old friend")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	timeRan := time.Since(start)
	w.Header().Set("X-Timing", fmt.Sprintf("%s", timeRan))
	w.Header().Set("X-Upstream-Proto", resp.Proto)
	w.Header().Set("X-Powered-By", "Nope/2.0")
	w.Header().Set("Server", "Nope/2.0")

	// NO MORE SETTING HEADERS!!!
	w.WriteHeader(resp.StatusCode)
	w.Write([]byte(body))
	log.Println("Time: " + fmt.Sprintf("%s", timeRan))
}

func handleHttpHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.RemoteAddr, r.URL)
		f(w, r)
	}
}

func main() {
	srv := &http.Server{
		Addr:           ":8080",
		Handler:        handleHttpHandler(handleHttpRequest),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, /* ??? FIXME */
	}
	srv.SetKeepAlivesEnabled(true)
	log.Println("Starting server on port 8080")
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
