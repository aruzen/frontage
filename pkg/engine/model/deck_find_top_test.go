package model

import (
	"reflect"
	"slices"
	"testing"
)

func TestFindTopUsesDeckTop(t *testing.T) {
	bottom := NewBaseCard("bottom", nil)
	middle := NewBaseCard("middle", nil)
	top := NewBaseCard("top", nil)
	deck := Cards{bottom, middle}
	deck.PushTop(Cards{top})
	for _, tc := range []struct {
		name string
		pred func(Card) bool
		want Card
	}{
		{"all match", func(Card) bool { return true }, top},
		{"skip top", func(c Card) bool { return c != top }, middle},
		{"bottom only", func(c Card) bool { return c == bottom }, bottom},
		{"no match", func(Card) bool { return false }, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			before := slices.Clone(deck)
			found, ok := deck.FindTop(tc.pred)
			if found != tc.want || ok != (tc.want != nil) {
				t.Fatalf("FindTop = %v, %v; want %v", found, ok, tc.want)
			}
			searched, foundBySearch := before.SearchTop(tc.pred)
			if searched != found || foundBySearch != ok {
				t.Fatal("FindTop and SearchTop disagree on the deck top")
			}
			if !reflect.DeepEqual(deck, Cards{bottom, middle, top}) {
				t.Fatal("FindTop changed the deck")
			}
		})
	}
	var empty Cards
	if found, ok := empty.FindTop(func(Card) bool { t.Fatal("empty deck evaluated predicate"); return true }); found != nil || ok {
		t.Fatal("empty deck returned a card")
	}
}
