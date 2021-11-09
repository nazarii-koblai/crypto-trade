package model

const (
	FBTC CurrencyType = iota
	FETH
)

// CurrencyType represents different currecy types.
type CurrencyType int

// String returns string representation of currency type.
func (ct CurrencyType) String() string {
	switch ct {
	case FBTC:
		return "FBTH"
	case FETH:
		return "FETH"
	default:
		return ""
	}
}
