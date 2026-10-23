package pve_test

import (
	"testing"

	"frontage/pkg"
	"frontage/pkg/engine/impl/card"
	"frontage/pkg/engine/model"
	"frontage/pkg/network/controller/pve"
	"frontage/pkg/network/controller/pvp"
	"frontage/pkg/network/repository"
	"github.com/google/uuid"
)

var demoTags = []pkg.LocalizeTag{"少年グルーシャ", "爛漫に咲く花・姫百子", "生まれたばかりの灯・ベビードレイク", "旗将炎猿・ドモルドス"}

func demoRepository(t *testing.T) *repository.CardRepository {
	t.Helper()
	repo := repository.NewCardRepository()
	for _, tag := range demoTags {
		piece := card.NewBasePiece(model.NewBaseCard(tag, model.Materials{model.PYRO: 2}), 4, 3, 2, []pkg.Point{{X: 1, Y: 0}}, nil)
		if err := repo.Insert(piece); err != nil {
			t.Fatal(err)
		}
	}
	return repo
}

func TestDemoDecksHaveEightRealCards(t *testing.T) {
	repo := demoRepository(t)
	for _, build := range []func() (model.Cards, model.Cards, model.Cards, model.Cards){
		func() (model.Cards, model.Cards, model.Cards, model.Cards) {
			return pve.DemoDecks(&pve.RequireContents{CardRepo: repo})
		},
		func() (model.Cards, model.Cards, model.Cards, model.Cards) {
			return pvp.DemoDecks(&pvp.RequireContents{CardRepo: repo})
		},
	} {
		a, b, c, d := build()
		for i, deck := range []model.Cards{a, b, c, d} {
			if len(deck) != 8 {
				t.Fatalf("deck %d contains %d cards; want 8", i, len(deck))
			}
			counts := make(map[pkg.LocalizeTag]int)
			for _, instance := range deck {
				if instance == nil {
					t.Fatal("deck contains a nil card")
				}
				counts[instance.LocalizeTag()]++
				piece, ok := instance.(card.Piece)
				if !ok || piece.HP() != 4 || piece.MP() != 3 || piece.ATK() != 2 || piece.PlayCost()[model.PYRO] != 2 {
					t.Fatal("new card lost its definition")
				}
			}
			group := (i / 2) * 2
			if counts[demoTags[group]] != 4 || counts[demoTags[group+1]] != 4 {
				t.Fatalf("wrong deck composition: %v", counts)
			}
		}
	}
}

func TestDemoDeckCardsHaveIndependentIDs(t *testing.T) {
	repo := demoRepository(t)
	seen := make(map[uuid.UUID]bool)
	for _, tag := range demoTags {
		template, _ := repo.Find(tag)
		seen[template.Id()] = true
	}
	for run := 0; run < 2; run++ {
		for _, build := range []func() (model.Cards, model.Cards, model.Cards, model.Cards){
			func() (model.Cards, model.Cards, model.Cards, model.Cards) {
				return pve.DemoDecks(&pve.RequireContents{CardRepo: repo})
			},
			func() (model.Cards, model.Cards, model.Cards, model.Cards) {
				return pvp.DemoDecks(&pvp.RequireContents{CardRepo: repo})
			},
		} {
			a, b, c, d := build()
			for _, deck := range []model.Cards{a, b, c, d} {
				for _, instance := range deck {
					if instance == nil {
						continue
					}
					if instance.Id() == uuid.Nil || seen[instance.Id()] {
						t.Fatalf("reused card ID %s", instance.Id())
					}
					seen[instance.Id()] = true
				}
			}
		}
	}
}
