// Package tables implements (URL, HTTP method) -> SermoCRM action lookup
// tables. A table is an RPC service's set of endpoints.
package tables

// Table is a mapping of URLs to Endpoints.
type Table map[string][]Endpoint

// Endpoint is an RPC endpoint.
type Endpoint struct {
	// Method is the case-sensitive method.
	Method string

	// Unauthenticated is true if the endpoint does not require authentication.
	Unauthenticated bool

	// Action is the Endpoint's action.
	Action string
}

// Mapping maps (URL, HTTP method) -> SermoCRM actions.
type Mapping struct {
	t Table
}

// NewMapping creates a new Mapping
func NewMapping(fns ...func() Table) (*Mapping, error) {
	return nil, nil
}

// Mapping finds an endpoint based on a URL and HTTP method pair.
func (m *Mapping) Lookup(url, method string) (Endpoint, error) {
	const errNotFound = err("url not found")

	eps, ok := m.t[url]
	if !ok {
		return Endpoint{}, errNotFound
	}
	for _, e := range eps {
		if e.Method == method {
			return e, nil
		}
	}
	return Endpoint{}, errNotFound
}

type err string

func (e err) Error() string { return string(e) }
