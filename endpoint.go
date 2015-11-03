package main

type Endpoint struct {
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}

type EndpointIndex struct {
	URL     string
	Method  string
	Variant string
}

type EndpointVariation struct {
	Variation string
	Count     int
}

type EndpointCache map[EndpointIndex]*Endpoint
type EndpointVariationSchedule map[EndpointIndex]EndpointVariation

func NewEndpoint() *Endpoint {
	return &Endpoint{}
}
