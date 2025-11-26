package translator

import (
	"frontage/pkg"
	"frontage/pkg/engine/impl/card"
	"frontage/pkg/engine/model"
	"frontage/pkg/network/data"
	"frontage/pkg/network/repository"
)

type CardTranslator struct {
	cardRepo      *repository.CardRepository
	materialsTran *MaterialsTranslator
}

func (cr *CardTranslator) ToPieceModel(d data.PieceCard) (card.MutablePiece, error) {
	m, err := cr.cardRepo.Find(pkg.LocalizeTag(d.Tag))
	if err != nil {
		return nil, err
	}
	p, ok := m.(card.Piece)
	if !ok {
		return nil, ErrBadCast
	}
	material, err := cr.materialsTran.ToModel(d.Cost)
	if err != nil {
		return nil, err
	}

	piece := p.Mirror(d.UUID)
	piece.SetPlaced(model.Placed(d.Placed))
	piece.SetPlayCost(material)
	piece.SetATK(d.Atk)
	piece.SetHP(d.Hp)
	piece.SetMP(d.Mp)
	return piece, nil
}

func (cr *CardTranslator) ToPieceData(m card.Piece) (data.PieceCard, error) {
	material, err := cr.materialsTran.ToData(m.PlayCost())
	if err != nil {
		return data.PieceCard{}, err
	}

	return data.PieceCard{
		Card: data.Card{
			Placed: int(m.Placed()),
			Type:   int(m.Type()),
			Tag:    string(m.LocalizeTag()),
			UUID:   m.Id(),
			Cost:   material,
		},
		Atk: m.ATK(),
		Hp:  m.HP(),
		Mp:  m.MP(),
	}, nil
}
