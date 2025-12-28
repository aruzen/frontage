package piece_action

import (
	"testing"

	"frontage/pkg"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
)

func newTestBoardWithPiece(t *testing.T, pos pkg.Point) (*model.Board, *model.BasePiece) {
	t.Helper()
	owner := model.NewProxyPlayer(uuid.New())
	players := [2]model.Player{owner, model.NewProxyPlayer(uuid.New())}
	board := model.NewBoard(model.NewBoardInfo(pkg.Size{Width: 5, Height: 5}, model.GENERATION_STRATEGY_SWAP), players)
	piece := model.NewBasePiece(uuid.New(), owner, pos, 10, 5, 3, nil, nil)
	if !board.SetPiece(pos, piece) {
		t.Fatalf("failed to set piece at %v", pos)
	}
	return board, piece
}

func TestPieceActionStatesToMapFromMap(t *testing.T) {
	// Summon
	{
		board, piece := newTestBoardWithPiece(t, pkg.Point{X: 1, Y: 2})
		state := PieceSummonActionState{pieceID: piece.Id(), point: piece.Position(), piece: piece}
		m := state.ToMap()
		var dst PieceSummonActionState
		if err := dst.FromMap(board, m); err != nil {
			t.Fatalf("summon FromMap error: %v", err)
		}
		if dst.pieceID != piece.Id() || dst.point != piece.Position() {
			t.Fatalf("summon state mismatch: %+v", dst)
		}
		if dst.piece == nil || dst.piece.Id() != piece.Id() {
			t.Fatalf("summon piece not restored")
		}
	}

	// Move
	{
		from := pkg.Point{X: 0, Y: 0}
		to := pkg.Point{X: 2, Y: 3}
		board, piece := newTestBoardWithPiece(t, from)
		state := PieceMoveActionState{pieceID: piece.Id(), from: from, to: to, actionCost: 2, piece: piece}
		m := state.ToMap()
		var dst PieceMoveActionState
		if err := dst.FromMap(board, m); err != nil {
			t.Fatalf("move FromMap error: %v", err)
		}
		if dst.pieceID != piece.Id() || dst.from != from || dst.to != to || dst.actionCost != 2 {
			t.Fatalf("move state mismatch: %+v", dst)
		}
		if dst.piece == nil || dst.piece.Id() != piece.Id() {
			t.Fatalf("move piece not restored")
		}
	}

	// Attack
	{
		pos := pkg.Point{X: 1, Y: 1}
		target := pkg.Point{X: 3, Y: 4}
		board, piece := newTestBoardWithPiece(t, pos)
		state := PieceAttackActionState{pieceID: piece.Id(), point: target, value: 7, actionCost: 1, piece: piece}
		m := state.ToMap()
		var dst PieceAttackActionState
		if err := dst.FromMap(board, m); err != nil {
			t.Fatalf("attack FromMap error: %v", err)
		}
		if dst.pieceID != piece.Id() || dst.point != target || dst.value != 7 || dst.actionCost != 1 {
			t.Fatalf("attack state mismatch: %+v", dst)
		}
		if dst.piece == nil || dst.piece.Id() != piece.Id() {
			t.Fatalf("attack piece not restored")
		}
	}

	// Operate
	{
		pos := pkg.Point{X: 2, Y: 2}
		board, piece := newTestBoardWithPiece(t, pos)
		state := PieceOperateActionState{pieceID: piece.Id(), piece: piece, value: 4}
		m := state.ToMap()
		var dst PieceOperateActionState
		if err := dst.FromMap(board, m); err != nil {
			t.Fatalf("operate FromMap error: %v", err)
		}
		if dst.pieceID != piece.Id() || dst.value != 4 {
			t.Fatalf("operate state mismatch: %+v", dst)
		}
		if dst.piece == nil || dst.piece.Id() != piece.Id() {
			t.Fatalf("operate piece not restored")
		}
	}

	// Death
	{
		board, piece := newTestBoardWithPiece(t, pkg.Point{X: 4, Y: 4})
		state := NewPieceDeathActionState(piece)
		m := state.ToMap()
		var dst PieceDeathActionState
		if err := dst.FromMap(board, m); err != nil {
			t.Fatalf("death FromMap error: %v", err)
		}
		if dst.pieceID != piece.Id() || dst.point != piece.Position() {
			t.Fatalf("death state mismatch: %+v", dst)
		}
		if dst.piece == nil || dst.piece.Id() != piece.Id() {
			t.Fatalf("death piece not restored")
		}
	}
}
