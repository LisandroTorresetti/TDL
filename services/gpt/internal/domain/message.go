package domain

type MessageRole string

const (
	SystemRole MessageRole = "system"
	UserRole   MessageRole = "user"
)

type Message struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

func NewMessage(role MessageRole, content string) Message {
	return Message{
		Role:    role,
		Content: content,
	}
}

// getSystemMessage returns the prompt message used to respond the user's messages
func getSystemMessage() Message {
	return Message{
		Role: SystemRole,
		Content: "You are an assistant used in a Telegram bot that can summarize news. " +
			"Be brief, I dont want a response with much more than 40 words. " +
			"Your response should be in the same language that the provided news are.",
	}
}
