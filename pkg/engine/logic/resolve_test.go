package logic_test

import (
	"frontage/pkg"
	"frontage/pkg/engine/impl/action/card_action"
	"frontage/pkg/engine/impl/action/common"
	"frontage/pkg/engine/impl/action/piece_action"
	"frontage/pkg/engine/impl/action/register"
	"frontage/pkg/engine/impl/card"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func newResolveBoard() (*model.Board, *model.BasePiece) {
	players := [2]model.Player{model.NewLocalPlayer(uuid.New(), nil, nil), model.NewLocalPlayer(uuid.New(), nil, nil)}
	board := model.NewBoard(model.NewBoardInfo(pkg.Size{Width: 3, Height: 3}, model.GENERATION_STRATEGY_CHAIN), players)
	piece := model.NewBasePiece(uuid.New(), players[0], pkg.Point{}, 10, 0, 3, nil, nil)
	board.SetPiece(pkg.Point{}, piece)
	return board, piece
}

func TestResolveTestIsolatesStateAndPreservesPendingEvents(t *testing.T) {
	for _, scenario := range []string{"single", "branch", "cancelled"} {
		t.Run(scenario, func(t *testing.T) {
			board, original := newResolveBoard()
			state := &piece_action.PieceOperateActionState{}
			if err := state.FromMap(board, map[string]interface{}{"piece_id": original.Id().String(), "value": 3}); err != nil {
				t.Fatal(err)
			}
			es := logic.NewEventSystem(board)
			es.Emit(logic.NewEffectEvent(piece_action.PieceHPDecreaseAction{}, state))
			wantHP, wantEffects := 7, 1
			switch scenario {
			case "branch":
				es.Chain(logic.NewEffectEvent(piece_action.PieceHPDecreaseAction{}, state))
				wantHP, wantEffects = 4, 2
			case "cancelled":
				es.Chain(logic.NewModifyEvent(common.CancelModifyAction{}, &common.CancelState{Reason: "test"}))
				wantHP, wantEffects = 10, 0
			}
			for simulation := 0; simulation < 2; simulation++ {
				result, effects := es.ResolveTest()
				if result == nil {
					t.Fatal("simulation returned no board")
				}
				piece, _ := result.GetPiece(original.Id())
				if piece.HP() != wantHP || len(effects) != wantEffects {
					t.Fatalf("simulation: HP=%d, effects=%d; want %d, %d", piece.HP(), len(effects), wantHP, wantEffects)
				}
				if original.HP() != 10 || state.Piece() != original || len(es.AppliedEffects) != 0 || len(es.Summaries) != 0 {
					t.Fatal("simulation changed the original board, state or result buffers")
				}
			}
			es.Resolve()
			actual, _ := es.Board.GetPiece(original.Id())
			if actual.HP() != wantHP || len(es.AppliedEffects) != wantEffects || original.HP() != 10 {
				t.Fatal("real resolution after simulation did not preserve the expected result and previous board")
			}
		})
	}
}

func TestResolveTestClonesLinkedAttackStates(t *testing.T) {
	register.Init()
	board, attacker := newResolveBoard()
	target := model.NewBasePiece(uuid.New(), board.Players()[1], pkg.Point{X: 1}, 9, 0, 1, nil, nil)
	board.SetPiece(target.Position(), target)
	state := &piece_action.PieceAttackActionState{}
	if err := state.FromMap(board, map[string]interface{}{
		"piece_id": attacker.Id().String(), "point": pkg.PointToMap(target.Position()), "value": 3, "action_cost": 1,
	}); err != nil {
		t.Fatal(err)
	}
	es := logic.NewEventSystem(board)
	es.Emit(logic.NewEffectEvent(piece_action.PieceAttackAction{}, state))
	pendingDamage := state.DecreaseState().ToMap()
	for simulation := 0; simulation < 2; simulation++ {
		result, effects := es.ResolveTest()
		if result == nil {
			t.Fatal("simulation returned no board")
		}
		resultTarget, _ := result.GetPiece(target.Id())
		resultAttacker, _ := result.GetPiece(attacker.Id())
		if resultTarget.HP() != 6 || resultAttacker.UsedActionCost() != 1 || len(effects) != 2 {
			t.Fatal("simulated attack did not resolve linked damage and action cost")
		}
		if target.HP() != 9 || attacker.UsedActionCost() != 0 || !reflect.DeepEqual(pendingDamage, state.DecreaseState().ToMap()) {
			t.Fatal("simulation changed original pieces or the pending child state")
		}
	}
	es.Resolve()
	actual, _ := es.Board.GetPiece(target.Id())
	if actual.HP() != 6 || target.HP() != 9 || attacker.UsedActionCost() != 0 {
		t.Fatal("real attack after simulation changed the previous board or used stale damage")
	}
}

func TestCardEffectUpdatesCopiedHolder(t *testing.T) {
	board, _ := newResolveBoard()
	player := board.Players()[0].(*model.LocalPlayer)
	original := card.NewBasePiece(model.NewBaseCard("test", nil), 9, 0, 1, nil, nil)
	player.Hand = model.Cards{original}
	state := card_action.NewPieceActionState(player, model.DECK_TYPE_HAND, original, 3)
	es := logic.NewEventSystem(board)
	es.Emit(logic.NewEffectEvent(card_action.PieceHPDecreaseAction{}, &state))
	result, _ := es.ResolveTest()
	if result == nil {
		t.Fatal("simulation returned no board")
	}
	resultPlayer, _ := result.FindPlayer(player.ID())
	resultCard := resultPlayer.(*model.LocalPlayer).Hand[0].(card.Piece)
	if resultCard.HP() != 6 || original.HP() != 9 {
		t.Fatal("simulated card effect did not isolate the original holder")
	}
	es.Resolve()
	actualPlayer, _ := es.Board.FindPlayer(player.ID())
	actual := actualPlayer.(*model.LocalPlayer).Hand[0].(card.Piece)
	if actual.HP() != 6 || original.HP() != 9 {
		t.Fatal("real card effect mutated the old holder or missed the new holder")
	}
}
