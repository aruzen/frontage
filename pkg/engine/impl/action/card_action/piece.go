package card_action

import (
	"fmt"
	"frontage/internal/event"
	"frontage/pkg"
	"frontage/pkg/engine/impl/action"
	"frontage/pkg/engine/impl/card"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
	"log/slog"
)

type PieceActionState struct {
	holder   uuid.UUID
	deckType model.DeckType
	card     uuid.UUID
	value    int
}

type PieceHPContext struct {
	*event.BaseEffectContext
	Value int
}

type PieceMPContext struct {
	*event.BaseEffectContext
	Value int
}

type PieceATKContext struct {
	*event.BaseEffectContext
	Value int
}

func (c PieceHPContext) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"value": c.Value,
	}
	if c.BaseEffectContext != nil {
		for k, v := range c.BaseEffectContext.ToMap() {
			result[k] = v
		}
	}
	return result
}

func (c *PieceHPContext) FromMap(m map[string]interface{}) error {
	if c.BaseEffectContext == nil {
		base := event.NewBaseEffectContext()
		c.BaseEffectContext = &base
	}
	if err := c.BaseEffectContext.FromMap(m); err != nil {
		return err
	}
	if v, ok := m["value"]; ok {
		num, err := pkg.ToInt(v)
		if err != nil {
			return fmt.Errorf("value: %w", err)
		}
		c.Value = num
	}
	return nil
}

func (c PieceMPContext) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"value": c.Value,
	}
	if c.BaseEffectContext != nil {
		for k, v := range c.BaseEffectContext.ToMap() {
			result[k] = v
		}
	}
	return result
}

func (c *PieceMPContext) FromMap(m map[string]interface{}) error {
	if c.BaseEffectContext == nil {
		base := event.NewBaseEffectContext()
		c.BaseEffectContext = &base
	}
	if err := c.BaseEffectContext.FromMap(m); err != nil {
		return err
	}
	if v, ok := m["value"]; ok {
		num, err := pkg.ToInt(v)
		if err != nil {
			return fmt.Errorf("value: %w", err)
		}
		c.Value = num
	}
	return nil
}

func (c PieceATKContext) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"value": c.Value,
	}
	if c.BaseEffectContext != nil {
		for k, v := range c.BaseEffectContext.ToMap() {
			result[k] = v
		}
	}
	return result
}

func (c *PieceATKContext) FromMap(m map[string]interface{}) error {
	if c.BaseEffectContext == nil {
		base := event.NewBaseEffectContext()
		c.BaseEffectContext = &base
	}
	if err := c.BaseEffectContext.FromMap(m); err != nil {
		return err
	}
	if v, ok := m["value"]; ok {
		num, err := pkg.ToInt(v)
		if err != nil {
			return fmt.Errorf("value: %w", err)
		}
		c.Value = num
	}
	return nil
}

type baseHPAction struct {
	logic.BaseAction[PieceActionState, PieceHPContext]
}

type baseMPAction struct {
	logic.BaseAction[PieceActionState, PieceMPContext]
}

type baseATKAction struct {
	logic.BaseAction[PieceActionState, PieceATKContext]
}

type PieceHPIncreaseAction struct{ baseHPAction }

type PieceHPDecreaseAction struct{ baseHPAction }

type PieceHPFixAction struct{ baseHPAction }

type PieceMPIncreaseAction struct{ baseMPAction }

type PieceMPDecreaseAction struct{ baseMPAction }

type PieceMPFixAction struct{ baseMPAction }

type PieceATKIncreaseAction struct{ baseATKAction }

type PieceATKDecreaseAction struct{ baseATKAction }

type PieceATKFixAction struct{ baseATKAction }

