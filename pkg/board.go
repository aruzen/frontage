package pkg

import (
	"frontage/pkg/card"
	"frontage/pkg/entity"
	"frontage/pkg/structure"
)

type BoardGenerationStrategy int

const (
	BOARD_GENERATION_STRATEGY_SWAP BoardGenerationStrategy = iota
	BOARD_GENERATION_STRATEGY_CHAIN
)

type BoardInfo struct {
	boardGenerationStrategy BoardGenerationStrategy
	size                    Size
}

type Board struct {
	info             *BoardInfo
	beforeGeneration *Board
	turn             int
	phase            int
	players          [2]*Player
	structures       [][]structure.Structure
	entities         [][]entity.Entity
}

func NewBoardInfo(size Size, strategy BoardGenerationStrategy) *BoardInfo {
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
			NewPlayer(make(card.Cards, 0), make(card.Cards, 0)),
			NewPlayer(make(card.Cards, 0), make(card.Cards, 0)),
		},
		turn:       0,
		structures: Make2D[structure.Structure](info.size, nil),
		entities:   Make2D[entity.Entity](info.size, nil),
	}
}

func (b *Board) Copy() *Board {
	return &Board{
		info:       b.info,
		turn:       b.turn,
		phase:      b.phase,
		players:    [2]*Player{b.players[0].Copy(), b.players[1].Copy()},
		structures: Copy2D(b.info.size, b.structures),
		entities:   Copy2D(b.info.size, b.entities),
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
		structures:       Copy2D(b.info.size, b.structures),
		entities:         Copy2D(b.info.size, b.entities),
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

func (b *Board) Size() Size {
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

func (b *Board) Structures() [][]structure.Structure {
	return b.structures
}

func (b *Board) Entities() [][]entity.Entity {
	return b.entities
}

func (b *Board) ActivePlayer() *Player {
	return b.players[0]
}

func (b *Board) ReplaceEntity(target entity.Entity) bool {
	if b == nil {
		return false
	}
	return Replace2D[entity.Entity](b.entities, target)
}

func (b *Board) ReplaceStructure(target structure.Structure) bool {
	if b == nil {
		return false
	}
	return Replace2D[structure.Structure](b.structures, target)
}
