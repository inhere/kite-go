package httptpl

// Response struct
type Response struct {
	Status      int    `json:"status"`
	Struct      any    `json:"struct"`
	ContentType string `json:"content_type"`
	// Required fields in response.
	Required []string `json:"required"`
	// Asserts {name: "$.name == inhere"}
	Asserts map[string]any `json:"asserts"`
}
