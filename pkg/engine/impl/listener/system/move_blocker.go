package system

import (
	"frontage/pkg/engine/impl/action"
	"frontage/pkg/engine/impl/action/piece_action"
	"frontage/pkg/engine/logic"
)

// MoveBlocker cancels move events when the target cell is occupied.
type MoveBlocker struct{}

func (MoveBlocker) PreListen(es *logic.EventSystem, act logic.Action, state interface{}) {
	effect, ok := act.(logic.EffectAction)
	if !ok || effect.Tag() != action.ENTITY_MOVE_ACTION {
		return
	}
	moveState, ok := state.(piece_action.PieceMoveActionState)
	if !ok || es == nil || es.Board == nil {
		return
	}
	to := moveState.To()
	size := es.Board.Size()
	if to.X < 0 || to.X >= size.Width || to.Y < 0 || to.Y >= size.Height {
		return
	}
	if es.Board.Entities()[to.X][to.Y] == nil {
		return
	}
	cancel := piece_action.MoveCancelModifyAction{}
	es.Chain(logic.NewModifyEvent(cancel, piece_action.MoveCancelState{Reason: "destination occupied"}))
}
