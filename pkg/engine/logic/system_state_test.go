package logic

import "testing"

func TestSystemNoticeStateToMapFromMap(t *testing.T) {
	state := SystemNoticeState{}
	m := state.ToMap()
	if len(m) != 0 {
		t.Fatalf("expected empty map, got %v", m)
	}
	var dst SystemNoticeState
	if err := dst.FromMap(nil, m); err != nil {
		t.Fatalf("FromMap error: %v", err)
	}
}
