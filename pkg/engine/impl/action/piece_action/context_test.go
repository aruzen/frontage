package piece_action_test

import (
	"testing"

	"frontage/internal/event"
	"frontage/pkg"
	"frontage/pkg/engine/impl/action/piece_action"
)

func TestPieceActionContextsToMapFromMap(t *testing.T) {
	// Summon
	{
		base := event.NewBaseEffectContext()
		base.Cancel()
		ctx := piece_action.PieceSummonActionContext{BaseEffectContext: base, Point: pkg.Point{X: 1, Y: 2}}
		m := ctx.ToMap()
		if m["canceled"] != true {
			t.Fatalf("summon canceled flag missing")
		}
		point := m["point"].(map[string]interface{})
		if point["x"] != 1 || point["y"] != 2 {
			t.Fatalf("summon point not serialized correctly: %v", point)
		}
		dst := &piece_action.PieceSummonActionContext{}
		if err := dst.FromMap(m); err != nil {
			t.Fatalf("summon FromMap error: %v", err)
		}
		if dst.Point != (pkg.Point{1, 2}) || !dst.IsCanceled() {
			t.Fatalf("summon context not restored: %+v", dst)
		}
	}

	// Move
	{
		base := event.NewBaseEffectContext()
		ctx := piece_action.PieceMoveActionContext{BaseEffectContext: base, Point: pkg.Point{X: 3, Y: 4}}
		m := ctx.ToMap()
		point := m["point"].(map[string]interface{})
		if point["x"] != 3 || point["y"] != 4 {
			t.Fatalf("move point not serialized correctly: %v", point)
		}
		dst := &piece_action.PieceMoveActionContext{}
		if err := dst.FromMap(m); err != nil {
			t.Fatalf("move FromMap error: %v", err)
		}
		if dst.Point != (pkg.Point{3, 4}) {
			t.Fatalf("move context not restored: %+v", dst)
		}
	}

	// Attack
	{
		base := event.NewBaseEffectContext()
		base.Cancel()
		ctx := piece_action.PieceAttackActionContext{BaseEffectContext: base, Point: pkg.Point{X: 5, Y: 6}, Value: 10}
		m := ctx.ToMap()
		point := m["point"].(map[string]interface{})
		if point["x"] != 5 || point["y"] != 6 || m["value"] != 10 {
			t.Fatalf("attack serialization incorrect: %v", m)
		}
		dst := &piece_action.PieceAttackActionContext{}
		if err := dst.FromMap(m); err != nil {
			t.Fatalf("attack FromMap error: %v", err)
		}
		if dst.Point != (pkg.Point{5, 6}) || dst.Value != 10 || !dst.IsCanceled() {
			t.Fatalf("attack context not restored: %+v", dst)
		}
	}
}
