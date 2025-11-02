package state

import (
	"fmt"

	"github.com/usernamenenad/ibft-go/internal/network"
)

func (ps *ProcessState) ValidatePrePrepareMessage(instance *ConsensusInstance, msg *network.IbftMessage) error {
	if msg.Round != instance.CurrentRound {
		return fmt.Errorf("message round %d does not match current instance round %d", msg.Round, instance.CurrentRound)
	}

	if ps.getLeader(instance.InstanceId, instance.CurrentRound) != msg.SenderProcessId {
		return fmt.Errorf("message of type PRE-PREPARE is not from leader for round %d", instance.CurrentRound)
	}

	if !ps.ValidityPredicate(msg.Value) {
		return fmt.Errorf("message does not satisfy external validity predicate")
	}

	if !ps.JustifyPrePrepare(msg) {
		return fmt.Errorf("message is not PRE-PREPARE justified")
	}

	return nil
}

func (ps *ProcessState) JustifyPrePrepare(msg *network.IbftMessage) bool {
	if msg.Round == 1 {
		return true
	}

	return true
}

func (ps *ProcessState) JustifyRoundChange() {}
