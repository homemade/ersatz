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

type EndpointCache map[EndpointIndex]*Endpoint

func NewEndpoint() *Endpoint {
	return &Endpoint{}
}
