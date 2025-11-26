package model

import "github.com/google/uuid"

const DEFAULT_GOD_BLESSING_POINTS = 3

type DeckType int

const (
	DECK_TYPE_MAIN DeckType = iota
	DECK_TYPE_SUB
	DECK_TYPE_HAND
	DECK_TYPE_GRAVEYARD
	DECK_TYPE_EXTRA
)

// Player はボード上で扱うプレイヤーの抽象インターフェース。
// デッキ操作を行うものは DeckPlayer を実装してください。
type Player interface {
	ID() uuid.UUID
	Copy() Player
	GodBlessing() int
	SetGodBlessing(int)
	Materials() Materials
	SetMaterials(Materials)
}

type DeckPlayer interface {
	Player
	GetDeck(deckType DeckType) *Cards
}

type ProxyPlayer struct {
	id          uuid.UUID
	godBlessing int
	materials   Materials
}

func NewProxyPlayer(id uuid.UUID) *ProxyPlayer {
	return &ProxyPlayer{
		id:          id,
		godBlessing: DEFAULT_GOD_BLESSING_POINTS,
		materials:   make(Materials),
	}
}

func (p *ProxyPlayer) ID() uuid.UUID {
	return p.id
}

func (p *ProxyPlayer) Copy() Player {
	return &ProxyPlayer{
		id:          p.id,
		godBlessing: p.godBlessing,
		materials:   p.materials.Copy(),
	}
}

func (p *ProxyPlayer) GodBlessing() int {
	return p.godBlessing
}

func (p *ProxyPlayer) SetGodBlessing(v int) {
	p.godBlessing = v
}

func (p *ProxyPlayer) Materials() Materials {
	return p.materials.Copy()
}

func (p *ProxyPlayer) SetMaterials(m Materials) {
	p.materials = m.Copy()
}

// LocalPlayer は従来のフル機能プレイヤー。
type LocalPlayer struct {
	id                 uuid.UUID
	beforePlayer       *LocalPlayer
	generationStrategy GenerationStrategy
	godBlessing        int
	MainDeck, SubDeck  Cards
	Hand, Graveyard    Cards
	Extra              Cards
	materials          Materials
}

var _ Player = &LocalPlayer{}

func NewLocalPlayer(id uuid.UUID, mainDeck Cards, subDeck Cards) *LocalPlayer {
	return &LocalPlayer{
		id:           id,
		beforePlayer: nil,
		godBlessing:  DEFAULT_GOD_BLESSING_POINTS,
		MainDeck:     mainDeck,
		SubDeck:      subDeck,
		materials:    make(Materials),
	}
}

func (p *LocalPlayer) ID() uuid.UUID {
	return p.id
}

func (p *LocalPlayer) Copy() Player {
	return &LocalPlayer{
		id:           p.id,
		beforePlayer: p.beforePlayer,
		godBlessing:  p.godBlessing,
		MainDeck:     p.MainDeck,
		SubDeck:      p.SubDeck,
		Hand:         p.Hand,
		Graveyard:    p.Graveyard,
		Extra:        p.Extra,
		materials:    p.materials.Copy(),
	}
}

func (p *LocalPlayer) Decks() []struct {
	DeckType DeckType
	Deck     *Cards
} {
	return []struct {
		DeckType DeckType
		Deck     *Cards
	}{
		{DECK_TYPE_MAIN, &p.MainDeck},
		{DECK_TYPE_SUB, &p.SubDeck},
		{DECK_TYPE_HAND, &p.Hand},
		{DECK_TYPE_GRAVEYARD, &p.Graveyard},
		{DECK_TYPE_EXTRA, &p.Extra},
	}
}

func (p *LocalPlayer) GetDeck(deckType DeckType) *Cards {
	switch deckType {
	case DECK_TYPE_MAIN:
		return &p.MainDeck
	case DECK_TYPE_SUB:
		return &p.SubDeck
	case DECK_TYPE_HAND:
		return &p.Hand
	case DECK_TYPE_GRAVEYARD:
		return &p.Graveyard
	case DECK_TYPE_EXTRA:
		return &p.Extra
	}
	return nil
}

func (p *LocalPlayer) FindCard(card Card) (DeckType, int, bool) {
	for _, deck := range p.Decks() {
		for i, c := range *deck.Deck {
			if c.Id() == card.Id() {
				return deck.DeckType, i, true
			}
		}
	}
	return 0, 0, false
}

func (p *LocalPlayer) HasCard(card Card) bool {
	_, _, found := p.FindCard(card)
	return found
}

func (p *LocalPlayer) SetCard(deckType DeckType, idx int, card Card) bool {
	deck := p.GetDeck(deckType)
	if deck == nil {
		return false
	}
	if idx < 0 || len(*deck) <= idx {
		return false
	}
	(*deck)[idx] = card
	return true
}

func (p *LocalPlayer) UpdateCard(card Card) bool {
	for _, deck := range p.Decks() {
		idx := -1
		for i, c := range *deck.Deck {
			if c.Id() == card.Id() {
				idx = i
			}
		}
		if idx != -1 {
			(*deck.Deck)[idx] = card
			return true
		}
	}
	return false
}

func (p *LocalPlayer) GodBlessing() int {
	return p.godBlessing
}

func (p *LocalPlayer) SetGodBlessing(v int) {
	p.godBlessing = v
}

func (p *LocalPlayer) Materials() Materials {
	return p.materials.Copy()
}

func (p *LocalPlayer) SetMaterials(m Materials) {
	p.materials = m.Copy()
}
