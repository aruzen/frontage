package translator

import (
	"frontage/pkg"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"frontage/pkg/network/data"
	"frontage/pkg/network/repository"
)

type ActionResultTranslator struct {
	effectRepo *repository.ActionRepository
}

type ActionSummaryTranslator struct {
	effectRepo *repository.ActionRepository
}

func NewActionResultTranslator(repo *repository.ActionRepository) *ActionResultTranslator {
	return &ActionResultTranslator{
		effectRepo: repo,
	}
}

func (t *ActionResultTranslator) ToModel(b *model.Board, d data.ActionResult) (logic.ActionResult, error) {
	action, err := t.effectRepo.FindEffect(logic.EffectActionTag(d.ActionTag))
	if err != nil {
		return logic.ActionResult{}, err
	}
	context := action.MakeContext(d.Context)
	if context == nil {
		return logic.ActionResult{}, ErrNewContextFailed
	}
	state := action.MakeState(b, d.State)
	if state == nil {
		return logic.ActionResult{}, ErrNewStateFailed
	}
	return logic.ActionResult{Action: action, State: state, Context: context, SummaryIdx: d.SummaryIdx}, nil
}

func (t *ActionResultTranslator) FromModel(l logic.ActionResult) (data.ActionResult, error) {
	result := data.ActionResult{}
	result.SummaryIdx = l.SummaryIdx
	result.Context = l.Context.ToMap()
	result.State = l.State.ToMap()
	result.ActionTag = string(l.Action.LocalizeTag())
	return result, nil
}

func NewActionSummaryTranslator(effectRepo *repository.ActionRepository) *ActionSummaryTranslator {
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
		return logic.ActionSummary{Action: action, Type: logic.SummaryType(d.Type), Data: d.Data}, nil
	}
	return logic.ActionSummary{}, ErrNotFound
}

func (t *ActionSummaryTranslator) FromModel(l logic.ActionSummary) (data.ActionSummary, error) {
	summary := data.ActionSummary{}
	summary.Data = l.Data
	summary.ActionTag = string(l.Action.LocalizeTag())
	summary.Type = string(l.Type)
	return summary, nil
}
