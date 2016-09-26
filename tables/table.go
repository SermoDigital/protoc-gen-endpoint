// Package tables implements (URL, HTTP method) -> SermoCRM Action lookup
// tables. A table is an RPC service's set of endpoints.
package tables

import "strings"

// Table maps URLs to Endpoints.
type Table map[string]Endpoint

// Action is an RPC Action.
type Action struct {
	// Name is the qualified name of the Action. E.g., users.Create.
	Name string

	Method string

	// Unauthenticated is true if the endpoint does not require authentication.
	Unauthenticated bool
}

// Endpoint is an HTTP endpoint.
type Endpoint struct {
	// Methods is all comma-delimited list of all the HTTP methods this endpoint
	// supports.
	Methods string

	// Actions are all the Actions that correspond to the given Endpoint.
	Actions []Action
}

// Add adds the Action to the Endpoint and adjusts the Action's S and E fields.
func (e *Endpoint) Add(act Action) {
	if e.Methods == "" {
		e.Methods = act.Method
	} else if !strings.Contains(e.Methods, act.Method) {
		e.Methods += "," + act.Method
	}
	e.Actions = append(e.Actions, act)
}

var emptyAct Action

func (e Endpoint) Find(method string) (Action, bool) {
	for _, act := range e.Actions {
		if method == act.Method {
			return act, true
		}
	}
	return emptyAct, false
}

// Mapping maps (URL, HTTP method) -> SermoCRM Actions.
type Mapping struct {
	t Table
}

// MakeMapping creates a new Mapping.
func MakeMapping(fns ...func() Table) Mapping {
	t := make(Table)
	for _, fn := range fns {
		for url, eps := range fn() {
			ep := t[url]
			for _, act := range eps.Actions {
				ep.Add(act)
			}
			t[url] = ep
		}
	}
	return Mapping{t: t}
}

// Mapping finds an endpoint based on a URL and HTTP method pair.
func (m Mapping) Get(url string) (Endpoint, bool) {
	ep, ok := m.t[url]
	return ep, ok
}
