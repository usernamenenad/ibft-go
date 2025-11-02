package state

import (
	"fmt"

	"github.com/usernamenenad/ibft-go/internal/network"
)

func (ps *ProcessState) HandleStartMessage(instanceId uint64, inputValue []byte) error {
	ps.Mu.Lock()
	defer ps.Mu.Unlock()

	consensusInstance := NewConsensusInstance(instanceId, inputValue)
	ps.Instances[instanceId] = consensusInstance

	if ps.isLeader(instanceId, 1) {
		msg := &network.IbftMessage{
			SenderProcessId: ps.ProcessId,
			Type:            network.MessageTypePrePrepare,
			InstanceId:      instanceId,
			Round:           1,
			Value:           inputValue,
		}

		if err := ps.Broadcast(msg); err != nil {
			return fmt.Errorf("failed to broadcast PRE-PREPARE: %w", err)
		}
	}

	timeout := ps.calculateTimeout(1)
	ps.startRoundChangeTimer(instanceId, timeout)

	return nil
}

func (ps *ProcessState) HandlePrePrepareMessage(msg *network.IbftMessage) error {
	ps.Mu.Lock()
	defer ps.Mu.Unlock()

	instance, ok := ps.Instances[msg.InstanceId]
	if !ok {
		instance := NewConsensusInstance(msg.InstanceId, msg.Value)
		ps.Instances[msg.InstanceId] = instance
	}

	if err := ps.ValidatePrePrepareMessage(instance, msg); err != nil {
		return err
	}

	prepareMsg := &network.IbftMessage{
		SenderProcessId: ps.ProcessId,
		Type:            network.MessageTypePrepare,
		InstanceId:      msg.InstanceId,
		Round:           msg.Round,
		Value:           msg.Value,
	}

	if err := ps.Broadcast(prepareMsg); err != nil {
		return fmt.Errorf("failed to broadcast PRE-PREPARE: %w", err)
	}

	timeout := ps.calculateTimeout(instance.CurrentRound)
	ps.startRoundChangeTimer(msg.InstanceId, timeout)
	return nil
}

func (ps *ProcessState) HandleRoundTimerExpiration(instanceId uint64) {}
