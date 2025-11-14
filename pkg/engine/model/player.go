package model

const DEFAULT_GOD_BLESSING_POINTS = 3

type DeckType int

const (
	DECK_TYPE_MAIN DeckType = iota
	DECK_TYPE_SUB
	DECK_TYPE_HAND
	DECK_TYPE_GRAVEYARD
	DECK_TYPE_EXTRA
)

type Player struct {
	beforePlayer                              *Player
	generationStrategy                        GenerationStrategy
	GodBlessing                               int
	MainDeck, SubDeck, Hand, Graveyard, Extra Cards
	Materials                                 Materials
}

func NewPlayer(mainDeck Cards, subDeck Cards) *Player {
	return &Player{
		beforePlayer: nil,
		GodBlessing:  DEFAULT_GOD_BLESSING_POINTS,
		MainDeck:     mainDeck,
		SubDeck:      subDeck,
	}
}

func (p *Player) Copy() *Player {
	return &Player{
		beforePlayer: p.beforePlayer,
		GodBlessing:  p.GodBlessing,
		MainDeck:     p.MainDeck,
		SubDeck:      p.SubDeck,
		Hand:         p.Hand,
		Graveyard:    p.Graveyard,
		Extra:        p.Extra,
		Materials:    p.Materials.Copy(),
	}
}

func (p *Player) Decks() []struct {
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

func (p *Player) GetDeck(deckType DeckType) *Cards {
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

func (p *Player) Find(card Card) (DeckType, int, bool) {
	for _, deck := range p.Decks() {
		for i, c := range *deck.Deck {
			if c.Id() == card.Id() {
				return deck.DeckType, i, true
			}
		}
	}
	return 0, 0, false
}

func (p *Player) HasCard(card Card) bool {
	_, _, found := p.Find(card)
	return found
}

func (p *Player) Set(deckType DeckType, idx int, card Card) bool {
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

func (p *Player) Update(card Card) bool {
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
