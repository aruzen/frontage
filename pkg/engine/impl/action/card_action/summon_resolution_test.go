package card_action

import (
	"testing"

	"frontage/pkg"
	"frontage/pkg/engine/impl/card"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
)

func TestSummonUsesCardAfterEarlierEffects(t *testing.T) {
	for _, strategy := range []model.GenerationStrategy{model.GENERATION_STRATEGY_CHAIN, model.GENERATION_STRATEGY_SWAP} {
		name := "chain"
		if strategy == model.GENERATION_STRATEGY_SWAP {
			name = "swap"
		}
		t.Run(name, func(t *testing.T) {
			original := card.NewBasePiece(model.NewBaseCard("test", nil), 1, 2, 3, nil, nil)
			player := model.NewLocalPlayer(uuid.New(), nil, nil)
			player.Hand = model.Cards{original}
			board := model.NewBoard(model.NewBoardInfo(pkg.Size{Width: 3, Height: 3}, strategy),
				[2]model.Player{player, model.NewProxyPlayer(uuid.New())})
			hp := NewPieceActionState(player, model.DECK_TYPE_HAND, original, 4)
			mp := NewPieceActionState(player, model.DECK_TYPE_HAND, original, 3)
			atk := NewPieceActionState(player, model.DECK_TYPE_HAND, original, 2)
			summon := NewPieceSummonActionState(player, model.DECK_TYPE_HAND, original, pkg.Point{X: 1, Y: 1})
			es := logic.NewEventSystem(board)
			es.Emit(logic.NewEffectEvent(PieceHPIncreaseAction{}, &hp))
			es.ChainLine(logic.NewEffectEvent(PieceMPIncreaseAction{}, &mp),
				logic.NewEffectEvent(PieceATKIncreaseAction{}, &atk),
				logic.NewEffectEvent(PieceSummonAction{}, &summon))

			check := func(result *model.Board) {
				t.Helper()
				if result == nil {
					t.Fatal("no result board")
				}
				piece, ok := result.GetPiece(summon.SummonID())
				if !ok {
					t.Fatal("summoned piece is missing")
				}
				if piece.HP() != 5 || piece.MP() != 5 || piece.ATK() != 5 {
					t.Errorf("summoned HP/MP/ATK = %d/%d/%d; want 5/5/5", piece.HP(), piece.MP(), piece.ATK())
				}
				owner, _ := result.FindPlayer(player.ID())
				if len(owner.(*model.LocalPlayer).Hand) != 0 {
					t.Error("summoned card remains in hand")
				}
				if piece.Owner() != owner {
					t.Error("summoned piece belongs to an older board generation")
				}
			}
			for i := 0; i < 2; i++ {
				trial, _ := es.ResolveTest()
				check(trial)
				if original.HP() != 1 || original.MP() != 2 || original.ATK() != 3 || len(player.Hand) != 1 {
					t.Fatal("simulation changed the original card or hand")
				}
			}
			es.Resolve()
			check(es.Board)
		})
	}
}
