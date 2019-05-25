package rooms

type FreeChatRoom struct {
	broadcaster *RoomBroadcaster
	clients     *[]ChatClient
}
