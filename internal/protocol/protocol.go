package protocol

type Message struct {
	Type string `json:"type"`
	ID   string `json:"id,omitempty"`
}
