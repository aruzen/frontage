package event_test

import (
	"testing"

	"frontage/internal/event"
)

func TestBaseEffectContextToMapFromMap(t *testing.T) {
	base := event.NewBaseEffectContext()
	if base.IsCanceled() {
		t.Fatalf("new context should not be canceled")
	}
	base.Cancel()
	m := base.ToMap()
	if v, ok := m["canceled"]; !ok || v != true {
		t.Fatalf("expected canceled=true in map, got %v", m["canceled"])
	}

	var dst event.BaseEffectContext
	if err := dst.FromMap(m); err != nil {
		t.Fatalf("FromMap error: %v", err)
	}
	if !dst.IsCanceled() {
		t.Fatalf("FromMap should restore canceled state")
	}
}
