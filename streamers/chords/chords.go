package chords

import "slices"

var (
	chordMajor = chordType{"Major", "M", []int{0, 4, 7}}
	chordAug   = chordType{"Augmented", "aug", []int{0, 4, 8}}
	// a.k.a sus, sus4
	chord4th = chordType{"Forth", "4", []int{0, 5, 7}}
	chord6th = chordType{"Sixth", "6", []int{0, 4, 7, 9}}
	// a.k.a dominant seventh
	chord7th      = chordType{"Seventh", "7", []int{0, 4, 7, 10}}
	chordMajor7th = chordType{"Major Seventh", "maj7", []int{0, 4, 7, 11}}

	chordMinor    = chordType{"Minor", "m", []int{0, 3, 7}}
	chordMinor7th = chordType{"Minor Seventh", "m7", []int{0, 3, 7, 10}}
	chordDim      = chordType{"Diminished", "dim", []int{0, 3, 6}}
)

type ChordType interface {
	Name() string
	Symbol() string
	Semitones() []int
	Equals(other ChordType) bool
}

type chordType struct {
	name      string
	symbol    string
	semitones []int
}

func (c chordType) Name() string {
	return c.name
}

func (c chordType) Symbol() string {
	return c.symbol
}

func (c chordType) Semitones() []int {
	semitonesCopy := make([]int, len(c.semitones))
	copy(semitonesCopy, c.semitones)
	return semitonesCopy
}

func (c chordType) String() string {
	return c.symbol
}

func (c chordType) Equals(other ChordType) bool {
	if other == nil {
		return false
	}
	return slices.Equal(c.semitones, other.Semitones())
}

func Equals(x, y ChordType) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	return x.Equals(y)
}

func ChordTypes() []ChordType {
	return []ChordType{chordMajor, chordMinor, chord4th, chord6th, chord7th, chordMajor7th, chordMinor7th, chordAug, chordDim}
}
