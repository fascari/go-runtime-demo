package mineparallel

type InputPayload struct {
	Data       string `json:"data"`
	Goroutines int    `json:"goroutines"`
}
