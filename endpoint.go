package main

type Endpoint struct {
	ResponseCode int               `json:"response_code"`
	Headers      map[string]string `json:"headers"`
	Body         interface{}       `json:"body"`
}

type EndpointIndex struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

type VariableEndpointIndex struct {
	EndpointIndex
	Variant string `json:"variant"`
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
