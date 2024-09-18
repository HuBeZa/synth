package streamers

import (
	"math"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
)

// TransitionLoop is a beep/v2/effects.TransitionFunc.
// For the transition period it runs linearly from zero to one, and then back to zero.
func TransitionLoop(percent float64) float64 {
	if percent <= 0.5 {
		return percent*2
	}
	return (1-percent)*2
}

type EffectFunc func(l, r, progress float64) (float64, float64)

type effectLoop struct {
	streamer       beep.Streamer
	pos            int
	length         int
	transitionFunc effects.TransitionFunc
	effectFunc     EffectFunc
}

func NewEffectLoop(streamer beep.Streamer, length int, transitionFunc effects.TransitionFunc, effectFunc EffectFunc) beep.Streamer {
	return &effectLoop{
		streamer:       streamer,
		length:         length,
		transitionFunc: transitionFunc,
		effectFunc:     effectFunc,
	}
}

func (e *effectLoop) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = e.streamer.Stream(samples)

	for i := 0; i < n; i++ {
		pos := e.pos + i
		progress := float64(pos) / float64(e.length)
		if progress > 1 {
			_, progress = math.Modf(progress)
		}

		progress = e.transitionFunc(progress)
		samples[i][0], samples[i][1] = e.effectFunc(samples[i][0], samples[i][1], progress)
	}

	e.pos += n
	e.pos %= e.length

	return
}

func (e *effectLoop) Err() error {
	return e.streamer.Err()
}