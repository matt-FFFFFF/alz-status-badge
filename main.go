package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type server struct {
	approvedVariantsUrl             *url.URL
	approvedVariants                *ApprovedVariants
	approvedVariantsRefreshInterval time.Duration
	listenAddress                   string
	router                          *http.ServeMux
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

type ApprovedVariants struct {
	list map[string]interface{}
	mu   sync.RWMutex
}

const DefaultApprovedVariantsUrl = "https://raw.githubusercontent.com/matt-FFFFFF/alz-status-badge/main/data/approved-variants.json"

// main is used to to call run(), which returns an error.
func main() {
	if err := run(); err != nil {
		log.Fatalf("%s\n", err)
	}
}

// run is used to start the server.
func run() error {
	// Set up the application configuration defaults and read environment variables
	u, err := url.Parse(DefaultApprovedVariantsUrl)
	if err != nil {
		return err
	}

	av := ApprovedVariants{
		list: make(map[string]interface{}),
	}

	s := &server{
		approvedVariantsUrl:             u,
		approvedVariants:                &av,
		approvedVariantsRefreshInterval: time.Minute * 15,
		listenAddress:                   ":8080",
		router:                          http.NewServeMux(),
	}

	if e, ok := os.LookupEnv("ALZSTATUSBADGE_LISTEN_ADDRESS"); ok {
		s.listenAddress = e
	}

	if e, ok := os.LookupEnv("ALZSTATUSBADGE_APPROVED_VARIANTS_URL"); ok {
		u, err := url.Parse(e)
		if err != nil {
			return err
		}
		s.approvedVariantsUrl = u
	}

	if e, ok := os.LookupEnv("ALZSTATUSBADGE_APPROVED_VARIANTS_REFRESH_INTERVAL"); ok {
		d, err := time.ParseDuration(e)
		if err != nil {
			return err
		}
		s.approvedVariantsRefreshInterval = d
	}

	log.Printf("About to listen on %s", s.listenAddress)
	go updateApprovedVariants(s)
	s.routes()
	return http.ListenAndServe(s.listenAddress, s)
}
