package lobby_api

import (
	"context"
	"frontage/pkg/network"
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
		go pve.Game(m.pveRepos, id, pve.DefaultGameInfo())
		network.SendPacket(id, lobby_api.CompleteMatchMakePacket{})
	case data.PvP:
		m.matchRepo.

		go pve.Game(m.pveRepos, id, pve.DefaultGameInfo())
		network.SendPacket(id, lobby_api.CompleteMatchMakePacket{})
	}
}
