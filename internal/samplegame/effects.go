package samplegame

import (
	"fmt"

	"frontage/internal/event"
	"frontage/pkg"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
)

// Action budget changes are a sample effect; the existing engine actions handle
// summoning, movement, attacks, HP/ATK changes, drawing and death.
type actionCostState struct {
	PieceID uuid.UUID
	Value   int
}

func (s *actionCostState) ToMap() map[string]interface{} {
	return map[string]interface{}{"piece_id": s.PieceID.String(), "value": s.Value}
}
func (s *actionCostState) FromMap(_ *model.Board, input map[string]interface{}) error {
	var err error
	s.PieceID, err = pkg.ToUUID(input["piece_id"])
	if err != nil {
		return fmt.Errorf("piece_id: %w", err)
	}
	s.Value, err = pkg.ToInt(input["value"])
	if err != nil {
		return fmt.Errorf("value: %w", err)
	}
	if s.Value < 0 {
		return fmt.Errorf("行動消費量は0以上です")
	}
	return nil
}

type actionCostContext struct {
	event.BaseEffectContext
	Value int
}

func (c *actionCostContext) ToMap() map[string]interface{} {
	result := c.BaseEffectContext.ToMap()
	result["value"] = c.Value
	return result
}
func (c *actionCostContext) FromMap(input map[string]interface{}) error {
	if err := c.BaseEffectContext.FromMap(input); err != nil {
		return err
	}
	value, err := pkg.ToInt(input["value"])
	if err != nil {
		return err
	}
	c.Value = value
	return nil
}

type actionCostEffect struct {
	logic.BaseAction[actionCostState, actionCostContext]
}

func (actionCostEffect) Tag() logic.EffectActionTag     { return "sample/action_cost" }
func (a actionCostEffect) LocalizeTag() pkg.LocalizeTag { return pkg.LocalizeTag(a.Tag()) }
func (actionCostEffect) Act(state logic.ActionState, _ logic.EffectAction, _ logic.EffectContext) (logic.EffectContext, logic.Summary) {
	s, ok := state.(*actionCostState)
	if !ok {
		return nil, nil
	}
	return &actionCostContext{Value: s.Value}, nil
}
func (a actionCostEffect) Solve(board *model.Board, state logic.ActionState, context logic.EffectContext) (*model.Board, logic.Summary) {
	s, c, ok := a.CastStateContext(state, context)
	if !ok || board == nil {
		return nil, nil
	}
	if _, exists := board.GetPiece(s.PieceID); !exists {
		return nil, nil
	}
	next := board.Next()
	piece, _ := next.GetPiece(s.PieceID)
	mutable := piece.(model.MutablePiece)
	mutable.SetUsedActionCost(c.Value)
	next.UpdatePiece(mutable)
	return next, logic.Summary{"piece_id": s.PieceID.String(), "used_action_cost": c.Value}
}
