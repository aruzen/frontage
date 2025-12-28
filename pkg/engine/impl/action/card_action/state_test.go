package card_action

import (
	"testing"

	"frontage/pkg"
	"frontage/pkg/engine/impl/card"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
)

func newTestBoardWithCard(t *testing.T) (*model.Board, *model.LocalPlayer, card.Piece) {
	t.Helper()
	base := model.NewBaseCard("test_card", model.Materials{model.PYRO: 1})
	piece := card.NewBasePiece(base, 1, 2, 3, nil, nil)
	mainDeck := model.Cards{piece}
	player := model.NewLocalPlayer(uuid.New(), mainDeck, nil)
	players := [2]model.Player{player, model.NewProxyPlayer(uuid.New())}
	board := model.NewBoard(model.NewBoardInfo(pkg.Size{Width: 5, Height: 5}, model.GENERATION_STRATEGY_SWAP), players)
	return board, player, piece
}

func TestCardActionStatesToMapFromMap(t *testing.T) {
	// PieceActionState
	{
		board, player, piece := newTestBoardWithCard(t)
		state := PieceActionState{
			holderID: player.ID(),
			cardID:   piece.Id(),
			deckType: model.DECK_TYPE_MAIN,
			value:    5,
			holder:   player,
			piece:    piece,
		}
		m := state.ToMap()
		var dst PieceActionState
		if err := dst.FromMap(board, m); err != nil {
			t.Fatalf("piece action FromMap error: %v", err)
		}
		if dst.holderID != player.ID() || dst.cardID != piece.Id() || dst.deckType != model.DECK_TYPE_MAIN || dst.value != 5 {
			t.Fatalf("piece action state mismatch: %+v", dst)
		}
		if dst.holder == nil || dst.piece == nil {
			t.Fatalf("piece action holder/piece not restored")
		}
	}

	// PieceSummonActionState
	{
		board, player, piece := newTestBoardWithCard(t)
		summonID := uuid.New()
		point := pkg.Point{X: 2, Y: 1}
		state := PieceSummonActionState{
			holderID: player.ID(),
			cardID:   piece.Id(),
			summonID: summonID,
			deckType: model.DECK_TYPE_MAIN,
			point:    point,
			holder:   player,
			piece:    piece,
		}
		m := state.ToMap()
		var dst PieceSummonActionState
		if err := dst.FromMap(board, m); err != nil {
			t.Fatalf("piece summon FromMap error: %v", err)
		}
		if dst.holderID != player.ID() || dst.cardID != piece.Id() || dst.summonID != summonID || dst.deckType != model.DECK_TYPE_MAIN || dst.point != point {
			t.Fatalf("piece summon state mismatch: %+v", dst)
		}
		if dst.holder == nil || dst.piece == nil {
			t.Fatalf("piece summon holder/piece not restored")
		}
	}
}
