package pve

import (
	"context"
	"frontage/pkg"
	"frontage/pkg/engine/impl/action"
	"frontage/pkg/engine/impl/action/player_action"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"frontage/pkg/network"
	"frontage/pkg/network/controller"
	"frontage/pkg/network/game_api"
	"frontage/pkg/network/repository"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type gameSession struct {
	ctx       context.Context
	rc        RequireContents
	id        uuid.UUID
	info      GameInfo
	input     <-chan network.UnsolvedPacket
	parsers   controller.GamePacketParsers
	es        *logic.EventSystem
	extraTime time.Duration
}

func Game(ctx context.Context, rc RequireContents, id uuid.UUID, info GameInfo) {
	input := repository.GetGameChannel(id)
	if input == nil || input.Chan == nil || ctx.Err() != nil {
		return
	}
	input.Living.Store(true)
	defer input.Living.Store(false)
	InitRequireRepositories(&rc)
	mainDeck, subDeck, npcMainDeck, npcSubDeck := DemoDecks(&rc)
	game := &gameSession{
		ctx: ctx, rc: rc, id: id, info: info, input: input.Chan,
		parsers: controller.MakeDefaultParsers(rc.CardRepo, rc.ActionRepo),
		es: logic.NewEventSystem(model.NewBoard(
			model.NewBoardInfo(pkg.Size{Width: 7, Height: 7}, model.GENERATION_STRATEGY_SWAP),
			[2]model.Player{
				model.NewLocalPlayer(id, mainDeck, subDeck),
				model.NewLocalPlayer(uuid.New(), npcMainDeck, npcSubDeck),
			})),
		extraTime: info.ExtraTimeLimit,
	}
	game.run()
}

func (g *gameSession) send(packet network.Packet) bool {
	return g.ctx.Err() == nil && repository.SendPacket(g.id, packet)
}

func (g *gameSession) doAction(tag logic.EffectActionTag, state logic.ActionState) bool {
	if g.ctx.Err() != nil {
		return false
	}
	effect, err := g.rc.ActionRepo.FindEffect(tag)
	if err != nil {
		slog.Error("Action not found", "tag", tag, "err", err)
		return false
	}
	g.es.Clean()
	g.es.Emit(logic.NewEffectEvent(effect, state))
	g.es.Resolve()
	packet, err := g.rc.ActEventDisp.DispatchPacket(g.es.AppliedEffects, g.es.Summaries)
	if err != nil {
		slog.Error("Failed to dispatch action", "err", err)
		return false
	}
	return g.send(packet)
}

func (g *gameSession) syncClient(cause int) bool {
	for attempt := 0; attempt < g.info.SyncResendMax; attempt++ {
		if !g.send(game_api.SyncPacket{SyncCause: cause}) {
			return false
		}
		timer := time.NewTimer(g.info.SyncResendDuration)
		acknowledged, alive := g.waitForSync(cause, timer.C)
		timer.Stop()
		if acknowledged || !alive {
			return acknowledged
		}
	}
	// A zero retry count explicitly disables synchronization (e.g. local tests).
	return g.info.SyncResendMax <= 0 && g.ctx.Err() == nil
}

func (g *gameSession) waitForSync(cause int, timeout <-chan time.Time) (bool, bool) {
	for {
		select {
		case packet, ok := <-g.input:
			if !ok {
				return false, false
			}
			if packet.Tag != network.SYNC_PACKET_TAG {
				continue
			}
			ack, err := g.parsers.Sync.ServePacket(packet.Body)
			if err == nil && ack == cause {
				return true, true
			}
		case <-timeout:
			return false, true
		case <-g.ctx.Done():
			return false, false
		}
	}
}

func (g *gameSession) run() {
	if !g.syncClient(game_api.SYNC_CAUSE_GAME_START) {
		return
	}
	packet, err := g.rc.GameInitDisp.DispatchPacket(g.es.Board, g.es.Board.Players()[0])
	if err != nil || !g.send(packet) || !g.syncClient(game_api.SYNC_CAUSE_INITIALIZE_END) {
		return
	}
	if !g.doAction(logic.GAME_START_ACTION_TAG, &logic.SystemNoticeState{}) {
		return
	}
	for g.ctx.Err() == nil {
		yourSide := g.es.Board.ActivePlayer().ID() == g.id
		if !g.startTurn(yourSide) {
			return
		}
		if yourSide && !g.playerTurn() {
			return
		}
		// Until NPC actions are implemented, it draws and passes immediately.
		turn := g.es.Board.Turn()
		if !g.doAction(logic.TURN_END_ACTION_TAG, &logic.SystemNoticeState{}) ||
			!g.send(game_api.TurnEndPacket{Turn: turn, IsYourSide: yourSide}) ||
			!g.syncClient(game_api.SYNC_CAUSE_TURN_END) {
			return
		}
	}
}

func (g *gameSession) startTurn(yourSide bool) bool {
	if !g.doAction(logic.TURN_START_ACTION_TAG, &logic.SystemNoticeState{}) ||
		!g.send(game_api.TurnStartPacket{
			Turn: g.es.Board.Turn(), IsYourSide: yourSide,
			TurnTime: int(g.info.TurnTimeLimit), ExtraTime: int(g.extraTime),
		}) {
		return false
	}
	if yourSide {
		if !g.drawPhase() {
			return false
		}
	} else if !g.doAction(action.PLAYER_DRAW_MAIN_DECK_ACTION, player_action.NewDrawActionState(g.es.Board.ActivePlayer())) {
		return false
	}
	return g.syncClient(game_api.SYNC_CAUSE_DRAW_PHASE_END)
}

func (g *gameSession) drawPhase() bool {
	timer := time.NewTimer(g.info.DrawPhaseLimit)
	defer timer.Stop()
	for {
		select {
		case packet, ok := <-g.input:
			if !ok {
				return false
			}
			if packet.Tag != network.DRAW_PHASE_PACKET_TAG {
				continue
			}
			drawType, err := g.parsers.DrawPhase.ServePacket(packet.Body)
			if err != nil {
				continue
			}
			var tag logic.EffectActionTag
			switch drawType {
			case game_api.DRAW_TYPE_MAIN_ONLY:
				tag = action.PLAYER_DRAW_MAIN_DECK_ACTION
			case game_api.DRAW_TYPE_SUB_ONLY:
				tag = action.PLAYER_DRAW_SUB_DECK_ACTION
			default:
				continue
			}
			return g.doAction(tag, player_action.NewDrawActionState(g.es.Board.ActivePlayer()))
		case <-timer.C:
			return true
		case <-g.ctx.Done():
			return false
		}
	}
}

func (g *gameSession) playerTurn() bool {
	start := time.Now()
	timer := time.NewTimer(g.info.TurnTimeLimit + g.extraTime)
	defer timer.Stop()
	for {
		select {
		case packet, ok := <-g.input:
			if !ok {
				return false
			}
			if packet.Tag != network.TURN_PASS_PACKET_TAG {
				continue
			}
			turn, err := g.parsers.TurnPass.ServePacket(packet.Body)
			if err != nil || turn != g.es.Board.Turn() {
				continue
			}
			if overtime := time.Since(start) - g.info.TurnTimeLimit; overtime > 0 {
				g.extraTime = max(0, g.extraTime-overtime)
			}
			return true
		case <-timer.C:
			g.extraTime = 0
			return true
		case <-g.ctx.Done():
			return false
		}
	}
}
