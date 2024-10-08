package clir

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrorRouteNotFound = Error("route not found")
)
