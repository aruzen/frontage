package common

import "testing"

func TestCancelStateToMapFromMap(t *testing.T) {
	state := CancelState{Reason: "blocked"}
	m := state.ToMap()
	if m["reason"] != "blocked" {
		t.Fatalf("reason not serialized: %v", m)
	}
	var dst CancelState
	if err := dst.FromMap(nil, m); err != nil {
		t.Fatalf("FromMap error: %v", err)
	}
	if dst.Reason != "blocked" {
		t.Fatalf("reason not restored: %v", dst.Reason)
	}

	var empty CancelState
	if err := empty.FromMap(nil, map[string]interface{}{}); err != nil {
		t.Fatalf("FromMap empty error: %v", err)
	}
	if empty.Reason != "" {
		t.Fatalf("unexpected reason: %v", empty.Reason)
	}
}
