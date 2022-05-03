package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/matt-FFFFFF/alz-status-badge/badge"
)

type ApprovedVariants map[string]interface{}

const ApprovedVariantsUrl = "https://raw.githubusercontent.com/matt-FFFFFF/alz-status-badge/main/data/approved-variants.json"

// main starts the server.
func main() {
	listenAddress := ":8080"
	http.HandleFunc("/api/badge", badgeApi)
	log.Printf("About to listen on %s", listenAddress)
	av := make(ApprovedVariants)
	go updateApprovedVariants(av)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

// badgeApi returns the badge for the given variant supplied in the query.
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
		log.Printf("Badge error: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(badge)

	log.Printf("End request from: %s (%s)", r.RemoteAddr, r.Header["X-Request-Id"])
}

// checkVariant returns the status of the given variant.
func checkVariant(variant string) bool {
	return variant == "canadapubsec"
}

func updateApprovedVariants(av ApprovedVariants, url string) {
	for {
		log.Printf("Checking approved variants")
		vl, err := getApprovedVariants(&url)
		if err != nil {
			log.Printf("Error downloading approved variants: %s", err.Error())
		}

		av["canadapubsec"] = nil
		var vs string
		for v := range av {
			vs = vs + (fmt.Sprintf("%s, ", v))
		}
		log.Printf("Approved variants: %s", vs)
		time.Sleep(time.Second * 15)
	}
}

func getApprovedVariants(url *string) ([]string, error) {
	resp, err := http.Get(*url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	v, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	vl := make([]string, 0)
	json.Unmarshal(v, &vl)

	return vl, nil
}
