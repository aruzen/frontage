package logic

// EnumerateEffectAction returns all system-level EffectActions.
func EnumerateEffectAction() []EffectAction {
	return []EffectAction{
		GAME_START_ACTION,
		GAME_FINISH_ACTION,
		TURN_START_ACTION,
		TURN_END_ACTION,
		PLAYER_WIN_ACTION,
		PLAYER_LOSE_ACTION,
	}
}
