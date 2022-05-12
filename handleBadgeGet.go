package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/matt-FFFFFF/alz-status-badge/badge"
)

// badgeApi returns a http.HandlerFunc for the given approved variants type supplied.
func (s *server) handleBadgeGet() http.HandlerFunc {
	// Use a closure over the server to return the http.HandlerFunc
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET requests
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		// When running in Azure container apps, we have a request header called X-Request-Id that we can use to improve logging.
		requestId := "none"
		if r, ok := r.Header["X-Request-Id"]; ok {
			requestId = r[0]
		}
		requestId = strings.Replace(requestId, "\n", "", -1)
		requestId = strings.Replace(requestId, "\r", "", -1)

		log.Printf("Received request from: %s (%s)", r.RemoteAddr, requestId)

		if len(s.approvedVariants.list) == 0 {
			log.Printf("Approved variants not yet loaded, return Internal Server Error (%s)", requestId)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Get the variant from the URL
		variant := r.URL.Query().Get("variant")
		if variant == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing variant parameter."))
			log.Printf("Missing variant parameter (%s)", requestId)
			return
		}

		// variant must be a-z, A-Z, 0-9 with a length of 1-32 characters
		re := regexp.MustCompile(`^\w{1,32}$`)

		if ok := re.Match([]byte(variant)); !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Invalid variant parameter '%s'. Must be 32 characters or less. a-z, A-Z, 0-9 only.", variant)))
			log.Printf("Invalid variant parameter (%s)", requestId)
			return
		}

		// Check if the variant is approved
		approved := checkVariant(s.approvedVariants, variant)
		log.Printf("Variant: %s approval is %t (%s)", variant, approved, requestId)

		// Make new badge
		br := badge.NewBadgeRequest()
		br.Style = badge.ForTheBadge
		br.Label = "ALZ-VARIANT"
		br.Message = variant
		br.Color = "success"
		if !approved {
			br.Color = "critical"
			br.Message = "NOT APPROVED"
		}
		badge, err := badge.Get(br)

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
func checkVariant(av *ApprovedVariants, variant string) bool {
	av.mu.RLock()
	defer av.mu.RUnlock()
	_, exists := av.list[variant]
	if !exists {
		return false
	}
	return true
}
