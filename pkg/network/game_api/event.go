package game_api

import "frontage/pkg/network/data"

type ActEventPayload struct {
	Result  data.ActionResult
	Summary []data.ActionSummary
}

type ActEventPacket struct {
	Events []ActEventPayload
}
