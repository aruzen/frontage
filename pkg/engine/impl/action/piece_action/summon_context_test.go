package piece_action

import (
	"encoding/json"
	"testing"

	"frontage/pkg"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
)

func roundTripSummonContext(t *testing.T, context *PieceSummonActionContext) *PieceSummonActionContext {
	t.Helper()
	data, err := json.Marshal(context.ToMap())
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	var restored PieceSummonActionContext
	if err := restored.FromMap(decoded); err != nil {
		t.Fatal(err)
	}
	return &restored
}

func TestSummonContextRoundTripRetainsPiece(t *testing.T) {
	for _, override := range []bool{false, true} {
		name := "state piece"
		if override {
			name = "different context piece"
		}
		t.Run(name, func(t *testing.T) {
			board, piece := newTestBoardWithPiece(t, pkg.Point{X: 1, Y: 1})
			target := pkg.Point{X: 3, Y: 3}
			state := &PieceSummonActionState{pieceID: piece.Id(), piece: piece, point: target}
			action := PieceSummonAction{}
			context, _ := action.Act(state, nil, nil)
			if override {
				piece = model.NewBasePiece(uuid.New(), board.Players()[1], pkg.Point{X: 2, Y: 2}, 10, 5, 3, nil, nil)
				board.SetPiece(piece.Position(), piece)
				context.(*PieceSummonActionContext).Piece = piece
			}
			restored := roundTripSummonContext(t, context.(*PieceSummonActionContext))
			if restored.ToMap()["piece_id"] != piece.Id().String() {
				t.Error("piece identity was lost during serialization")
			}
			next, _ := action.Solve(board, state, restored)
			if got := next.Entities()[target.X][target.Y]; got == nil || got.Id() != piece.Id() {
				t.Fatal("restored context did not place its piece")
			}
		})
	}
}

func TestSummonUnresolvedContextPreservesBoard(t *testing.T) {
	for _, id := range []string{"", uuid.New().String()} {
		board, occupant := newTestBoardWithPiece(t, pkg.Point{X: 1, Y: 1})
		input := map[string]interface{}{"point": pkg.PointToMap(occupant.Position())}
		if id != "" {
			input["piece_id"] = id
		}
		var context PieceSummonActionContext
		if err := context.FromMap(input); err != nil {
			t.Fatal(err)
		}
		state := &PieceSummonActionState{pieceID: occupant.Id(), piece: occupant}
		next, summary := (PieceSummonAction{}).Solve(board, state, &context)
		if next != board || summary != nil {
			t.Fatal("unresolved context changed the board")
		}
		if _, ok := next.GetPiece(occupant.Id()); !ok {
			t.Fatal("unresolved context removed an existing piece")
		}
	}
}

func TestSummonContextRejectsInvalidPieceID(t *testing.T) {
	_, piece := newTestBoardWithPiece(t, pkg.Point{X: 1, Y: 1})
	context := &PieceSummonActionContext{Piece: piece}
	input := map[string]interface{}{"point": pkg.PointToMap(piece.Position()), "piece_id": "invalid"}
	if err := context.FromMap(input); err == nil {
		t.Fatal("invalid piece ID was accepted")
	}
	if context.Piece != nil {
		t.Fatal("old piece reference survived restoration")
	}
}
