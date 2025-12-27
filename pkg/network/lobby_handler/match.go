package lobby_handler

import (
	"context"
	"encoding/json"
	"frontage/pkg/network/data"
	"frontage/pkg/network/lobby_api"
	"frontage/pkg/network/repository"
	"github.com/google/uuid"
	"log/slog"
)

type MatchMakeService interface {
	MatchMake(ctx context.Context, id uuid.UUID, matchType data.MatchType) error
}

type MatchMakeHandler struct {
	service   MatchMakeService
	matchRepo repository.MatchRepository
}

func NewLoginHandler(service MatchMakeService, matchRepo repository.MatchRepository) *MatchMakeHandler {
	return &MatchMakeHandler{service: service, matchRepo: matchRepo}
}

func (h *MatchMakeHandler) ServePacket(id uuid.UUID, data []byte) error {
	_, found := h.matchRepo.FindByPlayerID(id)
	if found {
		// TODO 失敗パケットを送るべき
		slog.Error("Already login", "playerID", id)
		return nil
	}
	var packet lobby_api.WaitMatchMakePacket
	if err := json.Unmarshal(data, &packet); err != nil {
		return err
	}
	return h.service.MatchMake(context.Background(), id, packet.Type)
}
