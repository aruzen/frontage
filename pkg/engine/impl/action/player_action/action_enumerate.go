package player_action

import "frontage/pkg/engine/logic"

// EnumerateEffectAction returns all EffectActions defined in this package.
func EnumerateEffectAction() []logic.EffectAction {
	return []logic.EffectAction{
		DrawMainDeckAction{},
		DrawSubDeckAction{},
	}
}
