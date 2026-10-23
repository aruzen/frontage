package player_action

import (
	"testing"

	"frontage/pkg"
	"frontage/pkg/engine/impl/card"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
)

func newTestBoardWithLocalPlayer(t *testing.T) (*model.Board, *model.LocalPlayer, model.Card, model.Card, model.Card) {
	t.Helper()
	main1 := card.NewBasePiece(model.NewBaseCard("main1", model.Materials{}), 1, 1, 1, nil, nil)
	main2 := card.NewBasePiece(model.NewBaseCard("main2", model.Materials{}), 1, 1, 1, nil, nil)
	sub1 := card.NewBasePiece(model.NewBaseCard("sub1", model.Materials{}), 1, 1, 1, nil, nil)

	local := model.NewLocalPlayer(uuid.New(), model.Cards{main1, main2}, model.Cards{sub1})
	local.Hand = model.Cards{}
	other := model.NewProxyPlayer(uuid.New())
	board := model.NewBoard(model.NewBoardInfo(pkg.Size{Width: 3, Height: 3}, model.GENERATION_STRATEGY_SWAP), [2]model.Player{local, other})
	return board, local, main1, main2, sub1
}

func TestDrawActionStateToMapFromMap(t *testing.T) {
	board, local, _, _, _ := newTestBoardWithLocalPlayer(t)
	state := NewDrawActionState(local)
	m := state.ToMap()

	var dst DrawActionState
	if err := dst.FromMap(board, m); err != nil {
		t.Fatalf("FromMap error: %v", err)
	}
	if dst.playerID != local.ID() {
		t.Fatalf("player id mismatch: got %v want %v", dst.playerID, local.ID())
	}
	if dst.player == nil || dst.player.ID() != local.ID() {
		t.Fatalf("player not restored")
	}
}

func TestDrawMainDeckActionSolveMovesCard(t *testing.T) {
	board, local, _, main2, _ := newTestBoardWithLocalPlayer(t)
	action := DrawMainDeckAction{}
	state := NewDrawActionState(local)

	next, summary := action.Solve(board, state, &DrawActionContext{})
	if next == nil {
		t.Fatalf("expected next board")
	}
	player, ok := next.FindPlayer(local.ID())
	if !ok {
		t.Fatalf("player not found")
	}
	localNext, ok := player.(*model.LocalPlayer)
	if !ok {
		t.Fatalf("player is not LocalPlayer")
	}
	if len(localNext.MainDeck) != 1 {
		t.Fatalf("main deck size: got %d want %d", len(localNext.MainDeck), 1)
	}
	if len(localNext.Hand) != 1 {
		t.Fatalf("hand size: got %d want %d", len(localNext.Hand), 1)
	}
	if localNext.Hand[0].Id() != main2.Id() {
		t.Fatalf("unexpected drawn card: got %v want %v", localNext.Hand[0].Id(), main2.Id())
	}
	if summary["card_id"] != main2.Id().String() {
		t.Fatalf("summary card_id mismatch")
	}
}

func TestDrawSubDeckActionSolveMovesCard(t *testing.T) {
	board, local, _, _, sub1 := newTestBoardWithLocalPlayer(t)
	action := DrawSubDeckAction{}
	state := NewDrawActionState(local)

	next, summary := action.Solve(board, state, &DrawActionContext{})
	if next == nil {
		t.Fatalf("expected next board")
	}
	player, ok := next.FindPlayer(local.ID())
	if !ok {
		t.Fatalf("player not found")
	}
	localNext, ok := player.(*model.LocalPlayer)
	if !ok {
		t.Fatalf("player is not LocalPlayer")
	}
	if len(localNext.SubDeck) != 0 {
		t.Fatalf("sub deck size: got %d want %d", len(localNext.SubDeck), 0)
	}
	if len(localNext.Hand) != 1 {
		t.Fatalf("hand size: got %d want %d", len(localNext.Hand), 1)
	}
	if localNext.Hand[0].Id() != sub1.Id() {
		t.Fatalf("unexpected drawn card: got %v want %v", localNext.Hand[0].Id(), sub1.Id())
	}
	if summary["card_id"] != sub1.Id().String() {
		t.Fatalf("summary card_id mismatch")
	}
}

func TestDrawActionEmptyDeckNoChange(t *testing.T) {
	main := card.NewBasePiece(model.NewBaseCard("main", model.Materials{}), 1, 1, 1, nil, nil)
	local := model.NewLocalPlayer(uuid.New(), model.Cards{}, model.Cards{})
	local.Hand = model.Cards{main}
	other := model.NewProxyPlayer(uuid.New())
	board := model.NewBoard(model.NewBoardInfo(pkg.Size{Width: 3, Height: 3}, model.GENERATION_STRATEGY_SWAP), [2]model.Player{local, other})

	action := DrawMainDeckAction{}
	state := NewDrawActionState(local)
	next, summary := action.Solve(board, state, &DrawActionContext{})
	if next != board {
		t.Fatalf("expected original board")
	}
	if summary != nil {
		t.Fatalf("expected nil summary")
	}
	if len(local.Hand) != 1 {
		t.Fatalf("hand size changed: %d", len(local.Hand))
	}
}
