package topics

// Topics is a structure represing kafka topic definition
type Topics struct {
	Topic     string   `json:"topic"`
	Room      string   `json:"room"`
	Event     string   `json:"event"`
	Enhancers []string `json:"enhancers"`
}
