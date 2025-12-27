package game_dispatcher

import (
	"encoding/json"
	"frontage/pkg/network/data"
	"frontage/pkg/network/game_api"
)

// ActEventDispatcher builds packets for action events.
type ActEventDispatcher struct{}

func NewActEventDispatcher() *ActEventDispatcher { return &ActEventDispatcher{} }

// DispatchPacket serializes the given payloads into a JSON packet.
func (d *ActEventDispatcher) DispatchPacket(events []game_api.ActEventPayload) ([]byte, error) {
	packet := game_api.ActEventPacket{Events: events}
	return json.Marshal(packet)
}

// Helper to build payload from result/summary slices if呼び出し側が使いたい場合。
func (d *ActEventDispatcher) BuildPayload(result data.ActionResult, summary []data.ActionSummary) game_api.ActEventPayload {
	return game_api.ActEventPayload{Result: result, Summary: summary}
}
