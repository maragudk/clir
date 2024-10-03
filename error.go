package clir

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrorNotFound = Error("not found")
)
