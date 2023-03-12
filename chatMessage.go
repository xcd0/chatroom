package chat

import (
	"time"

	"github.com/google/uuid"
)

// ChatMessageはチャットメッセージを表す
type ChatMessage struct {
	PostID    string    `json:"post_id"`   // 投稿UUID
	RoomID    string    `json:"room_id"`   // チャットルームUUID
	UserID    string    `json:"user_id"`   // 投稿者のユーザーUUID
	Content   string    `json:"content"`   // 投稿内容
	Timestamp time.Time `json:"timestamp"` // 投稿日時
}

func newChatMessage(roomID string, userUUID string, content string) *ChatMessage {
	id := uuid.NewString()
	chat.MessageList[id] = &ChatMessage{
		PostID:    id,
		RoomID:    roomID,
		UserID:    userUUID,
		Content:   content,
		Timestamp: time.Now(),
	}
	return chat.MessageList[id]
}
