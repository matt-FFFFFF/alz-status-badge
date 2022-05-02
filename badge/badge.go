package badge

import (
	"fmt"
	"io"
	"net/http"
)

type BadgeStyle int64

// BadgeStyle is the style of the badge. Used for shieldsio call
const (
	Plastic BadgeStyle = iota
	Flat
	FlatSquare
	ForTheBadge
	Social
)

// Stringer interface for BadgeStyle
func (s BadgeStyle) String() string {
	switch s {
	case Plastic:
		return "plastic"
	case Flat:
		return "flat"
	case FlatSquare:
		return "flat-square"
	case ForTheBadge:
		return "for-the-badge"
	case Social:
		return "social"
	}
	return ""
}

// MadeShieldsioBadge generates the badge for the given variant and approval status.
func MakeShieldsioBadge(variant string, approved bool) ([]byte, error) {
	var color string
	switch approved {
	case true:
		color = "success"
	case false:
		color = "critical"
		variant = "NOT VALID"
	}

	badgeUri := fmt.Sprintf("https://img.shields.io/badge/alz--variant-%s-%s.svg?style=%s", variant, color, ForTheBadge)
	resp, err := http.Get(badgeUri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	badge, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return badge, nil
}
