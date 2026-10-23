package piece_action

import (
	"frontage/pkg"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
	"testing"
)

func resolveDeath(board *model.Board, state *PieceDeathActionState) *model.Board {
	es := logic.NewEventSystem(board)
	es.Emit(logic.NewEffectEvent(PieceDeathAction{}, state))
	es.Resolve()
	return es.Board
}

func TestDeathMissingTargetPreservesReplacement(t *testing.T) {
	position := pkg.Point{X: 1, Y: 1}
	board, original := newTestBoardWithPiece(t, position)
	state := NewPieceDeathActionState(original)
	board.RemovePiece(position)
	replacement := model.NewBasePiece(uuid.New(), board.Players()[1], position, 10, 0, 1, nil, nil)
	board.SetPiece(position, replacement)
	result := resolveDeath(board, &state)
	if piece, ok := result.GetPiece(replacement.Id()); !ok || piece.HP() != 10 {
		t.Fatal("death of an absent piece removed its replacement")
	}
}

func TestDeathFollowsMovedTarget(t *testing.T) {
	from := pkg.Point{X: 1, Y: 1}
	to := pkg.Point{X: 2, Y: 1}
	board, target := newTestBoardWithPiece(t, from)
	state := NewPieceDeathActionState(target)
	board.SetPiece(to, target)
	replacement := model.NewBasePiece(uuid.New(), board.Players()[1], from, 10, 0, 1, nil, nil)
	board.SetPiece(from, replacement)
	result := resolveDeath(board, &state)
	if _, ok := result.GetPiece(target.Id()); ok {
		t.Fatal("moved target survived its death event")
	}
	if _, ok := result.GetPiece(replacement.Id()); !ok {
		t.Fatal("replacement at the previous position was removed")
	}
}

func TestDeathByPositionWithoutID(t *testing.T) {
	position := pkg.Point{X: 1, Y: 1}
	board, target := newTestBoardWithPiece(t, position)
	state := &PieceDeathActionState{}
	if err := state.FromMap(board, map[string]interface{}{"piece_id": uuid.Nil.String(), "point": pkg.PointToMap(position)}); err != nil {
		t.Fatal(err)
	}
	result := resolveDeath(board, state)
	if _, ok := result.GetPiece(target.Id()); ok {
		t.Fatal("position-only removal did not remove its target")
	}
}
