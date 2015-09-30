package multisafego

import (
	"net/url"
	"strconv"
)

// Gateway represents a gateway in multisafepay(IDEAL, Paypal)
type Gateway struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// Gateways will return the gateways for the api
// Use default values if no filters apply: Gateways("", "", 0)
// locale should be in ISO 639-1
// currency should be in ISO 4217
// amount should be in cents
func (m *MultiSafePay) Gateways(locale, currency string, amount int, country string) ([]Gateway, *APIError) {
	m.baseURL.Path = Path("/gateways")

	v := url.Values{}
	if locale != "" {
		v.Add("locale", locale)
	}

	if currency != "" {
		v.Add("currency", currency)
	}

	if amount != 0 {
		v.Add("amount", strconv.Itoa(amount))
	}

	if country != "" {
		v.Add("country", country)
	}

	m.baseURL.RawQuery = v.Encode()

	var gateways []Gateway
	err := m.Execute(m.baseURL, "GET", nil, &gateways)
	return gateways, err
}
