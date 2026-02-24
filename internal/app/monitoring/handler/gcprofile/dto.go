package gcprofile

type InputPayload struct {
	DurationSeconds int    `json:"duration_seconds"`
	ProfileType     string `json:"profile_type"` // "heap", "cpu", "goroutine", "allocs"
}
