package response

type AddLink struct {
	Alias string `json:"alias,omitempty"`
	Error string `json:"error,omitempty"`
}
