package logic

type MultiEffectAction interface {
	SubEffects(state interface{}) []EffectEvent
}
