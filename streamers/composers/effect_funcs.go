package composers

type EffectFunc func(l, r, progress float64) (float64, float64)

func GainTransitionEffect(startGain, endGain float64) EffectFunc {
	return func(l, r, progress float64) (float64, float64) {
		gain := startGain + (endGain-startGain)*progress
		return l * gain, r * gain
	}
}

// TODO: add PanTransitionEffect
