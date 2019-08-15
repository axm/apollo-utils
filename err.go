package apollo

import "strings"

type Error struct {
	message string
	err error
}

func NewError(message string, err error) *Error {
	return &Error{
		message: message,
		err:     err,
	}
}

func (error *Error) Error() string {
	builder := strings.Builder{}
	builder.WriteString(error.message)
	if error.err != nil {
		builder.WriteString("\n")
		builder.WriteString("Inner error:")
		builder.WriteString("\n")
		builder.WriteString(error.err.Error())
	}

	return builder.String()
}
