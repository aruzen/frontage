package piece_action

import "frontage/pkg/engine/logic"

// EnumerateEffectAction returns all EffectActions defined in this package.
func EnumerateEffectAction() []logic.EffectAction {
	return []logic.EffectAction{
		PieceSummonAction{},
		PieceMoveAction{},
		PieceAttackAction{},
		PieceInvasionAction{},
		PieceHPIncreaseAction{},
		PieceHPDecreaseAction{},
		PieceHPFixAction{},
		PieceMPIncreaseAction{},
		PieceMPDecreaseAction{},
		PieceMPFixAction{},
		PieceATKIncreaseAction{},
		PieceATKDecreaseAction{},
		PieceATKFixAction{},
	}
}
