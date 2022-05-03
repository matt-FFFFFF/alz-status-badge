package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sort"
	"time"
)

// updateApprovedVariants updates the approved variants map every 15 minutes. Use a goroutine to call this function.
func updateApprovedVariants(av ApprovedVariants, url string) {
	for {
		log.Printf("Checking approved variants")
		vl, err := getApprovedVariants(&url)
		if err != nil {
			log.Printf("Error downloading approved variants: %s", err.Error())
		}

		log.Printf("Approved variants downloaded: %d", len(vl))

		for _, v := range vl {
			av[v] = nil
		}

		vs := make([]string, len(av))
		i := 0
		for k := range av {
			vs[i] = k
			i++
		}

		sort.Strings(vs)

		log.Printf("%d approved variants: %v", len(vs), vs)
		time.Sleep(time.Minute * 15)
	}
}

// getApprovedVariants returns a list of approved variants from the given url.
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
	err = json.Unmarshal(v, &vl)
	if err != nil {
		return nil, err
	}

	return vl, nil
}
