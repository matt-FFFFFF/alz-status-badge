package main

import (
	"log"
	"net/http"
	"os"

	svg "github.com/ajstarks/svgo"
)

func main() {
	listenAddress := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddress = ":" + val
	}
	http.HandleFunc("/api/badge", badge)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddress, listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

func badge(w http.ResponseWriter, r *http.Request) {
	log.Printf("HTTP referrer: %s", r.Referer())
	w.Header().Set("Content-Type", "image/svg+xml")
	s := svg.New(w)
	s.Start(250, 100)
	s.Text(10, 20, r.Referer(), "font-family=\"Verdana\" font-size=\"12\" fill=\"blue\"")
	s.End()
}
