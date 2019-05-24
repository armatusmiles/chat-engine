package main

type IChatRoom interface {
	CountClients() uint32
	AddClient(*ChatClient) bool
	RemoveClientById(ID uint64) bool
}
