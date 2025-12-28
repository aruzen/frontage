package card_action

import (
	"errors"
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

type PieceSummonActionState struct {
	holderID uuid.UUID
	cardID   uuid.UUID
	summonID uuid.UUID
	deckType model.DeckType
	point    pkg.Point

	holder *model.LocalPlayer
	piece  card.Piece
}

func (s PieceSummonActionState) HolderId() uuid.UUID { return s.holderID }
func (s PieceSummonActionState) Holder() *model.LocalPlayer {
	return s.holder
}
func (s PieceSummonActionState) CardID() uuid.UUID   { return s.cardID }
func (s PieceSummonActionState) SummonID() uuid.UUID { return s.summonID }
func (s PieceSummonActionState) DeckType() model.DeckType {
	return s.deckType
}
func (s PieceSummonActionState) Point() pkg.Point { return s.point }
func (s PieceSummonActionState) PieceCard() card.Piece {
	return s.piece
}

func (s PieceSummonActionState) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"holder_id": s.holderID.String(),
		"card_id":   s.cardID.String(),
		"summon_id": s.summonID.String(),
		"deck_type": s.deckType,
		"point":     pkg.PointToMap(s.point),
	}
}

func (s *PieceSummonActionState) FromMap(b *model.Board, m map[string]interface{}) error {
	var err error
	s.holderID, err = pkg.ToUUID(m["holder_id"])
	if err != nil {
		return err
	}
	s.cardID, err = pkg.ToUUID(m["card_id"])
	if err != nil {
		return err
	}
	s.summonID, err = pkg.ToUUID(m["summon_id"])
	if err != nil {
		return err
	}
	val, err := pkg.ToInt(m["deck_type"])
	if err != nil {
		return fmt.Errorf("deck_type: %w", err)
	}
	s.deckType = model.DeckType(val)
	s.point, err = pkg.PointFromMap(m["point"])
	if err != nil {
		return fmt.Errorf("point: %w", err)
	}

	if b == nil {
		return nil
	}
	player, ok := b.FindPlayer(s.holderID)
	if !ok {
		return errors.New("player not found")
	}
	s.holder, ok = player.(*model.LocalPlayer)
	if !ok {
		return errors.New("player does not support deck operations")
	}
	deck := s.holder.GetDeck(s.deckType)
	if deck == nil {
		return errors.New("deck not found")
	}
	c, ok := deck.FindById(s.cardID)
	if !ok {
		return fmt.Errorf("deck does not contain card: %v", s.cardID)
	}
	s.piece, ok = c.(card.Piece)
	if !ok {
		return fmt.Errorf("card is not piece: %v", s.cardID)
	}
	return nil
}

func NewPieceSummonActionState(holder *model.LocalPlayer, deckType model.DeckType, c card.Piece, point pkg.Point) PieceSummonActionState {
	return PieceSummonActionState{
		holderID: holder.ID(),
		holder:   holder,
		cardID:   c.Id(),
		piece:    c,
		deckType: deckType,
		point:    point,
	}
}

type PieceSummonActionContext struct {
	event.BaseEffectContext
	Point pkg.Point
}

func (c PieceSummonActionContext) ToMap() map[string]interface{} {
	result := c.BaseEffectContext.ToMap()
	result["point"] = pkg.PointToMap(c.Point)
	return result
}

func (c *PieceSummonActionContext) FromMap(m map[string]interface{}) error {
	if err := c.BaseEffectContext.FromMap(m); err != nil {
		return err
	}
	p, err := pkg.PointFromMap(m["point"])
	if err != nil {
		return fmt.Errorf("point: %w", err)
	}
	c.Point = p
	return nil
}

type PieceSummonAction struct {
	logic.BaseAction[PieceSummonActionState, PieceSummonActionContext]
}

func (PieceSummonAction) Tag() logic.EffectActionTag { return action.CARD_PIECE_SUMMON_ACTION }
func (a PieceSummonAction) LocalizeTag() pkg.LocalizeTag {
	return pkg.LocalizeTag(a.Tag())
}

func (e PieceSummonAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if s, ok := state.(PieceSummonActionState); ok {
		return &PieceSummonActionContext{BaseEffectContext: event.BaseEffectContext{}, Point: s.point},
			logic.Summary{"point": pkg.PointToMap(s.point)}
	}
	slog.Warn("State was not PieceSummonActionState.")
	return nil, nil
}

func (e PieceSummonAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) (*model.Board, logic.Summary) {
	s, c, ok := e.CastStateContext(state, context)
	if !ok || board == nil {
		return board, nil
	}
	next := board.Next()
	owner, ok := next.FindPlayer(s.holderID)
	if !ok {
		return board, nil
	}
	local, ok := owner.(*model.LocalPlayer)
	if !ok {
		return board, nil
	}
	deck := local.GetDeck(s.deckType)
	if deck == nil {
		return board, nil
	}
	if !deck.RemoveById(s.cardID) {
		return board, nil
	}
	if s.piece == nil {
		return board, nil
	}
	piece := model.NewBasePiece(s.summonID, owner, c.Point, s.piece.HP(), s.piece.MP(), s.piece.ATK(), s.piece.LegalMoves(), s.piece.AttackRanges())
	piece.SetActionsUsedThisTurn(piece.MaxActionsPerTurn())
	if !next.SetPiece(c.Point, piece) {
		return board, nil
	}
	return next, logic.Summary{"point": pkg.PointToMap(c.Point), "card_id": s.cardID.String()}
}
