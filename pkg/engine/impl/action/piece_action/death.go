package piece_action

import (
	"fmt"
	"frontage/internal/event"
	"frontage/pkg"
	"frontage/pkg/engine/impl/action"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
)

type PieceDeathActionState struct {
	pieceID uuid.UUID
	point   pkg.Point
	piece   model.Piece
}

func (s PieceDeathActionState) PieceID() uuid.UUID { return s.pieceID }
func (s PieceDeathActionState) Point() pkg.Point   { return s.point }
func (s PieceDeathActionState) Piece() model.Piece { return s.piece }

func (s PieceDeathActionState) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"piece_id": s.pieceID.String(),
		"point":    pkg.PointToMap(s.point),
	}
}

func (s *PieceDeathActionState) FromMap(b *model.Board, m map[string]interface{}) error {
	id, err := pkg.ToUUID(m["piece_id"])
	if err != nil {
		return fmt.Errorf("piece_id: %w", err)
	}
	s.pieceID = id
	p, err := pkg.PointFromMap(m["point"])
	if err != nil {
		return fmt.Errorf("point: %w", err)
	}
	s.point = p
	if b != nil {
		if piece, ok := b.GetPiece(id); ok {
			s.piece = piece
		}
	}
	return nil
}

func NewPieceDeathActionState(piece model.Piece) PieceDeathActionState {
	return PieceDeathActionState{
		pieceID: piece.Id(),
		point:   piece.Position(),
		piece:   piece,
	}
}

type PieceDeathActionContext struct {
	event.BaseEffectContext
	Reason string
}

func (c PieceDeathActionContext) ToMap() map[string]interface{} {
	result := c.BaseEffectContext.ToMap()
	if c.Reason != "" {
		result["reason"] = c.Reason
	}
	return result
}

func (c *PieceDeathActionContext) FromMap(m map[string]interface{}) error {
	if err := c.BaseEffectContext.FromMap(m); err != nil {
		return err
	}
	if v, ok := m["reason"]; ok {
		if s, ok := v.(string); ok {
			c.Reason = s
		}
	}
	return nil
}

type PieceDeathAction struct {
	logic.BaseAction[PieceDeathActionState, PieceDeathActionContext]
}

func (PieceDeathAction) Tag() logic.EffectActionTag { return action.PIECE_DEATH_ACTION }
func (a PieceDeathAction) LocalizeTag() pkg.LocalizeTag {
	return pkg.LocalizeTag(a.Tag())
}

func (e PieceDeathAction) Act(state logic.ActionState, beforeAction logic.EffectAction, beforeContext logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if s, ok := state.(*PieceDeathActionState); ok {
		return &PieceDeathActionContext{BaseEffectContext: event.BaseEffectContext{}, Reason: ""},
			logic.Summary{"piece_id": s.pieceID.String()}
	}
	return nil, nil
}

func (e PieceDeathAction) Solve(board *model.Board, state logic.ActionState, context logic.EffectContext) (*model.Board, logic.Summary) {
	s, _, ok := e.CastStateContext(state, context)
	if !ok || board == nil {
		return board, nil
	}
	next := board.Next()
	if s.pieceID != uuid.Nil {
		if piece, ok := next.GetPiece(s.pieceID); ok {
			_ = next.RemovePiece(piece.Position())
		} else {
			_ = next.RemovePiece(s.point)
		}
		return next, logic.Summary{"piece_id": s.pieceID.String()}
	}
	_ = next.RemovePiece(s.point)
	return next, logic.Summary{"point": pkg.PointToMap(s.point)}
}

type PieceDeathByBurnAction struct{ PieceDeathAction }

func (PieceDeathByBurnAction) Tag() logic.EffectActionTag { return action.PIECE_DEATH_BY_BURN_ACTION }
func (a PieceDeathByBurnAction) LocalizeTag() pkg.LocalizeTag {
	return pkg.LocalizeTag(a.Tag())
}

func (e PieceDeathByBurnAction) Act(state logic.ActionState, beforeAction logic.EffectAction, beforeContext logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if s, ok := state.(*PieceDeathActionState); ok {
		return &PieceDeathActionContext{BaseEffectContext: event.BaseEffectContext{}, Reason: "burn"},
			logic.Summary{"piece_id": s.pieceID.String(), "reason": "burn"}
	}
	return nil, nil
}

type PieceDeathByAttackAction struct{ PieceDeathAction }

func (PieceDeathByAttackAction) Tag() logic.EffectActionTag {
	return action.PIECE_DEATH_BY_ATTACK_ACTION
}
func (a PieceDeathByAttackAction) LocalizeTag() pkg.LocalizeTag {
	return pkg.LocalizeTag(a.Tag())
}

func (e PieceDeathByAttackAction) Act(state logic.ActionState, beforeAction logic.EffectAction, beforeContext logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if s, ok := state.(*PieceDeathActionState); ok {
		return &PieceDeathActionContext{BaseEffectContext: event.BaseEffectContext{}, Reason: "attack"},
			logic.Summary{"piece_id": s.pieceID.String(), "reason": "attack"}
	}
	return nil, nil
}

type PieceDeathByInvasionAction struct{ PieceDeathAction }

func (PieceDeathByInvasionAction) Tag() logic.EffectActionTag {
	return action.PIECE_DEATH_BY_INVASION_ACTION
}
func (a PieceDeathByInvasionAction) LocalizeTag() pkg.LocalizeTag {
	return pkg.LocalizeTag(a.Tag())
}

func (e PieceDeathByInvasionAction) Act(state logic.ActionState, beforeAction logic.EffectAction, beforeContext logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if s, ok := state.(*PieceDeathActionState); ok {
		return &PieceDeathActionContext{BaseEffectContext: event.BaseEffectContext{}, Reason: "invasion"},
			logic.Summary{"piece_id": s.pieceID.String(), "reason": "invasion"}
	}
	return nil, nil
}
