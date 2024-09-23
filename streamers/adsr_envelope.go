package streamers

import (
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
)

func attackFunc() EffectFunc {
	return TransitionEffectFunc(0, 1)
}

func decayFunc(sustain float64) EffectFunc {
	return TransitionEffectFunc(1, sustain)
}

func releaseFunc(sustain float64) EffectFunc {
	return TransitionEffectFunc(sustain, 0)
}

func SetAttackDecaySustain(streamer beep.Streamer, attack int, decay int, sustain float64) beep.Streamer {
	if attack == 0 && decay == 0 && sustain == 1 {
		return streamer
	}

	return NewEffectsChain(streamer).
		Append(attack, effects.TransitionLinear, attackFunc()).
		Append(decay, effects.TransitionLinear, decayFunc(sustain)).
		Loop(false).
		Build()
}

func SetRelease(streamer beep.Streamer, sustain float64, release int) beep.Streamer {
	if release == 0 {
		return silenceStreamer
	}

	return NewEffectsChain(streamer).
		Append(release, effects.TransitionLinear, releaseFunc(sustain)).
		Loop(false).
		Build()
}
