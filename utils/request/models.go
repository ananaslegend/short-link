package request

type AddLink struct {
	Link  string `json:"link"`
	Alias string `json:"alias,omitempty"`
}
