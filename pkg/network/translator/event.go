package translator

import (
	"frontage/pkg/engine/logic"
	"frontage/pkg/network/data"
	"frontage/pkg/network/repository"
	"reflect"
)

type ActionResultTranslator struct {
	effectRepo repository.ActionRepository
}

type EventSummaryTranslator struct {
	effectRepo repository.ActionRepository
}

func NewActionResultTranslator(repo repository.ActionRepository) *ActionResultTranslator {
	return &ActionResultTranslator{
		effectRepo: repo,
	}
}

func (t *ActionResultTranslator) ToModel(d data.ActionResult) (logic.ActionResult, error) {
	action, err := t.effectRepo.FindEffect(logic.EffectActionTag(d.ActionTag))
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

func (t *ActionResultTranslator) FromModel(m logic.ActionResult) (data.ActionResult, error) {
	return data.ActionResult{
		ActionTag: string(m.Action.LocalizeTag()),
		Data:      m.Context.ToMap(),
	}, nil
}

func NewEventSummaryTranslator(effectRepo repository.ActionRepository) *EventSummaryTranslator {
	return &EventSummaryTranslator{
		effectRepo: effectRepo,
	}
}

/* TODO
func (t *EventSummaryTranslator) ToModel(d data.EventSummary) (logic.Action, map[string]interface{}, error) {
	panic("implement me")
}

func (t *EventSummaryTranslator) FromModel(d data.EventSummary) (logic.Action, map[string]interface{}, error) {
	panic("implement me")
}
*/
