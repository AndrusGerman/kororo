package models

import "kororo/internal/core/domain/types"

type ActionField struct {
	Name string `bson:"name" json:"name"`
}

type Action struct {
	*Base              `bson:",inline" json:",inline"`
	Description        string                   `bson:"description" json:"description"`
	ActionProccessType types.ActionProccessType `bson:"action_proccess_type" json:"action_proccess_type"`

	Fields                 []*ActionField           `bson:"fields" json:"fields"`
	ProcessLLMSystemPrompt string                   `bson:"process_llm_system_prompt" json:"process_llm_system_prompt"`
	ResponseType           types.ActionResponseType `bson:"response_type" json:"response_type"`

	Command *ActionCommand `bson:"command" json:"command"`
	Http    *ActionHttp    `bson:"http" json:"http"`
}

type ActionCommand struct {
	Command string `bson:"command" json:"command"`
}

type ActionHttp struct {
	Url                    string `bson:"url" json:"url"`
	Method                 string `bson:"method" json:"method"`
	Body                   string `bson:"body" json:"body"`
	HttpValueNameResponse  string `bson:"http_value_name_response" json:"http_value_name_response"`
	CheckLLMResponsePrompt string `bson:"check_llm_response_prompt" json:"check_llm_response_prompt"`
}
