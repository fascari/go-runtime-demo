package stresstest

type InputPayload struct {
	Allocations int `json:"allocations"`
	Goroutines  int `json:"goroutines"`
}
