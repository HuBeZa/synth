package frequencies

import (
	"fmt"
	"math"
)

// octave = 12 semitones
// octave multiplier = 2
// semitone multiplier = 12th root of 2 = ~1.0595
var semitoneMultiplier = math.Pow(2, 1.0/12)

type Frequency interface {
	Name() string
	Frequency() float64
	MidiID() int
	ShiftSemitone(semitonesDiff int) Frequency
	ShiftOctave(octavesDiff int) Frequency
	fmt.Stringer
}

type frequency struct {
	name      string
	frequency float64
	midiId    int
}

func New(freq float64) Frequency {
	if i, ok := knownFrequencyIndexes[freq]; ok {
		// return known frequency
		return knownFrequencies[i]
	}

	return frequency{
		name:      fmt.Sprintf("%vHz", freq),
		frequency: freq,
		midiId:    unknownMidiId,
	}
}

func (f frequency) Name() string {
	return f.name
}

func (f frequency) Frequency() float64 {
	return f.frequency
}

func (f frequency) MidiID() int {
	return f.midiId
}

func (f frequency) String() string {
	return f.name
}

func (f frequency) ShiftSemitone(semitonesDiff int) Frequency {
	if semitonesDiff == 0 {
		return f
	}

	if i, ok := knownFrequencyIndexes[f.frequency]; ok {
		newIndex := i + semitonesDiff
		if newIndex >= 0 && newIndex < len(knownFrequencies) {
			return knownFrequencies[newIndex]
		}
	}

	return shiftSemitoneUnknownFreq(f.frequency, semitonesDiff)
}

func (f frequency) ShiftOctave(octavesDiff int) Frequency {
	return f.ShiftSemitone(octavesDiff * 12)
}

func shiftSemitoneUnknownFreq(freq float64, semitonesDiff int) Frequency {
	return New(freq * math.Pow(semitoneMultiplier, float64(semitonesDiff)))
}
