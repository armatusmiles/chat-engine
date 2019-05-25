package clients

import "github.com/armatusmiles/gospeak-chat-engine/db/models"

type IChatClient interface {
	SendMessage(msg *models.Message)
}
