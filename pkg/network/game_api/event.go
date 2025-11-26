package game_api

import "frontage/pkg/network/data"

type ActEventPayload struct {
	result  data.ActionResult
	summary []data.EventSummary
}

type ActEventPacket struct {
	events []ActEventPayload
}
