package card

import (
	"frontage/pkg"
	"math/rand"
)

type Cards []*Card

func (cs *Cards) Shuffle(rand rand.Rand) {
	rand.Shuffle(len(*cs), func(i, j int) {
		(*cs)[i], (*cs)[j] = (*cs)[j], (*cs)[i]
	})
}

func (cs Cards) Find(count uint, pred func(*Card) bool) (Cards, bool) {
	result := make(Cards, 0)
	for _, card := range cs {
		if pred(card) {
			result = append(result, card)
			count--
			if count == 0 {
				return result, true
			}
		}
	}
	return result, len(result) > 0
}

func (cs Cards) FindTop(pred func(*Card) bool) (*Card, bool) {
	for _, card := range cs {
		if pred(card) {
			return card, true
		}
	}
	return nil, false
}

func (cs Cards) FindAll(pred func(*Card) bool) (Cards, bool) {
	result := make(Cards, 0)
	for _, card := range cs {
		if pred(card) {
			result = append(result, card)
		}
	}
	return result, len(result) > 0
}

func (cs *Cards) RemoveCard(target *Card) bool {
	for i := len(*cs) - 1; i >= 0; i-- {
		if (*cs)[i] == target {
			(*cs) = append((*cs)[:i], (*cs)[i+1:]...)
			return true
		}
	}
	return false
}

func (cs *Cards) RemoveCards(targets Cards) (int, bool) {
	count := 0
	for i := len(*cs) - 1; i >= 0; i-- {
		for _, t := range targets {
			if (*cs)[i] == t {
				(*cs) = append((*cs)[:i], (*cs)[i+1:]...)
				count++
			}
		}
	}
	return count, len(*cs) > 0
}

func (cs *Cards) Remove(require_count int, pred func(*Card) bool) (int, bool) {
	count := 0
	for i := len(*cs) - 1; i >= 0; i-- {
		if pred((*cs)[i]) {
			(*cs) = append((*cs)[:i], (*cs)[i+1:]...)
			count++
			if require_count == count {
				return count, true
			}
		}
	}
	return count, false
}

func (cs *Cards) RemoveTop(pred func(*Card) bool) bool {
	for i := len(*cs) - 1; i >= 0; i-- {
		if pred((*cs)[i]) {
			(*cs) = append((*cs)[:i], (*cs)[i+1:]...)
			return true
		}
	}
	return false
}

func (cs *Cards) RemoveAll(pred func(*Card) bool) int {
	count := 0
	for i := len(*cs) - 1; i >= 0; i-- {
		if pred((*cs)[i]) {
			(*cs) = append((*cs)[:i], (*cs)[i+1:]...)
			count++
		}
	}
	return count
}

func (cs *Cards) UnorderdRemoveCards(targets Cards) (int, bool) {
	count := 0
	for i := len(*cs) - 1; i >= 0; i-- {
		for _, t := range targets {
			if (*cs)[i] == t {
				(*cs)[i] = (*cs)[len(*cs)-1]
				*cs = (*cs)[:len(*cs)-1]
				count++
			}
		}
	}
	return count, len(*cs) > 0
}

func (cs *Cards) UnorderdRemoveTop(target *Card) {
	for i := len(*cs) - 1; i >= 0; i-- {
		if (*cs)[i] == target {
			(*cs)[i] = (*cs)[len(*cs)-1]
			*cs = (*cs)[:len(*cs)-1]
			return
		}
	}
}

func (cs *Cards) Search(count int, pred func(*Card) bool) Cards {
	var removed Cards
	for i := len(*cs) - 1; i >= 0 && count > 0; i-- {
		if pred((*cs)[i]) {
			removed = append(removed, (*cs)[i])
			*cs = append((*cs)[:i], (*cs)[i+1:]...)
			count--
		}
	}
	return removed
}

func (cs *Cards) SearchTop(pred func(*Card) bool) (*Card, bool) {
	for i := len(*cs) - 1; i >= 0; i-- {
		if pred((*cs)[i]) {
			card := (*cs)[i]
			*cs = append((*cs)[:i], (*cs)[i+1:]...)
			return card, true
		}
	}
	return nil, false
}

func (cs *Cards) SearchAll(pred func(*Card) bool) Cards {
	var removed Cards
	for i := len(*cs) - 1; i >= 0; i-- {
		if pred((*cs)[i]) {
			removed = append(removed, (*cs)[i])
			*cs = append((*cs)[:i], (*cs)[i+1:]...)
		}
	}
	return removed
}

func (cs *Cards) UnorderedSearch(count int, pred func(*Card) bool) Cards {
	var removed Cards
	for i := len(*cs) - 1; i >= 0 && count > 0; i-- {
		if pred((*cs)[i]) {
			removed = append(removed, (*cs)[i])
			(*cs)[i] = (*cs)[len(*cs)-1]
			*cs = (*cs)[:len(*cs)-1]
			count--
		}
	}
	return removed
}

func (cs *Cards) UnorderedSearchAll(pred func(*Card) bool) Cards {
	var removed Cards
	for i := len(*cs) - 1; i >= 0; i-- {
		if pred((*cs)[i]) {
			removed = append(removed, (*cs)[i])
			(*cs)[i] = (*cs)[len(*cs)-1]
			*cs = (*cs)[:len(*cs)-1]
		}
	}
	return removed
}

func (cs *Cards) UnorderedSearchTop(pred func(*Card) bool) (Card, bool) {
	for i := len(*cs) - 1; i >= 0; i-- {
		if pred((*cs)[i]) {
			card := (*cs)[i]
			(*cs)[i] = (*cs)[len(*cs)-1]
			*cs = (*cs)[:len(*cs)-1]
			return *card, true
		}
	}
	return nil, false
}

func (cs *Cards) Draw(rand rand.Rand) Card {
	card := (*cs)[len(*cs)-1]
	*cs = (*cs)[:len(*cs)-1]
	cs.Shuffle(rand)
	return *card
}

func (cs *Cards) PushTop(cards Cards) {
	*cs = append(*cs, cards...)
}

func (cs *Cards) PushBottom(cards Cards) {
	*cs = append(cards, *cs...)
}

func (cs *Cards) Replace(target *Card) bool {
	if cs == nil {
		return false
	}
	slice := []*Card(*cs)
	return pkg.ReplacePtr[Card](slice, target)
}