func (PieceHPIncreaseAction) Tag() logic.EffectActionTag { return action.CARD_PIECE_HP_INCREASE_ACTION }
func (PieceHPDecreaseAction) Tag() logic.EffectActionTag { return action.CARD_PIECE_HP_DECREASE_ACTION }
func (PieceHPFixAction) Tag() logic.EffectActionTag      { return action.CARD_PIECE_HP_FIX_ACTION }
func (PieceMPIncreaseAction) Tag() logic.EffectActionTag { return action.CARD_PIECE_MP_INCREASE_ACTION }
func (PieceMPDecreaseAction) Tag() logic.EffectActionTag { return action.CARD_PIECE_MP_DECREASE_ACTION }
func (PieceMPFixAction) Tag() logic.EffectActionTag      { return action.CARD_PIECE_MP_FIX_ACTION }
func (PieceATKIncreaseAction) Tag() logic.EffectActionTag {
	return action.CARD_PIECE_ATK_INCREASE_ACTION
}
func (PieceATKDecreaseAction) Tag() logic.EffectActionTag {
	return action.CARD_PIECE_ATK_DECREASE_ACTION
}
func (PieceATKFixAction) Tag() logic.EffectActionTag          { return action.CARD_PIECE_ATK_FIX_ACTION }
func (a PieceHPIncreaseAction) LocalizeTag() pkg.LocalizeTag  { return pkg.LocalizeTag(a.Tag()) }
func (a PieceHPDecreaseAction) LocalizeTag() pkg.LocalizeTag  { return pkg.LocalizeTag(a.Tag()) }
func (a PieceHPFixAction) LocalizeTag() pkg.LocalizeTag       { return pkg.LocalizeTag(a.Tag()) }
func (a PieceMPIncreaseAction) LocalizeTag() pkg.LocalizeTag  { return pkg.LocalizeTag(a.Tag()) }
func (a PieceMPDecreaseAction) LocalizeTag() pkg.LocalizeTag  { return pkg.LocalizeTag(a.Tag()) }
func (a PieceMPFixAction) LocalizeTag() pkg.LocalizeTag       { return pkg.LocalizeTag(a.Tag()) }
func (a PieceATKIncreaseAction) LocalizeTag() pkg.LocalizeTag { return pkg.LocalizeTag(a.Tag()) }
func (a PieceATKDecreaseAction) LocalizeTag() pkg.LocalizeTag { return pkg.LocalizeTag(a.Tag()) }
func (a PieceATKFixAction) LocalizeTag() pkg.LocalizeTag      { return pkg.LocalizeTag(a.Tag()) }

func (p PieceActionState) HolderId() uuid.UUID {
	return p.holder
}

func (p PieceActionState) Holder(b *model.Board) *model.LocalPlayer {
	player, ok := b.FindPlayer(p.holder)
	if !ok {
		slog.Warn("player not found.")
		return nil
	}

	deckPlayer, ok := player.(*model.LocalPlayer)
	if !ok {
		slog.Warn("player does not support deck operations.")
		return nil
	}
	return deckPlayer
}

func (p PieceActionState) CardID() uuid.UUID {
	return p.card
}

func (p PieceActionState) DeckAndCard(b *model.Board) (*model.Cards, card.Piece) {
	player := p.Holder(b)
	if player == nil {
		slog.Warn("player not found.")
		return nil, nil
	}

	deck := player.GetDeck(p.deckType)
	if deck == nil {
		slog.Warn("deck is nil for holder", "holder", p.holder, "deckType", p.deckType)
		return nil, nil
	}

	c, ok := deck.FindById(p.card)
	if !ok {
		slog.Warn("deck does not contain card", "card", p.card)
		return nil, nil
	}

	piece, ok := c.(card.Piece)
	if !ok {
		slog.Warn("card does not piece", "card", p.card)
		return nil, nil
	}

	return deck, piece
}

func (p PieceActionState) Value() int {
	return p.value
}

