package common

import "frontage/pkg/engine/logic"

// EnumerateModifyAction returns all ModifyActions defined in this package.
func EnumerateModifyAction() []logic.ModifyAction {
	return []logic.ModifyAction{
		CancelModifyAction{},
	}
}
