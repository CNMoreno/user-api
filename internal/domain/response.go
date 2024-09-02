package domain

// APIResponse response endpoints.
type APIResponse struct {
	Success bool    `json:"success"`
	Errors  *Errors `json:"errors,omitempty"`
	ID      string  `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	Email   string  `json:"email,omitempty"`
}

// Errors handles errors in endpoints.
type Errors struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Details string `json:"details,omitempty"`
}
