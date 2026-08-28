package response

type RespondableError interface {
	error
	StatusCode() int
	Message() string
	Details() any
}
