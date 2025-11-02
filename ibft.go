package ibft

import "github.com/usernamenenad/ibft-go/internal/state"

type Node struct {
	Addr         string
	ProcessState *state.ProcessState
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) Propose(value []byte) error {
	return nil
}
