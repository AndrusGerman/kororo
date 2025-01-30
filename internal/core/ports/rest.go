package ports

type RestAdapter interface {
	Post(url string, body any, out any) error
	Stream(url string, body any) (<-chan StreamRest, error)
}

type StreamRest interface {
	Parse(body any) error
}
