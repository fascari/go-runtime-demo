package gcfinalizers

type InputPayload struct {
	Count     int  `json:"count"`
	TriggerGC bool `json:"trigger_gc"`
}
