package common

import (
	"frontage/pkg"
	"frontage/pkg/engine/impl/action"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
)

type CancelState struct {
	Reason string
}

func (s CancelState) ToMap() map[string]interface{} {
	if s.Reason == "" {
		return map[string]interface{}{}
	}
	return map[string]interface{}{"reason": s.Reason}
}

func (s *CancelState) FromMap(_ *model.Board, m map[string]interface{}) error {
	if m == nil {
		return nil
	}
	if v, ok := m["reason"]; ok {
		if reason, ok := v.(string); ok {
			s.Reason = reason
		}
	}
	return nil
}

type CancelModifyAction struct {
	logic.BaseAction[CancelState, logic.EffectContext]
}

func (CancelModifyAction) Tag() logic.ModifyActionTag { return action.CANCEL_MODIFY_ACTION }
func (a CancelModifyAction) LocalizeTag() pkg.LocalizeTag {
	return pkg.LocalizeTag(a.Tag())
}

func (m CancelModifyAction) Modify(state logic.ActionState, context logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if context == nil {
		return nil, nil
	}
	context.Cancel()
	summary := logic.Summary{}
	switch s := state.(type) {
	case *CancelState:
		if s != nil && s.Reason != "" {
			summary["reason"] = s.Reason
		}
	}
	return context, summary
}
