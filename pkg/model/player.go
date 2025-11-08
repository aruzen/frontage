package model

const DEFAULT_GOD_BLESSING_POINTS = 3

type Player struct {
	GodBlessing                               int
	MainDeck, SubDeck, Hand, Graveyard, Extra Cards
	Materials                                 Materials
}

func NewPlayer(MainDeck Cards, SubDeck Cards) *Player {
	return &Player{
		GodBlessing: DEFAULT_GOD_BLESSING_POINTS,
		MainDeck:    MainDeck,
		SubDeck:     SubDeck,
	}
}

func (p *Player) Copy() *Player {
	return &Player{
		GodBlessing: p.GodBlessing,
		MainDeck:    p.MainDeck,
		SubDeck:     p.SubDeck,
		Hand:        p.Hand,
		Graveyard:   p.Graveyard,
		Extra:       p.Extra,
		Materials:   p.Materials.Copy(),
	}
}
