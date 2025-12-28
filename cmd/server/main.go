package main

import (
	"context"
	"fmt"
	"frontage/internal/log"
	"frontage/pkg/engine/impl/action"
	action_register "frontage/pkg/engine/impl/action/register"
	"frontage/pkg/engine/logic"
	"frontage/pkg/network"
	"frontage/pkg/network/controller"
	"frontage/pkg/network/controller/pve"
	"frontage/pkg/network/controller/pvp"
	"frontage/pkg/network/game_dispatcher"
	"frontage/pkg/network/lobby_handler"
	"frontage/pkg/network/lobby_service"
	"frontage/pkg/network/repository"
	"frontage/pkg/network/translator"
	"github.com/google/uuid"
	"log/slog"
	"net"
	"reflect"
)

type entryAndExitInfo struct {
	id      uuid.UUID
	isEntry bool
	channel chan network.UnsolvedPacket
}

var (
	matchRepo  = repository.NewMatchRepository()
	cardRepo   = repository.NewCardRepository()
	actionRepo *repository.ActionRepository
)

func main() {
	action_register.Init()
	actionRepo = repository.NewActionRepository(func(tag logic.ModifyActionTag) logic.ModifyAction {
		return action.FindActionModify(tag)
	}, func(tag logic.EffectActionTag) logic.EffectAction {
		return action.FindActionEffect(tag)
	})
	log.Init(true)

	systemCtx := context.Background()
	systemCtx, systemFinish := context.WithCancel(systemCtx)
	systemVisitPlayer := make(chan entryAndExitInfo)
	lobbyVisitPlayer := make(chan entryAndExitInfo)
	go systemLoop(systemCtx, systemVisitPlayer)
	go lobbyLoop(systemCtx, lobbyVisitPlayer)

	addr, err := net.ResolveTCPAddr("udp", ":8275")
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		slog.Error("Failed to listen TCP address", "err", err)
		return
	}
	if err != nil {
		slog.Error("Failed to resolve TCP address", "err", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		go receive(systemCtx, conn, systemVisitPlayer, lobbyVisitPlayer)
	}
	systemFinish()
}

func receive(ctx context.Context, conn net.Conn, systemVisitPlayer, lobbyVisitPlayer chan entryAndExitInfo) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			slog.Error("Failed to close TCP connection", "err", err)
		}
	}(conn)
	id, err := uuid.NewUUID()
	if err != nil {
		slog.Error("Failed to generate UUID", "err", err)
	}
	repository.AddConnection(id, conn)
	defer repository.RemoveConnection(id)

	systemChan := make(chan network.UnsolvedPacket)
	lobbyChan := make(chan network.UnsolvedPacket)
	gameChan := make(chan network.UnsolvedPacket)
	systemVisitPlayer <- entryAndExitInfo{
		id:      id,
		isEntry: true,
		channel: systemChan,
	}
	lobbyVisitPlayer <- entryAndExitInfo{
		id:      id,
		isEntry: true,
		channel: lobbyChan,
	}
	barrierGameChan := repository.AddGameChannel(id, gameChan)
	defer func() {
		systemVisitPlayer <- entryAndExitInfo{
			id:      id,
			isEntry: false,
			channel: nil,
		}
		lobbyVisitPlayer <- entryAndExitInfo{
			id:      id,
			isEntry: false,
			channel: nil,
		}
		repository.RemoveGameChannel(id)
	}()
	err = controller.ReceiveLoop(ctx, conn, systemChan, lobbyChan, barrierGameChan)
	if err != nil {
		slog.Error("Failed to receive loop packet", "err", err)
		return
	}
}

func systemLoop(ctx context.Context, visitChan chan entryAndExitInfo) {
	handlers := controller.SystemPacketHandlers{}

	cases := make([]reflect.SelectCase, 2)
	ids := make([]uuid.UUID, 2)
	cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())}
	ids[0] = uuid.Nil
	cases[1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(visitChan)}
	ids[1] = uuid.Nil

system_finish:
	for {
		idx, val, ok := reflect.Select(cases)
		if !ok {
			continue
		}
		switch idx {
		case 0:
			break system_finish
		case 1:
			visit := val.Interface().(entryAndExitInfo)
			if !visit.isEntry {
				cases = append(cases[:idx], cases[idx+1:]...)
				ids = append(ids[:idx], ids[idx+1:]...)
			} else {
				if visit.channel == nil {
					slog.Warn("Visit channel is nil", "id", visit.id)
					continue
				}
				cases = append(cases, reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(visit.channel),
				})
				ids = append(ids, visit.id)
			}
		default:
			packet := val.Interface().(network.UnsolvedPacket)
			id := ids[idx]
			if id == uuid.Nil {
				slog.Error("Invalid packet id", "id", id)
				continue
			}
			err := controller.DispatchSystemPacket(handlers, packet.Tag, id, packet.Body)
			if err != nil {
				slog.Error("Failed to dispatch system packet", "err", err)
				return
			}
		}
		fmt.Println("index:", idx, "value:", val.Interface())
	}
}

func lobbyLoop(ctx context.Context, visitChan chan entryAndExitInfo) {
	handlers := controller.LobbyPacketHandlers{
		MatchMake: lobby_handler.NewMatchMakeHandler(
			lobby_service.NewMatchMakeService(
				matchRepo,
				pvp.RequireContents{
					cardRepo,
				},
				pve.RequireContents{
					actionRepo,
					cardRepo,
					game_dispatcher.NewActEventDispatcher(
						translator.NewActionResultTranslator(actionRepo),
						translator.NewActionSummaryTranslator(actionRepo)),
					game_dispatcher.NewGameInitializeDispatcher(),
				}),
			matchRepo),
	}

	cases := make([]reflect.SelectCase, 2)
	ids := make([]uuid.UUID, 2)
	cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())}
	ids[0] = uuid.Nil
	cases[1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(visitChan)}
	ids[1] = uuid.Nil

system_finish:
	for {
		idx, val, ok := reflect.Select(cases)
		if !ok {
			continue
		}
		switch idx {
		case 0:
			break system_finish
		case 1:
			visit := val.Interface().(entryAndExitInfo)
			if !visit.isEntry {
				cases = append(cases[:idx], cases[idx+1:]...)
				ids = append(ids[:idx], ids[idx+1:]...)
			} else {
				if visit.channel == nil {
					slog.Warn("Visit channel is nil", "id", visit.id)
					continue
				}
				cases = append(cases, reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(visit.channel),
				})
				ids = append(ids, visit.id)
			}
		default:
			packet := val.Interface().(network.UnsolvedPacket)
			id := ids[idx]
			if id == uuid.Nil {
				slog.Error("Invalid packet id", "id", id)
				continue
			}
			err := controller.DispatchLobbyPacket(handlers, packet.Tag, id, packet.Body)
			if err != nil {
				slog.Error("Failed to dispatch system packet", "err", err)
				return
			}
		}
		fmt.Println("index:", idx, "value:", val.Interface())
	}
}
