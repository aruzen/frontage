package model

import (
	"frontage/pkg"
)

type BoardGenerationStrategy int

const (
	BOARD_GENERATION_STRATEGY_SWAP BoardGenerationStrategy = iota
	BOARD_GENERATION_STRATEGY_CHAIN
)

type BoardInfo struct {
	boardGenerationStrategy BoardGenerationStrategy
	size                    pkg.Size
}

type Board struct {
	info             *BoardInfo
	beforeGeneration *Board
	turn             int
	phase            int
	players          [2]*Player
	structures       [][]Structure
	entities         [][]Entity
}

func NewBoardInfo(size pkg.Size, strategy BoardGenerationStrategy) *BoardInfo {
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
		turn:       0,
		structures: pkg.Make2D[Structure](info.size, nil),
		entities:   pkg.Make2D[Entity](info.size, nil),
	}
}

func (b *Board) Copy() *Board {
	return &Board{
		info:       b.info,
		turn:       b.turn,
		phase:      b.phase,
		players:    [2]*Player{b.players[0].Copy(), b.players[1].Copy()},
		structures: pkg.Copy2D(b.info.size, b.structures),
		entities:   pkg.Copy2D(b.info.size, b.entities),
	}
}

func (b *Board) Sandbox() *Board {
	result := b.Copy()
	result.info = &BoardInfo{
		BOARD_GENERATION_STRATEGY_CHAIN,
		b.info.size,
	}

	return result
}

func (b *Board) Next() *Board {
	newBoard := Board{
		info:             b.info,
		beforeGeneration: b,
		turn:             b.turn,
		phase:            b.phase,
		players:          [2]*Player{b.players[0].Copy(), b.players[1].Copy()},
		structures:       pkg.Copy2D(b.info.size, b.structures),
		entities:         pkg.Copy2D(b.info.size, b.entities),
	}

	if b.info.boardGenerationStrategy == BOARD_GENERATION_STRATEGY_CHAIN || b.beforeGeneration == nil {
		return &newBoard
	} else if b.info.boardGenerationStrategy == BOARD_GENERATION_STRATEGY_SWAP {
		*b.beforeGeneration = newBoard
		return b.beforeGeneration
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

func (b *Board) ActivePlayer() *Player {
	return b.players[0]
}

func (b *Board) ReplaceEntity(old Entity, new Entity) bool {
	if b == nil {
		return false
	}
	return pkg.Replace2D[Entity](b.entities, old, new)
}

func (b *Board) ReplaceStructure(old Structure, new Structure) bool {
	if b == nil {
		return false
	}
	return pkg.Replace2D[Structure](b.structures, old, new)
}
