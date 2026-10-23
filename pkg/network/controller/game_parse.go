package controller

import (
	"frontage/pkg/network/game_handler"
	"frontage/pkg/network/repository"
)

type GamePacketParsers struct {
	MyDeckUpload      *game_handler.MyDeckUploadHandler
	TurnPass          *game_handler.TurnPassHandler
	DrawPhase         *game_handler.DrawPhaseHandler
	Sync              *game_handler.SyncHandler
	ActPieceSkill     *game_handler.ActPieceSkillHandler
	ActStructureSkill *game_handler.ActStructureSkillHandler
	PlayCard          *game_handler.PlayCardHandler
}

func MakeDefaultParsers(cardRepo *repository.CardRepository, actionRepo *repository.ActionRepository) GamePacketParsers {
	return GamePacketParsers{
		MyDeckUpload:      game_handler.NewMyDeckUploadHandler(cardRepo),
		TurnPass:          game_handler.NewTurnPassHandler(),
		DrawPhase:         game_handler.NewDrawPhaseHandler(),
		Sync:              game_handler.NewSyncHandler(),
		ActPieceSkill:     game_handler.NewActPieceSkillHandler(),
		ActStructureSkill: game_handler.NewActStructureSkillHandler(),
		PlayCard:          game_handler.NewPlayCardHandler(),
	}
}
