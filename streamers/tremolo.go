package streamers

import (
	"github.com/gopxl/beep/v2"

	"github.com/HuBeZa/synth/streamers/composers"
)

func Tremolo(streamer beep.Streamer, length int, startGain, endGain float64, pulsing bool) beep.Streamer {
	tremoloFunc := composers.GainTransitionEffect(startGain, endGain)

	if pulsing {
		return composers.NewEffectLoop(streamer, length, composers.TransitionEqualPower, tremoloFunc)
	}
	return composers.NewEffectLoop(streamer, length*2, composers.TransitionLoop(composers.TransitionLinear), tremoloFunc)
}
