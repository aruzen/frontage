package model

import (
	"frontage/pkg"
	"github.com/google/uuid"
)

type GenerationStrategy int

const (
	GENERATION_STRATEGY_SWAP GenerationStrategy = iota
	GENERATION_STRATEGY_CHAIN
	//	GENERATION_STRATEGY_MYSELF
)

type BoardInfo struct {
	boardGenerationStrategy GenerationStrategy
	size                    pkg.Size
}

type Board struct {
	info             *BoardInfo
	beforeGeneration *Board
	turn             int
	side             int
	players          [2]Player
	structures       [][]Structure
	structureTable   map[uuid.UUID]pkg.Point
	pieces           [][]Piece
	pieceTable       map[uuid.UUID]pkg.Point
}

func NewBoardInfo(size pkg.Size, strategy GenerationStrategy) *BoardInfo {
	return &BoardInfo{
		boardGenerationStrategy: strategy,
		size:                    size,
	}
}

func NewBoard(info *BoardInfo, players [2]Player) *Board {
	return &Board{
		info:           info,
		players:        players,
		turn:           0,
		structures:     pkg.Make2D[Structure](info.size, nil),
		structureTable: map[uuid.UUID]pkg.Point{},
		pieces:         pkg.Make2D[Piece](info.size, nil),
		pieceTable:     map[uuid.UUID]pkg.Point{},
	}
}

func (b *Board) Copy() *Board {
	return &Board{
		info:           b.info,
		turn:           b.turn,
		side:           b.side,
		players:        [2]Player{b.players[0].Copy(), b.players[1].Copy()},
		structures:     pkg.Copy2D(b.info.size, b.structures),
		structureTable: pkg.CopyMap(b.structureTable),
		pieces:         pkg.Copy2D(b.info.size, b.pieces),
		pieceTable:     pkg.CopyMap(b.pieceTable),
	}
}

func (b *Board) Overwrite(src *Board) {
	b.info = src.info
	b.turn = src.turn
	b.side = src.side
	b.players[0] = src.players[0].Copy()
	b.players[1] = src.players[1].Copy()
	pkg.Overwrite2D(b.pieces, src.pieces)
	b.pieceTable = pkg.CopyMap(src.pieceTable)
	pkg.Overwrite2D(b.structures, src.structures)
	b.structureTable = pkg.CopyMap(src.structureTable)
}

func (b *Board) Sandbox() *Board {
	result := b.Copy()
	result.info = &BoardInfo{
		GENERATION_STRATEGY_CHAIN,
		b.info.size,
	}

	return result
}

func (b *Board) Next() *Board {
	if b.info.boardGenerationStrategy == GENERATION_STRATEGY_CHAIN || b.beforeGeneration == nil {
		return &Board{
			info:             b.info,
			beforeGeneration: b,
			turn:             b.turn,
			side:             b.side,
			players:          [2]Player{b.players[0].Copy(), b.players[1].Copy()},
			structures:       pkg.Copy2D(b.info.size, b.structures),
			structureTable:   pkg.CopyMap(b.structureTable),
			pieces:           pkg.Copy2D(b.info.size, b.pieces),
			pieceTable:       pkg.CopyMap(b.pieceTable),
		}
	} else if b.info.boardGenerationStrategy == GENERATION_STRATEGY_SWAP {
		result := b.beforeGeneration
		result.Overwrite(b)
		result.beforeGeneration = b
		b.beforeGeneration = nil
		return result
	}
	return nil
}

func (b *Board) Info() BoardInfo {
	return *b.info
}

func (b *Board) Size() pkg.Size {
	return b.info.size
}

func (b *Board) Turn() int {
	return b.turn
}

func (b *Board) Phase() int {
	return b.side
}

func (b *Board) Players() [2]Player {
	return b.players
}

func (b *Board) FindPlayer(id uuid.UUID) (Player, bool) {
	for _, p := range b.players {
		if p != nil && p.ID() == id {
			return p, true
		}
	}
	return nil, false
}

func (b *Board) Structures() [][]Structure {
	return b.structures
}

