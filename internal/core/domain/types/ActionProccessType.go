package types

type ActionProccessType int

const (
	ActionProccessTypeRequestHttp ActionProccessType = 1
	ActionProccessTypeLLMResponse ActionProccessType = 2
	ActionProccessTypeCommand     ActionProccessType = 3
	ActionProccessTypeBasicFormat ActionProccessType = 4
)
