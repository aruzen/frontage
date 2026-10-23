package logic_test

import (
	"testing"

	"frontage/pkg"
	"frontage/pkg/engine/impl/action/common"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
)

func newCycleSystem() *logic.EventSystem {
	return logic.NewEventSystem(model.NewBoard(
		model.NewBoardInfo(pkg.Size{Width: 2, Height: 2}, model.GENERATION_STRATEGY_CHAIN), [2]model.Player{}))
}

func noticeEvent() *logic.EffectEvent {
	return logic.NewEffectEvent(logic.GameStartAction{}, &logic.SystemNoticeState{})
}

func TestEventCannotChainItselfOrAnAncestor(t *testing.T) {
	for _, scenario := range []string{"self_chain", "self_line", "ancestor"} {
		t.Run(scenario, func(t *testing.T) {
			es := newCycleSystem()
			root := noticeEvent()
			es.Emit(root)
			want := 1
			switch scenario {
			case "self_chain":
				es.Chain(root)
			case "self_line":
				es.ChainLine(root)
			case "ancestor":
				es.ChainLine(noticeEvent(), root)
				want = 2
			}
			if root.Parent() != nil {
				t.Fatal("root became its own descendant; resolving this tree would recurse indefinitely")
			}
			es.Resolve()
			if len(es.AppliedEffects) != want {
				t.Errorf("applied effects = %d; want %d", len(es.AppliedEffects), want)
			}
		})
	}
}

func TestModifierCannotChainItself(t *testing.T) {
	es := newCycleSystem()
	root := noticeEvent()
	es.Emit(root)
	modifier := logic.NewModifyEvent(common.CancelModifyAction{}, &common.CancelState{})
	es.ChainLine(modifier, modifier)
	if modifier.Parent() != root {
		t.Fatal("modifier became its own parent")
	}
	es.Resolve()
}

func TestDistinctEventsMayShareActionAndState(t *testing.T) {
	es := newCycleSystem()
	state := &logic.SystemNoticeState{}
	root := logic.NewEffectEvent(logic.GameStartAction{}, state)
	first := logic.NewEffectEvent(logic.GameStartAction{}, state)
	second := logic.NewEffectEvent(logic.GameStartAction{}, state)
	es.Emit(root)
	es.ChainLine(first, second)
	if first.Parent() != root || second.Parent() != first {
		t.Fatal("valid chain was rejected")
	}
	trial, effects := es.ResolveTest()
	if trial == nil || len(effects) != 3 {
		t.Fatal("valid shared-state chain could not be simulated")
	}
	es.Resolve()
	if len(es.AppliedEffects) != 3 {
		t.Fatal("valid shared-state chain did not resolve")
	}
}

func TestEventCannotBeAttachedToTwoParents(t *testing.T) {
	es := newCycleSystem()
	root, first, second, shared := noticeEvent(), noticeEvent(), noticeEvent(), noticeEvent()
	es.Emit(root)
	es.ChainLine(first, shared)
	es.ChainLine(second, shared, first)
	if shared.Parent() != first || first.Parent() != root {
		t.Fatal("reusing an attached event changed its parent or created a cycle through two branches")
	}
	es.Resolve()
	if len(es.AppliedEffects) != 4 {
		t.Fatalf("applied effects = %d; want each of the 4 distinct events once", len(es.AppliedEffects))
	}
}
