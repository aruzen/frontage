package logic_test

import (
	"testing"

	"frontage/pkg/engine/logic"
)

func TestSystemNoticeContextToMapFromMap(t *testing.T) {
	ctx := logic.SystemNoticeContext{}
	m := ctx.ToMap()
	if len(m) != 0 {
		t.Fatalf("SystemNoticeContext should serialize to empty map, got %v", m)
	}
	var dst logic.SystemNoticeContext
	if err := dst.FromMap(m); err != nil {
		t.Fatalf("FromMap error: %v", err)
	}
}
