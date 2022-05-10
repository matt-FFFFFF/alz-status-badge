package main

func (s *server) routes() {
	s.router.HandleFunc("/api/badge/", s.handleBadgeGet())
}
