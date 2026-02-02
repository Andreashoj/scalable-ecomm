package errors

type ErrorMessage string

const (
	NotFound   ErrorMessage = "Couldn't find what you were looking for."
	BadRequest ErrorMessage = "Something went wrong on the server"
)
