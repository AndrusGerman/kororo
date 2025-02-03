package ports

type RestAdapter interface {
	Post(url string, headers map[string]string, body any, out any) error
	Stream(url string, headers map[string]string, body any) (<-chan StreamRest, error)
}

type StreamRest interface {
	Parse(body any) error
}
