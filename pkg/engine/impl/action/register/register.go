package register

import (
	"frontage/pkg/engine/impl/action"
	"frontage/pkg/engine/impl/action/card_action"
	"frontage/pkg/engine/impl/action/common"
	"frontage/pkg/engine/impl/action/piece_action"
	"frontage/pkg/engine/logic"
	"log/slog"
	"sync"
)

var once sync.Once

// Init registers all known effect actions once.
func Init() {
	once.Do(func() {
		for _, a := range logic.EnumerateEffectAction() {
			if err := action.Register(a.Tag(), a); err != nil {
				slog.Warn("failed to register system action", "tag", a.Tag(), "err", err)
			}
		}
		for _, a := range card_action.EnumerateEffectAction() {
			if err := action.Register(a.Tag(), a); err != nil {
				slog.Warn("failed to register card action", "tag", a.Tag(), "err", err)
			}
		}
		for _, a := range piece_action.EnumerateEffectAction() {
			if err := action.Register(a.Tag(), a); err != nil {
				slog.Warn("failed to register piece action", "tag", a.Tag(), "err", err)
			}
		}
		for _, a := range common.EnumerateModifyAction() {
			if err := action.RegisterModify(a.Tag(), a); err != nil {
				slog.Warn("failed to register modify action", "tag", a.Tag(), "err", err)
			}
		}
	})
}
