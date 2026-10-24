package model

import (
	"fmt"
	"reflect"
	"slices"
	"testing"
)

func TestFindZeroCountPreservesDeck(t *testing.T) {
	deck := Cards{NewBaseCard("a", nil), NewBaseCard("b", nil)}
	before := slices.Clone(deck)
	calls := 0
	found, ok := deck.Find(0, func(Card) bool { calls++; return true })
	if len(found) != 0 || ok || calls != 0 || !reflect.DeepEqual(deck, before) {
		t.Fatalf("Find(0) = %v, %v; predicate calls = %d; deck = %v", found, ok, calls, deck)
	}
	if found, ok := deck.Find(1, func(Card) bool { return true }); len(found) != 1 || found[0] != deck[0] || !ok {
		t.Fatal("positive count no longer finds the requested card")
	}
}

func TestRemoveNonPositiveCountPreservesDeck(t *testing.T) {
	for _, count := range []int{0, -1, -100} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			deck := Cards{NewBaseCard("a", nil), NewBaseCard("b", nil)}
			before := slices.Clone(deck)
			calls := 0
			removed, ok := deck.Remove(count, func(Card) bool { calls++; return true })
			if removed != 0 || ok || calls != 0 || !reflect.DeepEqual(deck, before) {
				t.Fatalf("Remove(%d) = %d, %v; predicate calls = %d; deck = %v", count, removed, ok, calls, deck)
			}
			if removed, ok := deck.Remove(1, func(Card) bool { return true }); removed != 1 || !ok || len(deck) != 1 || deck[0] != before[0] {
				t.Fatal("positive count no longer removes exactly the top card")
			}
		})
	}
}
