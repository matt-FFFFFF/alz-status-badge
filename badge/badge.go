package badge

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

const shieldsIoDefaultEndpoint = "https://img.shields.io/static/v1"

func NewBadgeRequest() BadgeRequest {
	return BadgeRequest{
		Endpoint: shieldsIoDefaultEndpoint,
	}
}

// Get generates the badge and returns it as a byte slice.
func Get(br BadgeRequest) ([]byte, error) {
	u, err := url.Parse(br.Endpoint)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("label", br.Label)
	q.Add("message", br.Message)
	q.Add("color", br.Color)
	q.Add("style", br.Style.String())
	u.RawQuery = q.Encode()

	log.Printf("Badge request: %s", u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}