func updatePieceCard(board *model.Board, state *PieceActionState, mutate func(card.MutablePiece)) *model.Board {
	if board == nil || state == nil || mutate == nil {
		return board
	}

	next := board.Next()
	player := state.Holder(next)
	if player == nil {
		slog.Warn("player not found")
		return board
	}
	deck, card := state.DeckAndCard(next)
	if deck == nil || card == nil {
		slog.Warn("deck or card is nil")
		return board
	}
	mut := card.Copy()

	mutate(mut)

	if !deck.Update(mut) {
		slog.Warn("failed to update card in deck.", "deckType", state.deckType, "cardID", mut.Id())
	}

	return next
}

func NewPieceActionState(holder model.Player, deckType model.DeckType, card card.MutablePiece, value int) PieceActionState {
	return PieceActionState{
		deckType: deckType,
		holder:   holder.ID(),
		card:     card.Id(),
		value:    value,
	}
}

func (p baseHPAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if PieceState, ok := state.(PieceActionState); ok {
		base := event.NewBaseEffectContext()
		return &PieceHPContext{BaseEffectContext: &base, Value: PieceState.value}, logic.Summary{"value": PieceState.value}
	}
	slog.Warn("State was not PieceActionState.")
	return nil, nil
}

func (p PieceHPIncreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, hpContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceHPIncreaseAction.")
		return board, nil
	}

	board = updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetHP(m.HP() + hpContext.Value)
	})
	return board, logic.Summary{"delta_hp": hpContext.Value}
}

func (p PieceHPDecreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, hpContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceHPDecreaseAction.")
		return board, nil
	}

	board = updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetHP(m.HP() - hpContext.Value)
	})
	return board, logic.Summary{"delta_hp": -hpContext.Value}
}

func (p PieceHPFixAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, hpContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceHPFixAction.")
		return board, nil
	}

	board = updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetHP(hpContext.Value)
	})
	return board, logic.Summary{"set_hp": hpContext.Value}
}

func (p baseMPAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if PieceState, ok := state.(PieceActionState); ok {
		base := event.NewBaseEffectContext()
		return &PieceMPContext{BaseEffectContext: &base, Value: PieceState.value}, logic.Summary{"value": PieceState.value}
	}
	slog.Warn("State was not PieceActionState.")
	return nil, nil
}

func (p PieceMPIncreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, mpContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceMPIncreaseAction.")
		return board, nil
	}

	board = updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetMP(m.MP() + mpContext.Value)
	})
	return board, logic.Summary{"delta_mp": mpContext.Value}
}

func (p PieceMPDecreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, mpContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceMPDecreaseAction.")
		return board, nil
	}

	board = updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetMP(m.MP() - mpContext.Value)
	})
	return board, logic.Summary{"delta_mp": -mpContext.Value}
}

func (p PieceMPFixAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, mpContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceMPFixAction.")
		return board, nil
	}

	board = updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetMP(mpContext.Value)
	})
	return board, logic.Summary{"set_mp": mpContext.Value}
}

func (p baseATKAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if PieceState, ok := state.(PieceActionState); ok {
		base := event.NewBaseEffectContext()
		return &PieceATKContext{BaseEffectContext: &base, Value: PieceState.value}, logic.Summary{"value": PieceState.value}
	}
	slog.Warn("State was not PieceActionState.")
	return nil, nil
}

func (p PieceATKIncreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, atkContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceATKIncreaseAction.")
		return board, nil
	}

	board = updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetATK(m.ATK() + atkContext.Value)
	})
	return board, logic.Summary{"delta_atk": atkContext.Value}
}

func (p PieceATKDecreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, atkContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceATKDecreaseAction.")
		return board, nil
	}

	board = updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetATK(m.ATK() - atkContext.Value)
	})
	return board, logic.Summary{"delta_atk": -atkContext.Value}
}

func (p PieceATKFixAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, atkContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceATKFixAction.")
		return board, nil
	}

	board = updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetATK(atkContext.Value)
	})
	return board, logic.Summary{"set_atk": atkContext.Value}
}
