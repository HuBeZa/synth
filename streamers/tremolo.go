package streamers

import (
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
)

func Tremolo(streamer beep.Streamer, length int, startGain, endGain float64, pulsing bool) beep.Streamer {
	tremoloFunc := func(l, r, progress float64) (float64, float64) {
		gain := startGain + (endGain-startGain)*progress
		return l * gain, r * gain
	}

	if pulsing {
		return NewEffectLoop(streamer, length, effects.TransitionEqualPower, tremoloFunc)
	}
	return NewEffectLoop(streamer, length*2, TransitionLoop, tremoloFunc)
}