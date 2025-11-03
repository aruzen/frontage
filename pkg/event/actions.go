package event

import "frontage/pkg/event/action/card"

var (
	CARD_Piece_HP_INCREASE_ACTION = card.PieceHPIncreaseAction{}

	CARD_Piece_HP_DECREASE_ACTION = card.PieceHPDecreaseAction{}

	CARD_Piece_HP_FIX_ACTION = card.PieceHPFixAction{}

	CARD_Piece_MP_INCREASE_ACTION = card.PieceMPIncreaseAction{}

	CARD_Piece_MP_DECREASE_ACTION = card.PieceMPDecreaseAction{}

	CARD_Piece_MP_FIX_ACTION = card.PieceMPFixAction{}

	CARD_Piece_ATK_INCREASE_ACTION = card.PieceATKIncreaseAction{}

	CARD_Piece_ATK_DECREASE_ACTION = card.PieceATKDecreaseAction{}

	CARD_Piece_ATK_FIX_ACTION = card.PieceATKFixAction{}
)
