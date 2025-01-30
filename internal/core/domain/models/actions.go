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
}

type ActionCommand struct {
	Command string `bson:"command" json:"command"`
}
