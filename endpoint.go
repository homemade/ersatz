package main

// FIXME! Rename to EnpointFile
type Endpoint struct {
	ResponseCode int               `json:"response_code"`
	Headers      map[string]string `json:"headers"`
	Body         interface{}       `json:"body"`
}

type EndpointIndex struct {
	URL    string
	Method string
}

type VariableEndpointIndex struct {
	EndpointIndex
	Variant string
}

type EndpointVariation struct {
	Variation string
	Count     int
}

type EndpointCache map[VariableEndpointIndex]*Endpoint
type EndpointVariationSchedule map[EndpointIndex]EndpointVariation

func NewEndpoint() *Endpoint {
	return &Endpoint{
		Headers: make(map[string]string),
	}
}
