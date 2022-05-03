package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/matt-FFFFFF/alz-status-badge/badge"
)

type AppConfig struct {
	ApprovedVariantsUrl             *url.URL         `json:"approvedVariantsUrl"`
	ApprovedVariants                ApprovedVariants `json:"approvedVariants"`
	ApprovedVariantsRefreshInterval time.Duration    `json:"approvedVariantsRefreshInterval"`
	ListenAddress                   string           `json:"listenAddress"`
}

type ApprovedVariants map[string]interface{}

const DefaultApprovedVariantsUrl = "https://raw.githubusercontent.com/matt-FFFFFF/alz-status-badge/main/data/approved-variants.json"

// main starts the server.
func main() {
	// Set up the application configuration defaults and read environment variables
	u, err := url.Parse(DefaultApprovedVariantsUrl)
	if err != nil {
		log.Fatalf("Failed to parse default approved variants url: %s", err.Error())
	}

	config := AppConfig{
		ApprovedVariantsUrl:             u,
		ApprovedVariants:                make(ApprovedVariants),
		ApprovedVariantsRefreshInterval: time.Minute * 15,
		ListenAddress:                   ":8080",
	}

	if s, ok := os.LookupEnv("ALZSTATUSBADGE_LISTEN_ADDRESS"); ok {
		config.ListenAddress = s
	}

	if s, ok := os.LookupEnv("ALZSTATUSBADGE_APPROVED_VARIANTS_URL"); ok {
		u, err := url.Parse(s)
		if err != nil {
			log.Fatalf("Invalid ALZSTATUSBADGE_APPROVED_VARIANTS_URL: %s", err.Error())
		}
		config.ApprovedVariantsUrl = u
	}

	if s, ok := os.LookupEnv("ALZSTATUSBADGE_APPROVED_VARIANTS_REFRESH_INTERVAL"); ok {
		d, err := time.ParseDuration(s)
		if err != nil {
			log.Fatalf("Invalid ALZSTATUSBADGE_APPROVED_VARIANTS_REFRESH_INTERVAL: %s", err.Error())
		}
		config.ApprovedVariantsRefreshInterval = d
	}

	http.HandleFunc("/api/badge", badgeApi(&config))
	log.Printf("About to listen on %s", config.ListenAddress)
	go updateApprovedVariants(&config)
	log.Fatal(http.ListenAndServe(config.ListenAddress, nil))
}

// badgeApi returns a http.HandlerFunc for the given approved variants type supplied.
func badgeApi(config *AppConfig) http.HandlerFunc {
	// Use a closure over the AppConfig to return the http.HandlerFunc
	return func(w http.ResponseWriter, r *http.Request) {
		// When running in Azure container apps, we have a request header called X-Request-Id that we can use to improve logging.
		requestId := "none"
		if r, ok := r.Header["X-Request-Id"]; ok {
			requestId = r[0]
		}
		requestId = strings.Replace(requestId, "\n", "", -1)
		requestId = strings.Replace(requestId, "\r", "", -1)

		log.Printf("Received request from: %s (%s)", r.RemoteAddr, requestId)

		if len(config.ApprovedVariants) == 0 {
			log.Printf("Approved variants not yet loaded, return Internal Server Error (%s)", requestId)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

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

		approved := checkVariant(config.ApprovedVariants, variant)
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
