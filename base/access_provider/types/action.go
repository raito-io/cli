package types

//go:generate go run github.com/raito-io/enumer -type=Action -json -yaml -transform=lower
type Action int

const (
	Promise Action = iota // Deprecated promises are set on who item
	Grant
	Deny
	Mask
	Filtered
	Purpose // Deprecated purposes are now moved to be a grant category
	Share
)
