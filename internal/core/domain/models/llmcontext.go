package models

type RoleID string

const (
	UserRoleID      RoleID = "user"
	AssistantRoleID RoleID = "assistant"
	SystemRoleID    RoleID = "system"
)

type Message struct {
	RoleID  RoleID
	Content string
}

func NewMessage(content string, roleId RoleID) *Message {
	return &Message{
		RoleID:  roleId,
		Content: content,
	}
}

func NewMessageFromUser(content string) *Message {
	return NewMessage(content, UserRoleID)
}

func NewMessageFromAssistant(content string) *Message {
	return NewMessage(content, AssistantRoleID)
}

func NewMessageFromSystem(content string) *Message {
	return NewMessage(content, SystemRoleID)
}
