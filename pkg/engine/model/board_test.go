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
	return NewBoard(info, [2]*Player{{}, {}})
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

func TestBoardSetAndGetEntity(t *testing.T) {
	board := newTestBoard()
	owner := board.Players()[0]
	pos := pkg.Point{X: 2, Y: 0}
	entity := NewBaseEntity(owner, pkg.Point{X: 0, Y: 0}, 10, 5, 3, nil, nil)

	if !board.SetEntity(pos, entity) {
		t.Fatalf("SetEntity(%v) returned false", pos)
	}
	got, ok := board.GetEntity(entity.Id())
	if !ok {
		t.Fatalf("GetEntity(%v) returned !ok", entity.Id())
	}
	if got != entity {
		t.Fatalf("expected same entity pointer")
	}
	if board.entities[pos.X][pos.Y] != entity {
		t.Fatalf("entity grid not updated at %v", pos)
	}
	if storedPos, ok := board.entityTable[entity.Id()]; !ok || storedPos != pos {
		t.Fatalf("entity table mismatch: got %v want %v", storedPos, pos)
	}
	if entity.Position() != pos {
		t.Fatalf("SetEntity should sync MutableEntity position, got %v want %v", entity.Position(), pos)
	}
}

func TestBoardSetEntityMovesExistingID(t *testing.T) {
	board := newTestBoard()
	owner := board.Players()[0]
	from := pkg.Point{X: 0, Y: 1}
	to := pkg.Point{X: 1, Y: 1}
	entity := NewBaseEntity(owner, from, 8, 4, 2, nil, nil)

	if !board.SetEntity(from, entity) {
		t.Fatalf("SetEntity(%v) failed", from)
	}
	clone := entity.Copy()
	if !board.SetEntity(to, clone) {
		t.Fatalf("SetEntity(%v) with clone failed", to)
	}

	if board.entities[from.X][from.Y] != nil {
		t.Fatalf("old cell %v was not cleared", from)
	}
	if board.entities[to.X][to.Y] != clone {
		t.Fatalf("new cell %v missing clone", to)
	}
	if storedPos := board.entityTable[entity.Id()]; storedPos != to {
		t.Fatalf("entity table not updated, got %v want %v", storedPos, to)
	}
	if clone.Position() != to {
		t.Fatalf("MutableEntity position not updated, got %v want %v", clone.Position(), to)
	}
}

func TestBoardUpdateEntity(t *testing.T) {
	board := newTestBoard()
	owner := board.Players()[0]
	pos := pkg.Point{X: 1, Y: 2}
	entity := NewBaseEntity(owner, pos, 12, 6, 4, nil, nil)
	if !board.SetEntity(pos, entity) {
		t.Fatalf("SetEntity(%v) failed", pos)
	}

	updated := entity.Copy()
	updatedMutable, ok := updated.(MutableEntity)
	if !ok {
		t.Fatalf("Copy should return MutableEntity")
	}
	updatedMutable.SetHP(99)

	if !board.UpdateEntity(updated) {
		t.Fatalf("UpdateEntity should succeed for existing entity")
	}
	stored, ok := board.GetEntity(entity.Id())
	if !ok || stored != updated {
		t.Fatalf("entity not updated in board")
	}

	missing := NewBaseEntity(owner, pkg.Point{X: 0, Y: 0}, 5, 3, 1, nil, nil)
	if board.UpdateEntity(missing) {
		t.Fatalf("UpdateEntity should fail for unknown entity")
	}
}
