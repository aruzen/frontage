package model

import (
	"testing"
)

func TestDrawEmptyDeckReturnsNoCard(t *testing.T) {
	for _, deck := range []Cards{nil, {}} {
		if card := deck.Draw(); card != nil || len(deck) != 0 {
			t.Fatalf("empty draw = %v; deck length = %d", card, len(deck))
		}
	}
}

func TestDrawLastCardThenEmpty(t *testing.T) {
	card := NewBaseCard("last", nil)
	deck := Cards{card}
	if got := deck.Draw(); got != card || len(deck) != 0 {
		t.Fatalf("last draw = %v; deck length = %d", got, len(deck))
	}
	if got := deck.Draw(); got != nil {
		t.Fatalf("draw after exhaustion returned %v", got)
	}
}
