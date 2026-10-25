package player_action

import (
	"frontage/internal/event"
	"frontage/pkg"
	"frontage/pkg/engine/impl/action"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
	"log/slog"
)

type DrawActionState struct {
	playerID uuid.UUID
	player   model.Player
}

func (s DrawActionState) PlayerID() uuid.UUID { return s.playerID }
func (s DrawActionState) Player() model.Player {
	return s.player
}

func (s DrawActionState) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"player_id": s.playerID.String(),
	}
}

func (s *DrawActionState) FromMap(b *model.Board, m map[string]interface{}) error {
	s.player = nil
	id, err := pkg.ToUUID(m["player_id"])
	if err != nil {
		return err
	}
	s.playerID = id

	if b == nil {
		return nil
	}
	player, ok := b.FindPlayer(id)
	if !ok {
		slog.Warn("player not found", "player", id)
		return nil
	}
	s.player = player
	return nil
}

func NewDrawActionState(player model.Player) *DrawActionState {
	return &DrawActionState{
		playerID: player.ID(),
		player:   player,
	}
}

type DrawActionContext struct {
	event.BaseEffectContext
}

type DrawMainDeckAction struct {
	logic.BaseAction[DrawActionState, DrawActionContext]
}

type DrawSubDeckAction struct {
	logic.BaseAction[DrawActionState, DrawActionContext]
}

func (DrawMainDeckAction) Tag() logic.EffectActionTag { return action.PLAYER_DRAW_MAIN_DECK_ACTION }
func (DrawSubDeckAction) Tag() logic.EffectActionTag  { return action.PLAYER_DRAW_SUB_DECK_ACTION }
func (a DrawMainDeckAction) LocalizeTag() pkg.LocalizeTag {
	return pkg.LocalizeTag(a.Tag())
}
func (a DrawSubDeckAction) LocalizeTag() pkg.LocalizeTag { return pkg.LocalizeTag(a.Tag()) }

func (a DrawMainDeckAction) Act(state logic.ActionState, beforeAction logic.EffectAction, beforeContext logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if s, ok := state.(*DrawActionState); ok {
		return &DrawActionContext{BaseEffectContext: event.NewBaseEffectContext()},
			logic.Summary{"player_id": s.playerID.String(), "deck_type": int(model.DECK_TYPE_MAIN)}
	}
	return nil, nil
}

func (a DrawSubDeckAction) Act(state logic.ActionState, beforeAction logic.EffectAction, beforeContext logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if s, ok := state.(*DrawActionState); ok {
		return &DrawActionContext{BaseEffectContext: event.NewBaseEffectContext()},
			logic.Summary{"player_id": s.playerID.String(), "deck_type": int(model.DECK_TYPE_SUB)}
	}
	return nil, nil
}

func (a DrawMainDeckAction) Solve(board *model.Board, state logic.ActionState, context logic.EffectContext) (*model.Board, logic.Summary) {
	return solveDraw(board, state, model.DECK_TYPE_MAIN)
}

func (a DrawSubDeckAction) Solve(board *model.Board, state logic.ActionState, context logic.EffectContext) (*model.Board, logic.Summary) {
	return solveDraw(board, state, model.DECK_TYPE_SUB)
}

func solveDraw(board *model.Board, state logic.ActionState, deckType model.DeckType) (*model.Board, logic.Summary) {
	s, ok := state.(*DrawActionState)
	if !ok {
		return board, nil
	}
	next := board.Next()
	player, ok := next.FindPlayer(s.playerID)
	if !ok {
		slog.Warn("player not found", "player", s.playerID)
		return board, nil
	}
	local, ok := player.(*model.LocalPlayer)
	if !ok {
		slog.Warn("player is not LocalPlayer", "player", s.playerID)
		return board, nil
	}
	deck := local.GetDeck(deckType)
	if deck == nil {
		slog.Warn("deck not found", "player", s.playerID, "deckType", deckType)
		return board, nil
	}
	if len(*deck) == 0 {
		return board, nil
	}
	hand := local.GetDeck(model.DECK_TYPE_HAND)
	if hand == nil {
		slog.Warn("hand deck not found", "player", s.playerID)
		return board, nil
	}
	card := deck.Draw()
	hand.PushTop(model.Cards{card})
	return next, logic.Summary{
		"player_id": s.playerID.String(),
		"card_id":   card.Id().String(),
		"deck_type": int(deckType),
	}
}
