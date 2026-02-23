package errors

// TODO: fix shared namespace
import "errors"

type ErrorMessage string

const (
	NotFound     ErrorMessage = "Couldn't find what you were looking for."
	Unauthorized ErrorMessage = "Unauthorized request"
	BadRequest   ErrorMessage = "Something went wrong on the server"
)

func New(msg string) error {
	return errors.New(msg)
}
