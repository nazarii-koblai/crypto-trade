package model

// ErrorResponse represents error response sturct.
type ErrorResponse struct {
	Message string
}

// Error returns error message.
func (er ErrorResponse) Error() string {
	return er.Message
}
