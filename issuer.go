package multisafego

// Issuer represents an issuer for an Gateway
type Issuer struct {
	Code        string `json:code`
	Description string `json:description`
}

// Issuers will return the issuers for the supplied gateway id
func (m *MultiSafePay) Issuers(id string) ([]Issuer, *APIError) {
	m.baseURL.Path = Path("/issuers/" + id)

	var issuers []Issuer
	err := m.Execute(m.baseURL, "GET", nil, &issuers)
	return issuers, err
}
