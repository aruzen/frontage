package repository

import (
	"testing"

	"frontage/pkg/network"
	"github.com/google/uuid"
)

func TestGameChannelStatusLifetimes(t *testing.T) {
	var status GameChannelStatus
	active, stopped := status.Snapshot()
	if active || status.Load() {
		t.Fatal("new channel is active")
	}
	select {
	case <-stopped:
	default:
		t.Fatal("inactive channel notification is open")
	}
	for i := 0; i < 3; i++ {
		status.Store(true)
		active, current := status.Snapshot()
		if !active || current == stopped {
			t.Fatal("new game reused the previous stop notification")
		}
		status.Store(true)
		select {
		case <-current:
			t.Fatal("active game was stopped")
		default:
		}
		status.Store(false)
		status.Store(false)
		select {
		case <-current:
		default:
			t.Fatal("game stop did not notify waiters")
		}
		stopped = current
	}
}

func TestRemovingOrReplacingGameChannelStopsPrevious(t *testing.T) {
	id := uuid.New()
	t.Cleanup(func() { RemoveGameChannel(id) })
	for _, replace := range []bool{true, false} {
		previous := AddGameChannel(id, make(chan network.UnsolvedPacket))
		previous.Living.Store(true)
		_, stopped := previous.Living.Snapshot()
		if replace {
			AddGameChannel(id, make(chan network.UnsolvedPacket))
		} else {
			RemoveGameChannel(id)
		}
		select {
		case <-stopped:
		default:
			t.Fatal("removed channel still has pending delivery waiters")
		}
	}
}
