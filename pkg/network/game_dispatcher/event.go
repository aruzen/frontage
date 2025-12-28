package game_dispatcher

import (
	"frontage/pkg/engine/logic"
	"frontage/pkg/network"
	"frontage/pkg/network/data"
	"frontage/pkg/network/game_api"
	"frontage/pkg/network/translator"
)

// ActEventDispatcher builds packets for action events.
type ActEventDispatcher struct {
	ActionResultTrans  *translator.ActionResultTranslator
	ActionSummaryTrans *translator.ActionSummaryTranslator
}

func NewActEventDispatcher(ActionResultTrans *translator.ActionResultTranslator, ActionSummaryTrans *translator.ActionSummaryTranslator) *ActEventDispatcher {
	return &ActEventDispatcher{
		ActionResultTrans:  ActionResultTrans,
		ActionSummaryTrans: ActionSummaryTrans,
	}
}

// DispatchPacket serializes the given payloads into a JSON packet.
func (d *ActEventDispatcher) DispatchPacket(AppliedEffects []logic.ActionResult, Summaries []*[]logic.ActionSummary) (network.Packet, error) {
	packet := game_api.ActEventPacket{
		Events:    make([]data.ActionResult, len(AppliedEffects)),
		Summaries: make([][]data.ActionSummary, len(Summaries)),
	}
	return packet, nil
}
