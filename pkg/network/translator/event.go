package translator

import (
	"frontage/pkg/engine/logic"
	"frontage/pkg/network/data"
	"frontage/pkg/network/repository"
	"reflect"
)

type ActionResultTranslator struct {
	effectRepo repository.EffectRepository
}

type EventSummaryTranslator struct {
}

func NewActionResultTranslator(repo repository.EffectRepository) *ActionResultTranslator {
	return &ActionResultTranslator{
		effectRepo: repo,
	}
}

func (t *ActionResultTranslator) ToModel(d data.ActionResult) (logic.ActionResult, error) {
	action, err := t.effectRepo.Find(logic.EffectActionTag(d.ActionTag))
	if err != nil {
		return logic.ActionResult{}, err
	}
	context, ok := reflect.New(action.WantContext()).Interface().(logic.EffectContext)
	if !ok {
		return logic.ActionResult{}, ErrNewContextFailed
	}
	err = context.FromMap(d.Data)
	if err != nil {
		return logic.ActionResult{}, err
	}
	return logic.ActionResult{Action: action, Context: context}, err
}
