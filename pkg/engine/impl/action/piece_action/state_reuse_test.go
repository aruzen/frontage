package piece_action

import (
	"testing"

	"frontage/pkg"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
)

func TestPieceStateReuseClearsReferences(t *testing.T) {
	point := pkg.Point{X: 1, Y: 1}
	board, piece := newTestBoardWithPiece(t, point)
	summon := &PieceSummonActionState{pieceID: piece.Id(), point: point}
	move := &PieceMoveActionState{pieceID: piece.Id(), from: point, to: point}
	attack := &PieceAttackActionState{pieceID: piece.Id(), point: point, value: 1}
	operate := &PieceOperateActionState{pieceID: piece.Id(), value: 1}
	death := NewPieceDeathActionState(piece)
	for _, tc := range []struct {
		name   string
		state  logic.ActionState
		cached func() bool
	}{
		{"summon", summon, func() bool { return summon.piece != nil }},
		{"move", move, func() bool { return move.piece != nil }},
		{"attack", attack, func() bool { return attack.piece != nil }},
		{"operate", operate, func() bool { return operate.piece != nil }},
		{"death", &death, func() bool { return death.piece != nil }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			valid := tc.state.ToMap()
			for _, scenario := range []string{"nil board", "missing piece", "invalid input"} {
				t.Run(scenario, func(t *testing.T) {
					if err := tc.state.FromMap(board, valid); err != nil || !tc.cached() {
						t.Fatalf("initial restore failed: %v", err)
					}
					input := tc.state.ToMap()
					var nextBoard *model.Board = board
					switch scenario {
					case "nil board":
						nextBoard = nil
					case "missing piece":
						input["piece_id"] = uuid.New().String()
					case "invalid input":
						input["piece_id"] = "invalid"
					}
					err := tc.state.FromMap(nextBoard, input)
					if (err != nil) != (scenario == "invalid input") {
						t.Errorf("unexpected restore error: %v", err)
					}
					if tc.cached() {
						t.Fatal("reference from the previous restore survived")
					}
				})
			}
		})
	}
}
