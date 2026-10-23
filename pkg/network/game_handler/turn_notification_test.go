package game_handler_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"frontage/pkg/network/game_api"
	"frontage/pkg/network/game_dispatcher"
	"frontage/pkg/network/game_handler"
)

func TestTurnNotificationHandlersKeepAllFields(t *testing.T) {
	start := game_api.TurnStartPacket{Turn: 7, IsYourSide: true, TurnTime: 60000, ExtraTime: 3000}
	data, err := json.Marshal(start)
	if err != nil {
		t.Fatal(err)
	}
	gotStart, err := game_handler.NewTurnStartHandler().ServePacket(data)
	if err != nil || !reflect.DeepEqual(gotStart, start) {
		t.Fatalf("start = %+v, %v; want %+v", gotStart, err, start)
	}
	end := game_api.TurnEndPacket{Turn: 7, IsYourSide: true}
	data, err = json.Marshal(end)
	if err != nil {
		t.Fatal(err)
	}
	gotEnd, err := game_handler.NewTurnEndHandler().ServePacket(data)
	if err != nil || !reflect.DeepEqual(gotEnd, end) {
		t.Fatalf("end = %+v, %v; want %+v", gotEnd, err, end)
	}
}

func TestTurnNotificationRoundTrip(t *testing.T) {
	for _, side := range []bool{false, true} {
		start := game_api.TurnStartPacket{Turn: 12, IsYourSide: side, TurnTime: 123456789, ExtraTime: 4567}
		packet, err := game_dispatcher.NewTurnStartDispatcher().DispatchPacket(start)
		if err != nil || packet.PacketTag() != start.PacketTag() {
			t.Fatalf("start dispatch failed: %v", err)
		}
		data, err := json.Marshal(packet)
		if err != nil {
			t.Fatal(err)
		}
		restoredStart, err := game_handler.NewTurnStartHandler().ServePacket(data)
		if err != nil || restoredStart != start {
			t.Fatalf("start round trip = %+v, %v; want %+v", restoredStart, err, start)
		}
		end := game_api.TurnEndPacket{Turn: 12, IsYourSide: side}
		packet, err = game_dispatcher.NewTurnEndDispatcher().DispatchPacket(end)
		if err != nil || packet.PacketTag() != end.PacketTag() {
			t.Fatalf("end dispatch failed: %v", err)
		}
		data, err = json.Marshal(packet)
		if err != nil {
			t.Fatal(err)
		}
		restoredEnd, err := game_handler.NewTurnEndHandler().ServePacket(data)
		if err != nil || restoredEnd != end {
			t.Fatalf("end round trip = %+v, %v; want %+v", restoredEnd, err, end)
		}
	}
	if _, err := game_handler.NewTurnStartHandler().ServePacket([]byte("{")); err == nil {
		t.Fatal("invalid start JSON was accepted")
	}
	if _, err := game_handler.NewTurnEndHandler().ServePacket([]byte("{")); err == nil {
		t.Fatal("invalid end JSON was accepted")
	}
}
