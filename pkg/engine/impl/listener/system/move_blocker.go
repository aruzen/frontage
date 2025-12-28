package system

import (
	"frontage/pkg/engine/impl/action"
	"frontage/pkg/engine/impl/action/common"
	"frontage/pkg/engine/impl/action/piece_action"
	"frontage/pkg/engine/logic"
)

// MoveBlocker cancels move events when the target cell is occupied.
type MoveBlocker struct{}

func (MoveBlocker) PreListen(es *logic.EventSystem, act logic.Action, state logic.ActionState) {
	effect, ok := act.(logic.EffectAction)
	if !ok || effect.Tag() != action.ENTITY_MOVE_ACTION {
		return
	}
	var moveState *piece_action.PieceMoveActionState
	switch s := state.(type) {
	case *piece_action.PieceMoveActionState:
		moveState = s
	}
	if moveState == nil || es == nil || es.Board == nil {
		return
	}
	to := moveState.To()
	size := es.Board.Size()
	if to.X < 0 || to.X >= size.Width || to.Y < 0 || to.Y >= size.Height {
		return
	}
	if es.Board.Entities()[to.X][to.Y] == nil && es.Board.Structures()[to.X][to.Y] == nil {
		return
	}
	cancel := common.CancelModifyAction{}
	es.Chain(logic.NewModifyEvent(cancel, &common.CancelState{Reason: "destination occupied"}))
}
