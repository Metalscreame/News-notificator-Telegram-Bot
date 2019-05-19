package main

var chatMap = make(map[int64]struct{})

type Chat struct {}

func NewChat() Chat {
	return Chat{}
}

// AddChatToPull adds a chat to norify to chat pull
func (u *Chat) AddChatToPull(chatID int64) {
	chatMap[chatID] = struct{}{}
}

// GetChatIDs returns chat map
func (u *Chat) GetChatIDs() map[int64]struct{} {
	return chatMap
}