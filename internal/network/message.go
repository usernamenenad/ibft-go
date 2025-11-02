package network

import "time"

type IbftMessageType int

const (
	MessageTypePrePrepare IbftMessageType = iota
	MessageTypePrepare
	MessageTypeCommit
	MessageTypeRoundChange
)

type IbftMessage struct {
	SenderProcessId   uint64
	ReceiverProcessId uint64
	Type              IbftMessageType
	InstanceId        uint64
	Round             uint64
	Value             []byte
	PreparedRound     uint64
	PreparedValue     []byte
	Justification     []*IbftMessage
}

type Message struct {
	SenderProcessId   uint64
	ReceiverProcessId uint64
	Signature         []byte
	Timestamp         time.Time
	Payload           *IbftMessage
}
