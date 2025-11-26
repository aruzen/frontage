package model

import (
	"frontage/pkg"
	"testing"

	"github.com/google/uuid"
)

type dummyStructure struct {
	id uuid.UUID
}

func newDummyStructure() dummyStructure {
	return dummyStructure{id: uuid.New()}
}

func newDummyStructureWithID(id uuid.UUID) dummyStructure {
	return dummyStructure{id: id}
}

func (d dummyStructure) ID() uuid.UUID {
	return d.id
}

func newTestBoard() *Board {
	info := NewBoardInfo(pkg.Size{Width: 3, Height: 3}, GENERATION_STRATEGY_CHAIN)
	return NewBoard(info, [2]Player{NewProxyPlayer(uuid.New()), NewProxyPlayer(uuid.New())})
}

func TestBoardSetAndGetStructure(t *testing.T) {
	board := newTestBoard()
	pos := pkg.Point{X: 1, Y: 2}
	structure := newDummyStructure()

	if !board.SetStructure(pos, structure) {
		t.Fatalf("SetStructure(%v) returned false", pos)
	}

	got, ok := board.GetStructure(structure.ID())
	if !ok {
		t.Fatalf("GetStructure(%v) returned !ok", structure.ID())
	}
	if got != structure {
		t.Fatalf("expected same structure instance, got %v", got)
	}
	if board.structures[pos.X][pos.Y] != structure {
		t.Fatalf("structure grid not updated at %v", pos)
	}
	if storedPos, ok := board.structureTable[structure.ID()]; !ok || storedPos != pos {
		t.Fatalf("structure table mismatch: got %v want %v", storedPos, pos)
	}
}

func TestBoardSetStructureMovesExistingID(t *testing.T) {
	board := newTestBoard()
	id := uuid.New()
	original := newDummyStructureWithID(id)
	moved := newDummyStructureWithID(id)
	from := pkg.Point{X: 0, Y: 0}
	to := pkg.Point{X: 2, Y: 1}

	if !board.SetStructure(from, original) {
		t.Fatalf("SetStructure(%v) failed", from)
	}
	if !board.SetStructure(to, moved) {
		t.Fatalf("SetStructure(%v) failed", to)
	}

	if board.structures[from.X][from.Y] != nil {
		t.Fatalf("original cell %v was not cleared", from)
	}
	if board.structures[to.X][to.Y] != moved {
		t.Fatalf("new cell %v missing moved structure", to)
	}
	if storedPos := board.structureTable[id]; storedPos != to {
		t.Fatalf("structure table not updated, got %v want %v", storedPos, to)
	}
}

func TestBoardUpdateStructure(t *testing.T) {
	board := newTestBoard()
	pos := pkg.Point{X: 2, Y: 2}
	initial := newDummyStructure()
	if !board.SetStructure(pos, initial) {
		t.Fatalf("SetStructure(%v) failed", pos)
	}

	updated := newDummyStructureWithID(initial.ID())
	if !board.UpdateStructure(updated) {
		t.Fatalf("UpdateStructure should succeed for existing structure")
	}
	if board.structures[pos.X][pos.Y] != updated {
		t.Fatalf("structure at %v not updated", pos)
	}

	missing := newDummyStructure()
	if board.UpdateStructure(missing) {
		t.Fatalf("UpdateStructure should fail for unknown structure")
	}
}

func TestBoardSetAndGetPiece(t *testing.T) {
	board := newTestBoard()
	owner := board.Players()[0]
	pos := pkg.Point{X: 2, Y: 0}
	piece := NewBasePiece(owner, pkg.Point{X: 0, Y: 0}, 10, 5, 3, nil, nil)

	if !board.SetPiece(pos, piece) {
		t.Fatalf("SetPiece(%v) returned false", pos)
	}
	got, ok := board.GetPiece(piece.Id())
	if !ok {
		t.Fatalf("GetPiece(%v) returned !ok", piece.Id())
	}
	if got != piece {
		t.Fatalf("expected same piece pointer")
	}
	if board.pieces[pos.X][pos.Y] != piece {
		t.Fatalf("piece grid not updated at %v", pos)
	}
	if storedPos, ok := board.pieceTable[piece.Id()]; !ok || storedPos != pos {
		t.Fatalf("piece table mismatch: got %v want %v", storedPos, pos)
	}
	if piece.Position() != pos {
		t.Fatalf("SetPiece should sync MutablePiece position, got %v want %v", piece.Position(), pos)
	}
}

func TestBoardSetPieceMovesExistingID(t *testing.T) {
	board := newTestBoard()
	owner := board.Players()[0]
	from := pkg.Point{X: 0, Y: 1}
	to := pkg.Point{X: 1, Y: 1}
	piece := NewBasePiece(owner, from, 8, 4, 2, nil, nil)

	if !board.SetPiece(from, piece) {
		t.Fatalf("SetPiece(%v) failed", from)
	}
	clone := piece.Copy()
	if !board.SetPiece(to, clone) {
		t.Fatalf("SetPiece(%v) with clone failed", to)
	}

	if board.pieces[from.X][from.Y] != nil {
		t.Fatalf("old cell %v was not cleared", from)
	}
	if board.pieces[to.X][to.Y] != clone {
		t.Fatalf("new cell %v missing clone", to)
	}
	if storedPos := board.pieceTable[piece.Id()]; storedPos != to {
		t.Fatalf("piece table not updated, got %v want %v", storedPos, to)
	}
	if clone.Position() != to {
		t.Fatalf("MutablePiece position not updated, got %v want %v", clone.Position(), to)
	}
}

func TestBoardUpdatePiece(t *testing.T) {
	board := newTestBoard()
	owner := board.Players()[0]
	pos := pkg.Point{X: 1, Y: 2}
	piece := NewBasePiece(owner, pos, 12, 6, 4, nil, nil)
	if !board.SetPiece(pos, piece) {
		t.Fatalf("SetPiece(%v) failed", pos)
	}

	updated := piece.Copy()
	updatedMutable, ok := updated.(MutablePiece)
	if !ok {
		t.Fatalf("Copy should return MutablePiece")
	}
	updatedMutable.SetHP(99)

	if !board.UpdatePiece(updated) {
		t.Fatalf("UpdatePiece should succeed for existing piece")
	}
	stored, ok := board.GetPiece(piece.Id())
	if !ok || stored != updated {
		t.Fatalf("piece not updated in board")
	}

	missing := NewBasePiece(owner, pkg.Point{X: 0, Y: 0}, 5, 3, 1, nil, nil)
	if board.UpdatePiece(missing) {
		t.Fatalf("UpdatePiece should fail for unknown piece")
	}
}
