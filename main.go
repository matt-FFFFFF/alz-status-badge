package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/matt-FFFFFF/alz-status-badge/badge"
)

// main starts the server.
func main() {
	listenAddress := ":8080"
	http.HandleFunc("/api/badge", badgeApi)
	log.Printf("About to listen on %s", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

// badge returns the badge for the given status.
func badgeApi(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request from: %s (%s)", r.RemoteAddr, r.Header["X-Request-Id"])

	variant := r.URL.Query().Get("variant")
	if variant == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing variant parameter."))
		log.Printf("Missing variant parameter (%s)", r.Header["X-Request-Id"])
		return
	}

	re := regexp.MustCompile(`^\w{1,32}$`)

	if ok := re.Match([]byte(variant)); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Invalid variant parameter '%s'. Must be 32 characters or less. a-z, A-Z, 0-9 only.", variant)))
		log.Printf("Invalid variant parameter '%s' (%s)", variant, r.Header["X-Request-Id"])
		return
	}

	approved := checkVariant(variant)
	log.Printf("Variant: %s approval is %t (%s)", variant, approved, r.Header["X-Request-Id"])

	badge, err := badge.MakeShieldsioBadge(variant, approved)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(badge)

	//s := svg.New(w)
	//makeBadge(variant, approved, s)

	log.Printf("End request from: %s (%s)", r.RemoteAddr, r.Header["X-Request-Id"])
}

// checkVariant returns the status of the given variant.
func checkVariant(variant string) bool {
	return variant == "canadapubsec"
}
