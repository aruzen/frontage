package model

import "testing"

func TestRemoveCardsDuplicateTargets(t *testing.T) {
	a := NewBaseCard("a", nil)
	b := NewBaseCard("b", nil)
	c := NewBaseCard("c", nil)
	missing := NewBaseCard("missing", nil)
	for _, operation := range []struct {
		name   string
		remove func(*Cards, Cards) (int, bool)
	}{
		{"ordered", (*Cards).RemoveCards},
		{"unordered", (*Cards).UnorderdRemoveCards},
	} {
		t.Run(operation.name, func(t *testing.T) {
			for _, tc := range []struct {
				name                string
				deck, targets, want Cards
				removed             int
			}{
				{"last card", Cards{a}, Cards{a, a}, nil, 1},
				{"keep other cards", Cards{a, b, c}, Cards{b, b, missing}, Cards{a, c}, 1},
				{"multiple targets", Cards{a, b, c}, Cards{a, b, b, c, c}, nil, 3},
				{"no matches", Cards{a, c}, Cards{b, b}, Cards{a, c}, 0},
			} {
				t.Run(tc.name, func(t *testing.T) {
					deck := append(Cards(nil), tc.deck...)
					removed, remaining := operation.remove(&deck, tc.targets)
					if removed != tc.removed || remaining != (len(tc.want) > 0) || len(deck) != len(tc.want) {
						t.Fatalf("got removed=%d remaining=%v deck=%v; want removed=%d deck=%v", removed, remaining, deck, tc.removed, tc.want)
					}
					counts := make(map[Card]int)
					for _, card := range deck {
						counts[card]++
					}
					for i, card := range tc.want {
						if operation.name == "ordered" && deck[i] != card {
							t.Fatal("ordered removal changed the remaining order")
						}
						counts[card]--
					}
					for _, count := range counts {
						if count != 0 {
							t.Fatal("removal changed a non-target card")
						}
					}
				})
			}
		})
	}
}
