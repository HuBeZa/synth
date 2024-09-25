package composers

import (
	"math"
	"strconv"

	"github.com/gopxl/beep/v2/effects"
)

type TransitionType int

const (
	Linear TransitionType = iota
	EqualPower
	Exponential
)

func (t TransitionType) Func() TransitionFunc {
	switch t {
	case Linear:
		return TransitionLinear
	case EqualPower:
		return TransitionEqualPower
	case Exponential:
		return TransitionExponential
	default:
		return nil
	}
}

func (t TransitionType) String() string {
	switch t {
	case Linear:
		return "linear"
	case EqualPower:
		return "sin"
	case Exponential:
		return "exp"
	default:
		return strconv.Itoa(int(t))
	}
}

func (t TransitionType) Equals(other TransitionType) bool {
	return t == other
}

func TransitionTypes() []TransitionType {
	return []TransitionType{Linear, EqualPower, Exponential}
}

// All funcs in this file are a beep/v2/effects.TransitionFunc
// Some functions are just a wrapper of beep funcs.

type TransitionFunc effects.TransitionFunc

func TransitionLinear(percent float64) float64 {
	return effects.TransitionLinear(percent)
}

func TransitionEqualPower(percent float64) float64 {
	return effects.TransitionEqualPower(percent)
}

func TransitionExponential(percent float64) float64 {
	return math.Pow(percent, 2)
}

// TransitionLoop runs f up from 0 to 1, and back to zero
func TransitionLoop(f effects.TransitionFunc) TransitionFunc {
	return func(percent float64) float64 {
		if percent <= 0.5 {
			return f(percent * 2)
		}
		return f((1 - percent) * 2)
	}
}
