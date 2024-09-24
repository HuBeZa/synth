package composers

import "github.com/gopxl/beep/v2/effects"

// All funcs in this file are a beep/v2/effects.TransitionFunc
// Some functions are just a wrapper of beep funcs.

type TransitionFunc effects.TransitionFunc

func TransitionLinear(percent float64) float64 {
	return effects.TransitionLinear(percent)
}

func TransitionEqualPower(percent float64) float64 {
	return effects.TransitionEqualPower(percent)
}

// TODO: add TransitionExponential func

// TransitionLoop runs f up from 0 to 1, and back to zero
func TransitionLoop(f effects.TransitionFunc) TransitionFunc {
	return func(percent float64) float64 {
		if percent <= 0.5 {
			return f(percent * 2)
		}
		return f((1 - percent) * 2)
	}
}
