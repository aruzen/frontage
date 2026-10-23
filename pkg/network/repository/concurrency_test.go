package repository

import (
	"fmt"
	"sync"
	"testing"

	"frontage/pkg"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
)

func TestRepositoriesSupportConcurrentConnections(t *testing.T) {
	cards := NewCardRepository()
	matches := NewMatchRepository()
	var workers sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		workers.Add(1)
		go func(worker int) {
			defer workers.Done()
			for i := 0; i < 50; i++ {
				tag := pkg.LocalizeTag(fmt.Sprintf("%d/%d", worker, i))
				card := model.NewBaseCard(tag, nil)
				if err := cards.Insert(card); err != nil {
					t.Error(err)
					return
				}
				if found, err := cards.Find(tag); err != nil || found.Id() != card.Id() {
					t.Errorf("card lookup failed: %v", err)
					return
				}
				match, first, second := uuid.New(), uuid.New(), uuid.New()
				matches.Insert(match, first, second)
				if found, ok := matches.FindByPlayerID(first); !ok || found != match {
					t.Error("match lookup failed")
					return
				}
				matches.Remove(match)
				if _, ok := matches.FindByPlayerID(second); ok {
					t.Error("removed match retained a player")
					return
				}
			}
		}(worker)
	}
	workers.Wait()
	if len(matches.matchIds) != 0 || len(matches.matchTable) != 0 {
		t.Fatal("completed work retained matches")
	}
}
