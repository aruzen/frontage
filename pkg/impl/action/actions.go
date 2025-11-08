package action

import (
	"frontage/pkg/impl/action/card_action"
)

var (
	CARD_Piece_HP_INCREASE_ACTION = card_action.PieceHPIncreaseAction{}

	CARD_Piece_HP_DECREASE_ACTION = card_action.PieceHPDecreaseAction{}

	CARD_Piece_HP_FIX_ACTION = card_action.PieceHPFixAction{}

	CARD_Piece_MP_INCREASE_ACTION = card_action.PieceMPIncreaseAction{}

	CARD_Piece_MP_DECREASE_ACTION = card_action.PieceMPDecreaseAction{}

	CARD_Piece_MP_FIX_ACTION = card_action.PieceMPFixAction{}

	CARD_Piece_ATK_INCREASE_ACTION = card_action.PieceATKIncreaseAction{}

	CARD_Piece_ATK_DECREASE_ACTION = card_action.PieceATKDecreaseAction{}

	CARD_Piece_ATK_FIX_ACTION = card_action.PieceATKFixAction{}
)
