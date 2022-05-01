package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	svg "github.com/ajstarks/svgo"
)

// main starts the server.
func main() {
	listenAddress := ":8080"
	http.HandleFunc("/api/badge", badge)
	log.Printf("About to listen on %s", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

// badge returns the badge for the given status.
func badge(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "image/svg+xml")
	s := svg.New(w)

	approved := checkVariant(variant)
	log.Printf("Variant: %s approval is %t (%s)", variant, approved, r.Header["X-Request-Id"])
	makeBadge(variant, approved, s)
	log.Printf("End request from: %s (%s)", r.RemoteAddr, r.Header["X-Request-Id"])
}

// makeBadge creates a badge for the given status.
func makeBadge(text string, approved bool, s *svg.SVG) {
	var rightFill string
	switch approved {
	case true:
		rightFill = "#4c1"
	case false:
		rightFill = "#c41"
		text = "NOT APPROVED"
	}

	w := 200
	h := 20
	split := 0.35
	rightRectX := int(float64(w) * split)
	rightRectW := w - rightRectX
	pathStyle := fmt.Sprintf("M%d 0h4v20h-4z", rightRectX)
	leftTextX := rightRectX / 2
	rightTextX := ((w - rightRectX) / 2) + rightRectX
	sc := []svg.Offcolor{
		{Offset: 0, Color: "#bbb", Opacity: 0.1},
		{Offset: 100, Opacity: 0.1},
	}
	s.Start(w, h)
	s.LinearGradient("g", 0, 0, 0, uint8(w), sc)
	s.Rect(0, 0, w, h, "fill=\"#555\" rx=\"3\"")
	s.Rect(rightRectX, 0, rightRectW, h, fmt.Sprintf("fill=\"%s\" rx=\"3\"", rightFill)) // This is the green rectangle
	s.Path(pathStyle, fmt.Sprintf("fill=\"%s\"", rightFill))
	s.Rect(0, 0, w, h, "fill=\"url(#g)\" rx=\"3\"")
	s.Gstyle("fill: rgb(255, 255, 255);	text-anchor: middle; font-family: DejaVu Sans, Verdana, Geneva, sans-serif;	font-size: 11;")
	s.Text(leftTextX, 15, "alz-variant", "fill=\"#010101\"	fill-opacity=\"0.3\"")
	s.Text(leftTextX, 14, "alz-variant")
	s.Text(rightTextX, 15, text, "fill=\"#010101\"	fill-opacity=\"0.3\"")
	s.Text(rightTextX, 14, text)
	s.Gend()
	s.End()
}

// checkVariant returns the status of the given variant.
func checkVariant(variant string) bool {
	return variant == "canadapubsec"
}
