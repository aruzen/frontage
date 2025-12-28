package lobby_service

import (
	"context"
	"errors"
	"frontage/pkg/network/controller/pve"
	"frontage/pkg/network/controller/pvp"
	"frontage/pkg/network/data"
	"frontage/pkg/network/lobby_api"
	"frontage/pkg/network/repository"
	"github.com/google/uuid"
)

type MatchMakeService struct {
	matchRepo *repository.MatchRepository
	pvpRepos  pvp.RequireRepositories
	pveRepos  pve.RequireRepositories
}

func NewMatchMakeService(matchRepo *repository.MatchRepository, pvp pvp.RequireRepositories, pve pve.RequireRepositories) *MatchMakeService {
	return &MatchMakeService{matchRepo, pvp, pve}
}

func (m MatchMakeService) MatchMake(ctx context.Context, id uuid.UUID, matchType data.MatchType) error {
	switch matchType {
	case data.PvE:
		matchID := uuid.New()
		if m.matchRepo != nil {
			m.matchRepo.Insert(matchID, id, uuid.Nil)
		}
		go pve.Game(ctx, m.pveRepos, id, pve.DefaultGameInfo())
		repository.SendPacket(id, lobby_api.CompleteMatchMakePacket{MatchID: matchID})
		return nil
	case data.PvP:
		return errors.New("pvp match not implemented")
	}
	return errors.New("unknown match type")
}
