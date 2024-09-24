package composers

import (
	"math"

	"github.com/gopxl/beep/v2"
)

type effectLoop struct {
	streamer       beep.Streamer
	pos            int
	length         int
	transitionFunc TransitionFunc
	effectFunc     EffectFunc
}

func NewEffectLoop(streamer beep.Streamer, length int, transitionFunc TransitionFunc, effectFunc EffectFunc) beep.Streamer {
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
