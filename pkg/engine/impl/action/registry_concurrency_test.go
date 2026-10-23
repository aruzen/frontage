package action_test

import (
	"fmt"
	"sync"
	"testing"

	"frontage/pkg/engine/impl/action"
	"frontage/pkg/engine/impl/action/common"
	"frontage/pkg/engine/logic"
	"github.com/google/uuid"
)

func TestConcurrentActionRegistrationAndLookup(t *testing.T) {
	prefix := "test/" + uuid.NewString()
	start := make(chan struct{})
	var workers sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		workers.Add(1)
		go func(worker int) {
			defer workers.Done()
			<-start
			for i := 0; i < 32; i++ {
				tag := fmt.Sprintf("%s/%d/%d", prefix, worker, i)
				effectTag, modifyTag := logic.EffectActionTag(tag), logic.ModifyActionTag(tag)
				if err := action.Register(effectTag, logic.GameStartAction{}); err != nil {
					t.Error(err)
				}
				if err := action.RegisterModify(modifyTag, common.CancelModifyAction{}); err != nil {
					t.Error(err)
				}
				if action.FindActionEffect(effectTag) == nil || action.FindActionModify(modifyTag) == nil {
					t.Errorf("registered action %s is missing", tag)
				}
			}
		}(worker)
	}
	close(start)
	workers.Wait()
}
