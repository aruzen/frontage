package piece_action

import (
	"frontage/pkg"
	"frontage/pkg/engine/impl/action"
	"frontage/pkg/engine/logic"
)

type MoveCancelState struct {
	Reason string
}

type MoveCancelModifyAction struct {
	logic.BaseAction[MoveCancelState, PieceMoveActionContext]
}

func (MoveCancelModifyAction) Tag() logic.ModifyActionTag { return action.MOVE_CANCEL_MODIFY_ACTION }
func (a MoveCancelModifyAction) LocalizeTag() pkg.LocalizeTag {
	return pkg.LocalizeTag(a.Tag())
}

func (m MoveCancelModifyAction) Modify(state interface{}, context logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if context == nil {
		return nil, nil
	}
	context.Cancel()
	if s, ok := state.(MoveCancelState); ok && s.Reason != "" {
		return context, logic.Summary{"reason": s.Reason}
	}
	return context, nil
}
