package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/matt-FFFFFF/alz-status-badge/badge"
)

type ApprovedVariants map[string]interface{}

const ApprovedVariantsUrl = "https://raw.githubusercontent.com/matt-FFFFFF/alz-status-badge/main/data/approved-variants.json"

// main starts the server.
func main() {
	listenAddress := ":8080"
	av := make(ApprovedVariants)
	badgeApiFunc := badgeApi(av)
	http.HandleFunc("/api/badge", badgeApiFunc)
	log.Printf("About to listen on %s", listenAddress)
	go updateApprovedVariants(av, ApprovedVariantsUrl)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

// badgeApi returns a http.HandlerFunc for the given approved variants type supplied.
func badgeApi(av ApprovedVariants) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header["X-Request-Id"][0]
		requestId = strings.Replace(requestId, "\n", "", -1)
		requestId = strings.Replace(requestId, "\r", "", -1)
		log.Printf("Received request from: %s (%s)", r.RemoteAddr, requestId)

		variant := r.URL.Query().Get("variant")
		if variant == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing variant parameter."))
			log.Printf("Missing variant parameter (%s)", requestId)
			return
		}

		re := regexp.MustCompile(`^\w{1,32}$`)

		if ok := re.Match([]byte(variant)); !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Invalid variant parameter '%s'. Must be 32 characters or less. a-z, A-Z, 0-9 only.", variant)))
			log.Printf("Invalid variant parameter '%s' (%s)", variant, requestId)
			return
		}

		approved := checkVariant(av, variant)
		log.Printf("Variant: %s approval is %t (%s)", variant, approved, requestId)

		badge, err := badge.MakeShieldsioBadge(variant, approved)
		if err != nil {
			log.Printf("Badge error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(badge)

		log.Printf("End request from: %s (%s)", r.RemoteAddr, requestId)
	}
}

// checkVariant returns the status of the given variant.
func checkVariant(av ApprovedVariants, variant string) bool {
	_, exists := av[variant]
	if !exists {
		return false
	}
	return true
}
