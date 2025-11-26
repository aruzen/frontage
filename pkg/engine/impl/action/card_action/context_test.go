package card_action_test

import (
	"testing"

	"frontage/internal/event"
	"frontage/pkg/engine/impl/action/card_action"
)

func TestCardPieceContextsToMapFromMap(t *testing.T) {
	tests := []struct {
		name  string
		make  func() interface{ ToMap() map[string]interface{} }
		load  func(map[string]interface{}) (bool, error)
		value int
	}{
		{
			name:  "HP",
			value: 5,
			make: func() interface{ ToMap() map[string]interface{} } {
				base := event.NewBaseEffectContext()
				base.Cancel()
				return &card_action.PieceHPContext{BaseEffectContext: &base, Value: 5}
			},
			load: func(m map[string]interface{}) (bool, error) {
				dst := &card_action.PieceHPContext{}
				if err := dst.FromMap(m); err != nil {
					return false, err
				}
				return dst.Value == 5 && dst.IsCanceled(), nil
			},
		},
		{
			name:  "MP",
			value: 7,
			make: func() interface{ ToMap() map[string]interface{} } {
				base := event.NewBaseEffectContext()
				base.Cancel()
				return &card_action.PieceMPContext{BaseEffectContext: &base, Value: 7}
			},
			load: func(m map[string]interface{}) (bool, error) {
				dst := &card_action.PieceMPContext{}
				if err := dst.FromMap(m); err != nil {
					return false, err
				}
				return dst.Value == 7 && dst.IsCanceled(), nil
			},
		},
		{
			name:  "ATK",
			value: 9,
			make: func() interface{ ToMap() map[string]interface{} } {
				base := event.NewBaseEffectContext()
				base.Cancel()
				return &card_action.PieceATKContext{BaseEffectContext: &base, Value: 9}
			},
			load: func(m map[string]interface{}) (bool, error) {
				dst := &card_action.PieceATKContext{}
				if err := dst.FromMap(m); err != nil {
					return false, err
				}
				return dst.Value == 9 && dst.IsCanceled(), nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.make().ToMap()
			if m["value"] != tt.value {
				t.Fatalf("value not serialized, got %v", m["value"])
			}
			if v, ok := m["canceled"]; !ok || v != true {
				t.Fatalf("canceled flag missing or false: %v", v)
			}
			ok, err := tt.load(m)
			if err != nil {
				t.Fatalf("FromMap error: %v", err)
			}
			if !ok {
				t.Fatalf("FromMap did not restore context correctly")
			}
		})
	}
}
