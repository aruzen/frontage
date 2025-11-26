package game_api

type ActEventPayload struct {
	ActionTag string
	data      map[string]interface{}
}

type ActEventPacket struct {
	events []ActEventPayload
}
