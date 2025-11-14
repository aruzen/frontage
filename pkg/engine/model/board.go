package model

import (
	"frontage/pkg/engine"
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
	phase            int
	players          [2]*Player
	structures       [][]Structure
	structureTable   map[uuid.UUID]pkg.Point
	entities         [][]Entity
	entityTable      map[uuid.UUID]pkg.Point
}

func NewBoardInfo(size pkg.Size, strategy GenerationStrategy) *BoardInfo {
	return &BoardInfo{
		boardGenerationStrategy: strategy,
		size:                    size,
	}
}

func NewBoard(info *BoardInfo) *Board {
	return &Board{
		info: info,
		players: [2]*Player{
			// TODO
			NewPlayer(make(Cards, 0), make(Cards, 0)),
			NewPlayer(make(Cards, 0), make(Cards, 0)),
		},
		turn:           0,
		structures:     pkg.Make2D[Structure](info.size, nil),
		structureTable: map[uuid.UUID]pkg.Point{},
		entities:       pkg.Make2D[Entity](info.size, nil),
		entityTable:    map[uuid.UUID]pkg.Point{},
	}
}

func (b *Board) Copy() *Board {
	return &Board{
		info:           b.info,
		turn:           b.turn,
		phase:          b.phase,
		players:        [2]*Player{b.players[0].Copy(), b.players[1].Copy()},
		structures:     pkg.Copy2D(b.info.size, b.structures),
		structureTable: pkg.CopyMap(b.structureTable),
		entities:       pkg.Copy2D(b.info.size, b.entities),
		entityTable:    pkg.CopyMap(b.entityTable),
	}
}

func (b *Board) Overwrite(src *Board) {
	b.info = src.info
	b.turn = src.turn
	b.phase = src.phase
	b.players[0] = src.players[0].Copy()
	b.players[1] = src.players[1].Copy()
	pkg.Overwrite2D(b.entities, src.entities)
	b.entityTable = pkg.CopyMap(src.entityTable)
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
			phase:            b.phase,
			players:          [2]*Player{b.players[0].Copy(), b.players[1].Copy()},
			structures:       pkg.Copy2D(b.info.size, b.structures),
			structureTable:   pkg.CopyMap(b.structureTable),
			entities:         pkg.Copy2D(b.info.size, b.entities),
			entityTable:      pkg.CopyMap(b.entityTable),
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
	return b.phase
}

func (b *Board) Players() [2]*Player {
	return b.players
}

func (b *Board) Structures() [][]Structure {
	return b.structures
}

func (b *Board) Entities() [][]Entity {
	return b.entities
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

func (b *Board) GetEntity(id uuid.UUID) (Entity, bool) {
	if b == nil {
		return nil, false
	}
	pos, ok := b.entityTable[id]
	if !ok || !b.inBounds(pos) {
		if ok {
			delete(b.entityTable, id)
		}
		return nil, false
	}
	entity := b.entities[pos.X][pos.Y]
	if entity == nil || entity.Id() != id {
		delete(b.entityTable, id)
		return nil, false
	}
	return entity, true
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

// SetEntity Move操作について, 元々配置されているEntity(UUIDで判別)がSetされる時元いた位置にnilを代入する
func (b *Board) SetEntity(pos pkg.Point, entity Entity) bool {
	if b == nil || !b.inBounds(pos) {
		return false
	}
	if current := b.entities[pos.X][pos.Y]; current != nil {
		delete(b.entityTable, current.Id())
	}
	b.entities[pos.X][pos.Y] = entity
	if entity == nil {
		return true
	}
	id := entity.Id()
	if prevPos, ok := b.entityTable[id]; ok && (prevPos != pos) {
		if b.inBounds(prevPos) {
			if prev := b.entities[prevPos.X][prevPos.Y]; prev != nil && prev.Id() == id {
				b.entities[prevPos.X][prevPos.Y] = nil
			}
		}
	}
	b.entityTable[id] = pos
	if mutable, ok := entity.(MutableEntity); ok {
		mutable.SetPosition(pos)
	}
	return true
}

func (b *Board) RemoveEntity(pos pkg.Point) bool {
	return b.SetEntity(pos, nil)
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

func (b *Board) UpdateEntity(entity Entity) bool {
	if b == nil || entity == nil {
		return false
	}
	id := entity.Id()
	pos, ok := b.entityTable[id]
	if !ok || !b.inBounds(pos) {
		if ok {
			delete(b.entityTable, id)
		}
		return false
	}
	return b.SetEntity(pos, entity)
}

func (b *Board) ActivePlayer() *Player {
	return b.players[0]
}

func (b *Board) inBounds(pos pkg.Point) bool {
	if b == nil {
		return false
	}
	return pos.X >= 0 && pos.X < b.info.size.Width && pos.Y >= 0 && pos.Y < b.info.size.Height
}
