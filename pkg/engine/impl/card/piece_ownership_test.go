package card

import (
	"reflect"
	"testing"

	"frontage/pkg"
	"frontage/pkg/engine/model"
)

func TestNewPieceDoesNotKeepCallersRanges(t *testing.T) {
	moves := []pkg.Point{{X: 1}}
	ranges := []pkg.Point{{Y: 1}}
	first := NewBasePiece(model.NewBaseCard("first", nil), 3, 1, 2, moves, ranges)
	moves[0] = pkg.Point{X: 2}
	ranges[0] = pkg.Point{Y: 2}
	second := NewBasePiece(model.NewBaseCard("second", nil), 3, 1, 2, moves, ranges)
	if !reflect.DeepEqual(first.LegalMoves(), []pkg.Point{{X: 1}}) ||
		!reflect.DeepEqual(first.AttackRanges(), []pkg.Point{{Y: 1}}) {
		t.Error("reusing a definition buffer changed the first card's movement or attack range")
	}
	if !reflect.DeepEqual(second.LegalMoves(), moves) || !reflect.DeepEqual(second.AttackRanges(), ranges) {
		t.Error("second card did not receive its own ranges")
	}
}
