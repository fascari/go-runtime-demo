package stresstest

type InputPayload struct {
	Allocations int    `json:"allocations"`
	Goroutines  int    `json:"goroutines"`
	Pattern     string `json:"pattern"` // "short-lived", "long-lived", "mixed" (default: "short-lived")
}