func (b *Board) Entities() [][]Piece {
	return b.pieces
}

func (b *Board) GetStructure(id uuid.UUID) (Structure, bool) {
	if b == nil {
		return nil, false
	}
	pos, ok := b.structureTable[id]
	if !ok || !b.inBounds(pos) {
		if ok {
			delete(b.structureTable, id)
		}
		return nil, false
	}
	structure := b.structures[pos.X][pos.Y]
	if structure == nil || structure.ID() != id {
		delete(b.structureTable, id)
		return nil, false
	}
	return structure, true
}

func (b *Board) GetPiece(id uuid.UUID) (Piece, bool) {
	if b == nil {
		return nil, false
	}
	pos, ok := b.pieceTable[id]
	if !ok || !b.inBounds(pos) {
		if ok {
			delete(b.pieceTable, id)
		}
		return nil, false
	}
	piece := b.pieces[pos.X][pos.Y]
	if piece == nil || piece.Id() != id {
		delete(b.pieceTable, id)
		return nil, false
	}
	return piece, true
}

func (b *Board) SetStructure(pos pkg.Point, structure Structure) bool {
	if b == nil || !b.inBounds(pos) {
		return false
	}
	if current := b.structures[pos.X][pos.Y]; current != nil {
		delete(b.structureTable, current.ID())
	}
	b.structures[pos.X][pos.Y] = structure
	if structure == nil {
		return true
	}
	id := structure.ID()
	if prevPos, ok := b.structureTable[id]; ok && (prevPos != pos) {
		if b.inBounds(prevPos) {
			if prev := b.structures[prevPos.X][prevPos.Y]; prev != nil && prev.ID() == id {
				b.structures[prevPos.X][prevPos.Y] = nil
			}
		}
	}
	b.structureTable[id] = pos
	return true
}

func (b *Board) RemoveStructure(pos pkg.Point) bool {
	return b.SetStructure(pos, nil)
}

// SetPiece Move操作について, 元々配置されているPiece(UUIDで判別)がSetされる時元いた位置にnilを代入する
func (b *Board) SetPiece(pos pkg.Point, piece Piece) bool {
	if b == nil || !b.inBounds(pos) {
		return false
	}
	if current := b.pieces[pos.X][pos.Y]; current != nil {
		delete(b.pieceTable, current.Id())
	}
	b.pieces[pos.X][pos.Y] = piece
	if piece == nil {
		return true
	}
	id := piece.Id()
	if prevPos, ok := b.pieceTable[id]; ok && (prevPos != pos) {
		if b.inBounds(prevPos) {
			if prev := b.pieces[prevPos.X][prevPos.Y]; prev != nil && prev.Id() == id {
				b.pieces[prevPos.X][prevPos.Y] = nil
			}
		}
	}
	b.pieceTable[id] = pos
	if mutable, ok := piece.(MutablePiece); ok {
		mutable.SetPosition(pos)
	}
	return true
}

func (b *Board) RemovePiece(pos pkg.Point) bool {
	return b.SetPiece(pos, nil)
}

func (b *Board) UpdateStructure(structure Structure) bool {
	if b == nil || structure == nil {
		return false
	}
	id := structure.ID()
	pos, ok := b.structureTable[id]
	if !ok || !b.inBounds(pos) {
		if ok {
			delete(b.structureTable, id)
		}
		return false
	}
	return b.SetStructure(pos, structure)
}

func (b *Board) UpdatePiece(piece Piece) bool {
	if b == nil || piece == nil {
		return false
	}
	id := piece.Id()
	pos, ok := b.pieceTable[id]
	if !ok || !b.inBounds(pos) {
		if ok {
			delete(b.pieceTable, id)
		}
		return false
	}
	return b.SetPiece(pos, piece)
}

func (b *Board) ActivePlayer() Player {
	return b.players[0]
}

func (b *Board) inBounds(pos pkg.Point) bool {
	if b == nil {
		return false
	}
	return pos.X >= 0 && pos.X < b.info.size.Width && pos.Y >= 0 && pos.Y < b.info.size.Height
}
