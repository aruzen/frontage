package piece_action_test

import (
	"testing"

	"frontage/pkg"
	"frontage/pkg/engine/impl/action/piece_action"
	"frontage/pkg/engine/impl/action/register"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
)

func TestRepeatedAttackDoesNotDamageTargetThatMovedAway(t *testing.T) {
	register.Init()
	for _, effect := range []logic.EffectAction{piece_action.PieceAttackAction{}, piece_action.PieceInvasionAction{}} {
		t.Run(string(effect.Tag()), func(t *testing.T) {
			testAttackAfterTargetMoved(t, effect)
		})
	}
}

func testAttackAfterTargetMoved(t *testing.T, repeated logic.EffectAction) {
	t.Helper()
	players := [2]model.Player{model.NewProxyPlayer(uuid.New()), model.NewProxyPlayer(uuid.New())}
	board := model.NewBoard(model.NewBoardInfo(pkg.Size{Width: 3, Height: 3}, model.GENERATION_STRATEGY_CHAIN), players)
	attacker := model.NewBasePiece(uuid.New(), players[0], pkg.Point{}, 10, 0, 3, nil, nil)
	target := model.NewBasePiece(uuid.New(), players[1], pkg.Point{X: 1}, 10, 0, 1, nil, nil)
	board.SetPiece(attacker.Position(), attacker)
	board.SetPiece(target.Position(), target)
	state := &piece_action.PieceAttackActionState{}
	if err := state.FromMap(board, map[string]interface{}{
		"piece_id": attacker.Id().String(), "point": pkg.PointToMap(target.Position()), "value": 3, "action_cost": 1,
	}); err != nil {
		t.Fatal(err)
	}
	es := logic.NewEventSystem(board)
	es.Emit(logic.NewEffectEvent(piece_action.PieceAttackAction{}, state))
	es.Resolve()
	victim, _ := es.Board.GetPiece(target.Id())
	if victim.HP() != 7 {
		t.Fatalf("first attack HP = %d; want 7", victim.HP())
	}

	es.Board = es.Board.Next()
	victim, _ = es.Board.GetPiece(target.Id())
	es.Board.SetPiece(pkg.Point{X: 2, Y: 2}, victim)
	es.Clean()
	es.Emit(logic.NewEffectEvent(repeated, state))
	es.Resolve()
	victim, _ = es.Board.GetPiece(target.Id())
	if victim.HP() != 7 {
		t.Errorf("attack on an empty cell damaged a previously targeted piece: HP = %d; want 7", victim.HP())
	}
	actor, _ := es.Board.GetPiece(attacker.Id())
	if actor.UsedActionCost() != 1 {
		t.Errorf("attack on an empty cell consumed action cost: %d; want 1", actor.UsedActionCost())
	}
}

func TestInvasionWithMissingPieceDoesNotPanic(t *testing.T) {
	register.Init()
	players := [2]model.Player{model.NewProxyPlayer(uuid.New()), model.NewProxyPlayer(uuid.New())}
	board := model.NewBoard(model.NewBoardInfo(pkg.Size{Width: 3, Height: 3}, model.GENERATION_STRATEGY_CHAIN), players)
	target := model.NewBasePiece(uuid.New(), players[1], pkg.Point{X: 1}, 10, 0, 1, nil, nil)
	board.SetPiece(target.Position(), target)
	state := &piece_action.PieceAttackActionState{}
	if err := state.FromMap(board, map[string]interface{}{
		"piece_id": uuid.NewString(), "point": pkg.PointToMap(target.Position()), "value": 3, "action_cost": 1,
	}); err != nil {
		t.Fatal(err)
	}
	es := logic.NewEventSystem(board)
	es.Emit(logic.NewEffectEvent(piece_action.PieceInvasionAction{}, state))
	es.Resolve()
	victim, exists := es.Board.GetPiece(target.Id())
	if !exists || victim.HP() != 10 || victim.Position() != target.Position() {
		t.Fatal("invasion by a missing piece changed the target")
	}
}
