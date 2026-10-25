package model

import (
	"reflect"
	"testing"
)

func TestSetPlayCostDoesNotKeepCallersMap(t *testing.T) {
	card := NewBaseCard("cost", nil)
	cost := Materials{PYRO: 2}
	card.SetPlayCost(cost)
	cost[PYRO] = 9
	cost[HYDRO] = 1
	if got := card.PlayCost(); !reflect.DeepEqual(got, Materials{PYRO: 2}) {
		t.Fatalf("changing the caller's map changed card cost: %v", got)
	}
}

func TestPushBottomDoesNotChangeSourceDeck(t *testing.T) {
	for _, spareCapacity := range []bool{false, true} {
		name := "full_capacity"
		if spareCapacity {
			name = "spare_capacity"
		}
		t.Run(name, func(t *testing.T) {
			a, b, c := NewBaseCard("a", nil), NewBaseCard("b", nil), NewBaseCard("c", nil)
			d := NewBaseCard("d", nil)
			source := Cards{a, b, c}
			incoming := source[:1:1]
			if spareCapacity {
				incoming = source[:1]
			}
			deck := Cards{d}
			deck.PushBottom(incoming)
			if !reflect.DeepEqual(source, Cards{a, b, c}) {
				t.Error("PushBottom overwrote another deck through the input slice")
			}
			if !reflect.DeepEqual(deck, Cards{a, d}) {
				t.Fatal("PushBottom changed card order")
			}
			source[0] = c
			if deck[0] != a {
				t.Error("destination deck still shares the caller's slice")
			}
		})
	}
}

func TestPushBottomHandlesItsOwnCards(t *testing.T) {
	a, b := NewBaseCard("a", nil), NewBaseCard("b", nil)
	deck := Cards{a, b}
	deck.PushBottom(deck)
	if !reflect.DeepEqual(deck, Cards{a, b, a, b}) {
		t.Fatal("PushBottom lost or reordered overlapping input")
	}
}
