package translator

import (
	"frontage/pkg"
	"frontage/pkg/engine/logic"
	"frontage/pkg/network/data"
	"frontage/pkg/network/repository"
	"reflect"
)

type ActionResultTranslator struct {
	effectRepo repository.ActionRepository
}

type ActionSummaryTranslator struct {
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

func NewActionSummaryTranslator(effectRepo repository.ActionRepository) *ActionSummaryTranslator {
	return &ActionSummaryTranslator{
		effectRepo: effectRepo,
	}
}

func (t *ActionSummaryTranslator) ToModel(d data.ActionSummary) (logic.ActionSummary, error) {
	action, err := t.effectRepo.Find(pkg.LocalizeTag(d.ActionTag))
	if err != nil {
		return logic.ActionSummary{}, err
	}
	switch logic.SummaryType(d.Type) {
	case logic.SUMMARY_TYPE_ACT:
		fallthrough
	case logic.SUMMARY_TYPE_MODIFY:
		fallthrough
	case logic.SUMMARY_TYPE_SOLVE:
		return logic.ActionSummary{Action: action, Type: logic.SummaryType(d.Type), Data: d.Data}, ErrNotFound
	}
	return logic.ActionSummary{}, ErrNotFound
}

func (t *ActionSummaryTranslator) FromModel(m logic.ActionSummary) (data.ActionSummary, error) {
	return data.ActionSummary{
		ActionTag: string(m.Action.LocalizeTag()),
		Type:      string(m.Type),
		Data:      m.Data,
	}, nil
}
