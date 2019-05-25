package websocketext

type Code int

const (
	RoomEnter Code = iota + 4000
	RoomExit
	RoomCreate
)
