package game_api

import "frontage/pkg/network/data"

type ActEventPayload struct {
	Result  data.ActionResult
	Summary []data.EventSummary
}

type ActEventPacket struct {
	Events []ActEventPayload
}
