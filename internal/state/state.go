package state

import (
	"math"
	"sync"
	"time"

	"github.com/usernamenenad/ibft-go/internal/network"
)

type ConsensusInstance struct {
	InstanceId    uint64
	CurrentRound  uint64
	PreparedRound uint64
	PreparedValue []byte
	InputValue    []byte
	DecisionRound uint64
	DecidedValue  []byte
	CommitQuorum  []*network.IbftMessage
}

type ProcessState struct {
	ProcessId          uint64
	TotalProcesses     uint64
	MaxFaultyProcesses uint64
	QuorumSize         uint64
	Instances          map[uint64]*ConsensusInstance
	Timer              *RoundChangeTimer
	MessageLog         []*network.IbftMessage
	ValidityPredicate  func([]byte) bool
	Mu                 sync.RWMutex
	Network            *network.Network
}

func NewProcess(processId uint64, totalProcesses uint64, validityPredicate func(value []byte) bool) *ProcessState {
	maxFaultyProcesses := uint64(math.Floor(float64(totalProcesses-1) / 3))
	quorumSize := (totalProcesses+maxFaultyProcesses)/2 + 1

	return &ProcessState{
		ProcessId:          processId,
		TotalProcesses:     totalProcesses,
		MaxFaultyProcesses: maxFaultyProcesses,
		QuorumSize:         quorumSize,
		Instances:          make(map[uint64]*ConsensusInstance),
		Timer:              NewRoundChangeTimer(),
		MessageLog:         make([]*network.IbftMessage, 0),
		ValidityPredicate:  validityPredicate,
		Mu:                 sync.RWMutex{},
	}
}

func NewConsensusInstance(instanceId uint64, inputValue []byte) *ConsensusInstance {
	return &ConsensusInstance{
		InstanceId:    instanceId,
		CurrentRound:  1,
		DecisionRound: 0xFFFFFFFFFFFFFFFF,
		PreparedRound: 0,
		PreparedValue: nil,
		DecidedValue:  nil,
		InputValue:    inputValue,
		CommitQuorum:  make([]*network.IbftMessage, 0),
	}
}

func (ps *ProcessState) Broadcast(msg *network.IbftMessage) error {
	return (*ps.Network).Broadcast(msg)
}

func (ps *ProcessState) startRoundChangeTimer(instanceId uint64, timeout time.Duration) {
	ps.Timer.Start(timeout, func() {
		ps.HandleRoundTimerExpiration(instanceId)
	})
}

func (ps *ProcessState) calculateTimeout(round uint64) time.Duration {
	return ps.Timer.calculateTimeout(round)
}

func (ps *ProcessState) getLeader(instanceId uint64, round uint64) uint64 {
	return (instanceId + round) % ps.TotalProcesses
}

func (ps *ProcessState) isLeader(instanceId uint64, round uint64) bool {
	return ps.ProcessId == ps.getLeader(instanceId, round)
}
