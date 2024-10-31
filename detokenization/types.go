package detokenization

// RequestPayload represents the payload structure for the API request.
type RequestPayload struct {
	DataElement string   `json:"data_element"`
	Data        []string `json:"data"`
}

// Response represents the structure of the API response.
type Response struct {
	Encoding string   `json:"encoding"`
	Results  []string `json:"results"`
	Success  string   `json:"success"`
}
