package common

import (
	"frontage/pkg"
	"frontage/pkg/engine/impl/action"
	"frontage/pkg/engine/logic"
)

type CancelState struct {
	Reason string
}

type CancelModifyAction struct {
	logic.BaseAction[CancelState, logic.EffectContext]
}

func (CancelModifyAction) Tag() logic.ModifyActionTag { return action.CANCEL_MODIFY_ACTION }
func (a CancelModifyAction) LocalizeTag() pkg.LocalizeTag {
	return pkg.LocalizeTag(a.Tag())
}

func (m CancelModifyAction) Modify(state interface{}, context logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if context == nil {
		return nil, nil
	}
	context.Cancel()
	if s, ok := state.(CancelState); ok && s.Reason != "" {
		return context, logic.Summary{"reason": s.Reason}
	}
	return context, nil
}
