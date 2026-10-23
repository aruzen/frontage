package model

import (
	"fmt"
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
	piece := NewBasePiece(uuid.New(), owner, pkg.Point{X: 0, Y: 0}, 10, 5, 3, nil, nil)

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
	piece := NewBasePiece(uuid.New(), owner, from, 8, 4, 2, nil, nil)

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
	piece := NewBasePiece(uuid.New(), owner, pos, 12, 6, 4, nil, nil)
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

	missing := NewBasePiece(uuid.New(), owner, pkg.Point{X: 0, Y: 0}, 5, 3, 1, nil, nil)
	if board.UpdatePiece(missing) {
		t.Fatalf("UpdatePiece should fail for unknown piece")
	}
}

func TestBoardCopiesIsolateMutableState(t *testing.T) {
	for _, strategy := range []GenerationStrategy{GENERATION_STRATEGY_CHAIN, GENERATION_STRATEGY_SWAP} {
		t.Run(fmt.Sprint(strategy), func(t *testing.T) {
			board := NewBoard(NewBoardInfo(pkg.Size{Width: 3, Height: 3}, strategy),
				[2]Player{NewLocalPlayer(uuid.New(), Cards{NewBaseCard("test", nil)}, nil), NewProxyPlayer(uuid.New())})
			piece := NewBasePiece(uuid.New(), board.Players()[0], pkg.Point{}, 10, 5, 3, nil, nil)
			board.SetPiece(pkg.Point{}, piece)
			for _, copyBoard := range []struct {
				name string
				copy func() *Board
			}{{"Copy", board.Copy}, {"Sandbox", board.Sandbox}, {"Next", board.Next}} {
				t.Run(copyBoard.name, func(t *testing.T) {
					copied := copyBoard.copy()
					copiedPiece, _ := copied.GetPiece(piece.Id())
					copiedPiece.(MutablePiece).SetHP(7)
					copied.SetPiece(pkg.Point{X: 1, Y: 1}, copiedPiece)
					copiedPiece.Owner().SetGodBlessing(0)
					deck := copied.Players()[0].(DeckPlayer).GetDeck(DECK_TYPE_MAIN)
					(*deck)[0].(MutableCard).SetPlayCost(Materials{PYRO: 9})
					if piece.HP() != 10 || piece.Position() != (pkg.Point{}) {
						t.Fatalf("copy changed original piece: HP=%d, position=%v", piece.HP(), piece.Position())
					}
					if board.Players()[0].GodBlessing() != DEFAULT_GOD_BLESSING_POINTS {
						t.Fatal("copy changed original player through piece owner")
					}
					if copiedPiece.Owner() != copied.Players()[0] {
						t.Fatal("copied piece owner is not the copied player")
					}
					originalDeck := board.Players()[0].(DeckPlayer).GetDeck(DECK_TYPE_MAIN)
					if len((*originalDeck)[0].PlayCost()) != 0 {
						t.Fatal("copy changed original card")
					}
				})
			}
			// Exercise the reuse branch of the swap generation strategy too.
			next := board.Next()
			latest := next.Next()
			latestPiece, _ := latest.GetPiece(piece.Id())
			latestPiece.(MutablePiece).SetHP(1)
			previousPiece, _ := next.GetPiece(piece.Id())
			if previousPiece.HP() != 10 {
				t.Fatal("reused board generation shares mutable pieces with its predecessor")
			}
		})
	}
}

func TestBoardAdvanceTurnAlternatesActivePlayer(t *testing.T) {
	board := newTestBoard()
	players := board.Players()
	for turn := 0; turn < 6; turn++ {
		if board.Turn() != turn || board.ActivePlayer().ID() != players[turn%2].ID() {
			t.Fatalf("turn %d: got turn %d, player %v", turn, board.Turn(), board.ActivePlayer().ID())
		}
		board = board.AdvanceTurn()
	}
}
