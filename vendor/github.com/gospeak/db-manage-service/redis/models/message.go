package models

type Message struct {
	RoomID    uint64 `json:"room_id"`
	UserID    uint64 `json:"user_id"`
	SendTime  int    `json:"send_time"` // timestamp
	Message   string `json:"message"`
	MessageID int    `json:"message_id"`
}
