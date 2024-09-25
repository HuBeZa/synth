package streamers

import (
	"github.com/gopxl/beep/v2"

	"github.com/HuBeZa/synth/streamers/composers"
)

func attackFunc() composers.EffectFunc {
	return composers.GainTransitionEffect(0, 1)
}

func decayFunc(sustain float64) composers.EffectFunc {
	return composers.GainTransitionEffect(1, sustain)
}

func releaseFunc(sustain float64) composers.EffectFunc {
	return composers.GainTransitionEffect(sustain, 0)
}

func SetAttackDecaySustain(streamer beep.Streamer, attack int, attackType composers.TransitionType, decay int, decayType composers.TransitionType, sustain float64) beep.Streamer {
	if attack == 0 && decay == 0 && sustain == 1 {
		return streamer
	}

	return composers.NewEffectsChain(streamer).
		Append(attack, attackType.Func(), attackFunc()).
		Append(decay, decayType.Func(), decayFunc(sustain)).
		Loop(false).
		Build()
}

func SetRelease(streamer beep.Streamer, sustain float64, release int, releaseType composers.TransitionType) beep.Streamer {
	if release == 0 {
		return silenceStreamer
	}

	return composers.NewEffectsChain(streamer).
		Append(release, releaseType.Func(), releaseFunc(sustain)).
		Loop(false).
		Build()
}
