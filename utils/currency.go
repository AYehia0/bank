// supported currencies
package utils

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	EGP = "EGP"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, EGP, CAD:
		return true
	}
	return false
}
