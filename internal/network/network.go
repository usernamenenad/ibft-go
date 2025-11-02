package network

type Network interface {
	Broadcast(*IbftMessage) error
}
