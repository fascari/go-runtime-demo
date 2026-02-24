package gcbenchmark

type InputPayload struct {
	Allocations int    `json:"allocations"`
	SizeKB      int    `json:"size_kb"`
	Pattern     string `json:"pattern"` // "short-lived", "long-lived", "mixed"
}
