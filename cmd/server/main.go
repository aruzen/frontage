package main

import (
	"context"
	"flag"
	"frontage/internal/log"
	"frontage/pkg/engine/impl/action"
	action_register "frontage/pkg/engine/impl/action/register"
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
	"os"
	"os/signal"
	"syscall"
)

type entryAndExitInfo struct {
	id      uuid.UUID
	isEntry bool
	channel chan network.UnsolvedPacket
	ctx     context.Context
}

var (
	matchRepo  = repository.NewMatchRepository()
	cardRepo   = repository.NewCardRepository()
	actionRepo *repository.ActionRepository
)

type ServerInfo struct {
	Addr        string
	TLSEnabled  bool
	TLSCertFile string
	TLSKeyFile  string
}

func main() {
	var info ServerInfo
	flag.StringVar(&info.Addr, "addr", ":8275", "TCP listen address")
	flag.StringVar(&info.TLSCertFile, "tls-cert", "", "TLS certificate file (requires -tls-key)")
	flag.StringVar(&info.TLSKeyFile, "tls-key", "", "TLS private key file (requires -tls-cert)")
	flag.Parse()
	info.TLSEnabled = info.TLSCertFile != "" || info.TLSKeyFile != ""
	log.Init(true)
	action_register.Init()
	actionRepo = repository.NewActionRepository(action.FindActionModify, action.FindActionEffect)
	listener, err := newListener(info)
	if err != nil {
		slog.Error("Failed to start server", "err", err)
		return
	}
	defer listener.Close()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := serve(ctx, listener); err != nil {
		slog.Error("Server stopped", "err", err)
	}
}

func systemLoop(ctx context.Context, visitChan chan entryAndExitInfo) {
	handlers := controller.SystemPacketHandlers{}
	packetLoop(ctx, visitChan, func(_ context.Context, id uuid.UUID, packet network.UnsolvedPacket) error {
		return controller.DispatchSystemPacket(handlers, packet.Tag, id, packet.Body)
	}, nil)
}

func lobbyLoop(ctx context.Context, visitChan chan entryAndExitInfo) {
	matches := matchRepo
	handlers := controller.LobbyPacketHandlers{
		MatchMake: lobby_handler.NewMatchMakeHandler(
			lobby_service.NewMatchMakeService(matches,
				pvp.RequireContents{CardRepo: cardRepo},
				pve.RequireContents{
					ActionRepo: actionRepo,
					CardRepo:   cardRepo,
					ActEventDisp: game_dispatcher.NewActEventDispatcher(
						translator.NewActionResultTranslator(actionRepo),
						translator.NewActionSummaryTranslator(actionRepo)),
					GameInitDisp: game_dispatcher.NewGameInitializeDispatcher(),
				}), matches),
	}
	packetLoop(ctx, visitChan, func(playerCtx context.Context, id uuid.UUID, packet network.UnsolvedPacket) error {
		return controller.DispatchLobbyPacketContext(playerCtx, handlers, packet.Tag, id, packet.Body)
	}, func(id uuid.UUID) {
		if match, ok := matches.FindByPlayerID(id); ok {
			matches.Remove(match)
		}
	})
}
