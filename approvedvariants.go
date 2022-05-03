package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func updateApprovedVariants(av ApprovedVariants, url string) {
	for {
		log.Printf("Checking approved variants")
		vl, err := getApprovedVariants(&url)
		if err != nil {
			log.Printf("Error downloading approved variants: %s", err.Error())
		}

		for _, v := range vl {
			av[v] = nil
		}

		var vs string
		for v := range av {
			vs = vs + (fmt.Sprintf("%s, ", v))
		}

		log.Printf("Approved variants: %s", vs)
		time.Sleep(time.Minute * 15)
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
	err = json.Unmarshal(v, &vl)
	if err != nil {
		return nil, err
	}

	return vl, nil
}
