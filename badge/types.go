package badge

type BadgeRequest struct {
	Endpoint string
	Label    string
	Message  string
	Color    string
	Style    BadgeStyle
}

type BadgeStyle int64

// BadgeStyle is the style of the badge. Used for shields.io call.
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
