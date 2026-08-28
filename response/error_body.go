package response

type ErrorBody struct {
	Message string `json:"message"`
	Details any    `json:"details"`
}
