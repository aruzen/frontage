package logic_test

import (
	"testing"

	"frontage/pkg"
	"frontage/pkg/engine/impl/action/piece_action"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
)

func withoutPanic(t *testing.T, run func()) {
	t.Helper()
	defer func() {
		if value := recover(); value != nil {
			t.Errorf("unexpected panic: %v", value)
		}
	}()
	run()
}

func TestIntegrityCheckRejectsTypedNil(t *testing.T) {
	var action *piece_action.PieceAttackAction
	var state *piece_action.PieceAttackActionState
	for _, tc := range []struct {
		name   string
		action logic.Action
		state  logic.ActionState
	}{
		{"action", action, &piece_action.PieceAttackActionState{}},
		{"state", piece_action.PieceAttackAction{}, state},
	} {
		t.Run(tc.name, func(t *testing.T) {
			withoutPanic(t, func() {
				if logic.IntegrityCheck(tc.action, tc.state) {
					t.Error("typed nil passed the integrity check")
				}
			})
		})
	}
}

func TestCastStateContextRejectsTypedNil(t *testing.T) {
	base := logic.BaseAction[piece_action.PieceOperateActionState, piece_action.PieceHPContext]{}
	state := &piece_action.PieceOperateActionState{}
	context := &piece_action.PieceHPContext{}
	for _, tc := range []struct {
		name    string
		state   *piece_action.PieceOperateActionState
		context *piece_action.PieceHPContext
	}{
		{"state", nil, context},
		{"context", state, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, ok := base.CastStateContext(tc.state, tc.context); ok {
				t.Error("typed nil was reported as a usable state/context pair")
			}
		})
	}
	gotState, gotContext, ok := base.CastStateContext(state, context)
	if !ok || gotState != state || gotContext != context {
		t.Fatal("valid state/context was rejected or replaced")
	}
}

type validationListener struct {
	model.Player
	notifications *int
}

func (p *validationListener) Copy() model.Player {
	return &validationListener{Player: p.Player.Copy(), notifications: p.notifications}
}

func (p *validationListener) PreListen(*logic.EventSystem, logic.Action, logic.ActionState) {
	(*p.notifications)++
}

func newValidationSystem() (*logic.EventSystem, *int) {
	count := new(int)
	player := &validationListener{Player: model.NewProxyPlayer(uuid.New()), notifications: count}
	board := model.NewBoard(model.NewBoardInfo(pkg.Size{Width: 2, Height: 2}, model.GENERATION_STRATEGY_CHAIN),
		[2]model.Player{player, model.NewProxyPlayer(uuid.New())})
	return logic.NewEventSystem(board), count
}

func TestInvalidEventsAreRejectedBeforeListeners(t *testing.T) {
	var action *piece_action.PieceAttackAction
	var state *piece_action.PieceAttackActionState
	for _, tc := range []struct {
		name  string
		event *logic.EffectEvent
	}{
		{"nil_event", nil},
		{"nil_action", logic.NewEffectEvent(nil, &logic.SystemNoticeState{})},
		{"typed_nil_action", logic.NewEffectEvent(action, &piece_action.PieceAttackActionState{})},
		{"nil_state", logic.NewEffectEvent(piece_action.PieceAttackAction{}, nil)},
		{"typed_nil_state", logic.NewEffectEvent(piece_action.PieceAttackAction{}, state)},
		{"wrong_state", logic.NewEffectEvent(piece_action.PieceAttackAction{}, &logic.SystemNoticeState{})},
	} {
		for _, operation := range []string{"emit", "chain", "chain_line"} {
			t.Run(tc.name+"/"+operation, func(t *testing.T) {
				es, count := newValidationSystem()
				es.Emit(logic.NewEffectEvent(logic.GameStartAction{}, &logic.SystemNoticeState{}))
				withoutPanic(t, func() {
					switch operation {
					case "emit":
						es.Emit(tc.event)
					case "chain":
						es.Chain(tc.event)
					case "chain_line":
						es.ChainLine(tc.event)
					}
					es.Resolve()
				})
				if *count != 1 {
					t.Errorf("invalid event reached a listener: notifications = %d; want 1", *count)
				}
				wantEffects := 1
				if operation == "emit" {
					wantEffects = 0
				}
				if len(es.AppliedEffects) != wantEffects {
					t.Errorf("applied effects = %d; want %d", len(es.AppliedEffects), wantEffects)
				}
			})
		}
	}
}

func TestChainsBeforeEmitAreIgnored(t *testing.T) {
	for _, operation := range []string{"chain", "chain_line"} {
		t.Run(operation, func(t *testing.T) {
			es, count := newValidationSystem()
			event := logic.NewEffectEvent(logic.GameStartAction{}, &logic.SystemNoticeState{})
			withoutPanic(t, func() {
				if operation == "chain" {
					es.Chain(event)
				} else {
					es.ChainLine(event)
				}
			})
			if *count != 0 {
				t.Error("event without a parent reached a listener")
			}
		})
	}
}
