package streamers

import (
	"fmt"
	"strconv"

	"github.com/gopxl/beep/v2/generators"
)

type Waveform int

const (
	Unknown Waveform = iota - 1
	Sine
	Triangle
	Square
	Sawtooth
	ReversedSawtooth
)

var (
	waveformToStr = map[Waveform]string{
		Sine:             "sine",
		Triangle:         "triangle",
		Square:           "square",
		Sawtooth:         "sawtooth",
		ReversedSawtooth: "reversed sawtooth",
	}

	waveformsToGenerator = map[Waveform]StreamerGeneratorFunc{
		Sine:             generators.SineTone,
		Triangle:         generators.TriangleTone,
		Square:           generators.SquareTone,
		Sawtooth:         generators.SawtoothTone,
		ReversedSawtooth: generators.SawtoothToneReversed,
	}
)

func AllWaveforms() []Waveform {
	return []Waveform{Sine, Triangle, Square, Sawtooth, ReversedSawtooth}
}

func (w Waveform) String() string {
	if s, ok := waveformToStr[w]; ok {
		return s
	}
	return strconv.Itoa(int(w))
}

func (w Waveform) streamerGenerator() (StreamerGeneratorFunc, error) {
	generator, ok := waveformsToGenerator[w]
	if !ok {
		return nil, fmt.Errorf("waveform unknown")
	}
	return generator, nil
}
