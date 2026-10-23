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
	pvpRepos  pvp.RequireContents
	pveRepos  pve.RequireContents
}

func NewMatchMakeService(matchRepo *repository.MatchRepository, pvp pvp.RequireContents, pve pve.RequireContents) *MatchMakeService {
	return &MatchMakeService{matchRepo, pvp, pve}
}

func (m MatchMakeService) MatchMake(ctx context.Context, id uuid.UUID, matchType data.MatchType) error {
	switch matchType {
	case data.PvE:
		channel := repository.GetGameChannel(id)
		if channel == nil || channel.Chan == nil {
			return errors.New("game channel not found")
		}
		matchID := uuid.New()
		if m.matchRepo != nil {
			m.matchRepo.Insert(matchID, id, uuid.Nil)
		}
		channel.Living.Store(true)
		if !repository.SendPacket(id, lobby_api.CompleteMatchMakePacket{MatchID: matchID}) {
			channel.Living.Store(false)
			if m.matchRepo != nil {
				m.matchRepo.Remove(matchID)
			}
			return errors.New("failed to send match confirmation")
		}
		go func() {
			defer channel.Living.Store(false)
			if m.matchRepo != nil {
				defer m.matchRepo.Remove(matchID)
			}
			pve.Game(ctx, m.pveRepos, id, pve.DefaultGameInfo())
		}()
		return nil
	case data.PvP:
		return errors.New("pvp match not implemented")
	}
	return errors.New("unknown match type")
}
