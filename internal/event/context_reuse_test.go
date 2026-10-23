package event_test

import (
	"testing"

	"frontage/internal/event"
)

func TestContextRestoreDoesNotRetainPreviousCancellation(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input map[string]interface{}
	}{
		{"nil map", nil},
		{"omitted", map[string]interface{}{}},
		{"null", map[string]interface{}{"canceled": nil}},
		{"false", map[string]interface{}{"canceled": false}},
		{"true", map[string]interface{}{"canceled": true}},
		{"invalid", map[string]interface{}{"canceled": "true"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reused := event.NewBaseEffectContext()
			reused.Cancel()
			var fresh event.BaseEffectContext
			freshErr := fresh.FromMap(tc.input)
			reusedErr := reused.FromMap(tc.input)
			if (freshErr == nil) != (reusedErr == nil) || fresh.IsCanceled() != reused.IsCanceled() {
				t.Fatalf("restore depends on previous cancellation: fresh=%v/%v reused=%v/%v", fresh.IsCanceled(), freshErr, reused.IsCanceled(), reusedErr)
			}
		})
	}
}
