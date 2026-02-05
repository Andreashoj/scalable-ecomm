package errors

type ErrorMessage string

const (
	NotFound     ErrorMessage = "Couldn't find what you were looking for."
	Unauthorized ErrorMessage = "Unauthorized request"
	BadRequest   ErrorMessage = "Something went wrong on the server"
)
