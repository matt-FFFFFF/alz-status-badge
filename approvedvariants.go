package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"time"
)

// updateApprovedVariants updates the approved variants map every 15 minutes. Use a goroutine to call this function as it will run forever.
func updateApprovedVariants(s *server) {
	for {
		log.Printf("Checking approved variants")
		vl, err := getApprovedVariants(s.approvedVariantsUrl)
		if err != nil {
			log.Printf("Error downloading approved variants: %s", err.Error())
		}

		if len(vl) == 0 {
			log.Printf("No approved variants found")
			break
		}

		log.Printf("%d approved variants downloaded", len(vl))

		nv := make(map[string]interface{})
		for _, v := range vl {
			nv[v] = nil
		}

		s.approvedVariants.mu.Lock()
		s.approvedVariants.list = nv
		s.approvedVariants.mu.Unlock()

		// Create a sorted list of unique approved variants for logging.
		vs := make([]string, len(nv))
		i := 0
		for k := range nv {
			vs[i] = k
			i++
		}

		sort.Strings(vs)

		log.Printf("%d approved variants processed: %v", len(vs), vs)
		time.Sleep(s.approvedVariantsRefreshInterval)
	}
}

// getApprovedVariants returns a list of approved variants from the given url.
func getApprovedVariants(url *url.URL) ([]string, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	v, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// The response should be a JSON array of strings. If it isn't, return an error.
	vl := make([]string, 0)
	err = json.Unmarshal(v, &vl)
	if err != nil {
		return nil, err
	}

	return vl, nil
}
