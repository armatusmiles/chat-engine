package clients

import "github.com/gospeak/protorepo/dbmanage"

type IChatClient interface {
	SendMessage(msg *dbmanage.ChatMessage)
}
