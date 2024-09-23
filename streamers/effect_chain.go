package streamers

import (
	"math"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
)

type EffectsChainBuilder interface {
	Append(length int, transitionFunc effects.TransitionFunc, effectFunc EffectFunc) EffectsChainBuilder
	Loop(loop bool) EffectsChainBuilder
	Build() beep.Streamer
}

type effectArgs struct {
	length         int
	transitionFunc effects.TransitionFunc
	effectFunc     EffectFunc
}

type effectsChain struct {
	streamer    beep.Streamer
	timePos     int
	effectPos   int
	effectFinal bool
	effects     []effectArgs
	loop        bool
}

func NewEffectsChain(streamer beep.Streamer) EffectsChainBuilder {
	return &effectsChain{
		streamer: streamer,
	}
}

func (e *effectsChain) Append(length int, transitionFunc effects.TransitionFunc, effectFunc EffectFunc) EffectsChainBuilder {
	e.effects = append(e.effects, effectArgs{length, transitionFunc, effectFunc})
	return e
}

func (e *effectsChain) Loop(loop bool) EffectsChainBuilder {
	e.loop = loop
	return e
}

func (e *effectsChain) Build() beep.Streamer {
	return e
}

func (e *effectsChain) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = e.streamer.Stream(samples)

	currEffect := e.effects[e.effectPos]
	for i := 0; i < n; i++ {
		progress := 1.0
		if !e.effectFinal {
			e.timePos++
			progress = float64(e.timePos) / float64(currEffect.length)
			if progress > 1 {
				if e.nextEffect() {
					// move to next effect
					currEffect = e.effects[e.effectPos]
					e.timePos = 1
					_, progress = math.Modf(progress)
				} else {
					// sustaining on last effect and final time
					e.effectFinal = true
					progress = 1
				}
			}

			progress = currEffect.transitionFunc(progress)
		}

		samples[i][0], samples[i][1] = currEffect.effectFunc(samples[i][0], samples[i][1], progress)
	}

	return
}

func (e *effectsChain) nextEffect() bool {
	if e.effectPos < len(e.effects)-1 {
		e.effectPos++
		return true
	}
	if e.loop {
		e.effectPos = 0
		return true
	}

	// resume with current effect
	return false
}

func (e *effectsChain) Err() error {
	return e.streamer.Err()
}
