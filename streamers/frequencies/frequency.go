package frequencies

import (
	"fmt"
	"math"
	"slices"
)

const (
	unknownMidiId = -999
)

// source: https://en.wikipedia.org/wiki/Scientific_pitch_notation#Table_of_note_frequencies
// TODO: change to funcs to make readonly
var (
	CMinus1      Frequency = frequency{`C-1`, 8.175799, 0}
	DFlatMinus1  Frequency = frequency{`C♯/D♭-1`, 8.661957, 1}
	CSharpMinus1 Frequency = DFlatMinus1
	DMinus1      Frequency = frequency{`D-1`, 9.177024, 2}
	DSharpMinus1 Frequency = frequency{`E♭/D♯-1`, 9.722718, 3}
	EFlatMinus1  Frequency = DSharpMinus1
	EMinus1      Frequency = frequency{`E-1`, 10.30086, 4}
	FMinus1      Frequency = frequency{`F-1`, 10.91338, 5}
	GFlatMinus1  Frequency = frequency{`F♯/G♭-1`, 11.56233, 6}
	FSharpMinus1 Frequency = GFlatMinus1
	GMinus1      Frequency = frequency{`G-1`, 12.24986, 7}
	GSharpMinus1 Frequency = frequency{`A♭/G♯-1`, 12.97827, 8}
	AFlatMinus1  Frequency = GSharpMinus1
	AMinus1      Frequency = frequency{`A-1`, 13.75000, 9}
	ASharpMinus1 Frequency = frequency{`B♭/A♯-1`, 14.56762, 10}
	BFlatMinus1  Frequency = ASharpMinus1
	BMinus1      Frequency = frequency{`B-1`, 15.43385, 11}

	C0      Frequency = frequency{`C0`, 16.35160, 12}
	DFlat0  Frequency = frequency{`C♯/D♭0`, 17.32391, 13}
	CSharp0 Frequency = DFlat0
	D0      Frequency = frequency{`D0`, 18.35405, 14}
	DSharp0 Frequency = frequency{`E♭/D♯0`, 19.44544, 15}
	EFlat0  Frequency = DSharp0
	E0      Frequency = frequency{`E0`, 20.60172, 16}
	F0      Frequency = frequency{`F0`, 21.82676, 17}
	GFlat0  Frequency = frequency{`F♯/G♭0`, 23.12465, 18}
	FSharp0 Frequency = GFlat0
	G0      Frequency = frequency{`G0`, 24.49971, 19}
	GSharp0 Frequency = frequency{`A♭/G♯0`, 25.95654, 20}
	AFlat0  Frequency = GSharp0
	A0      Frequency = frequency{`A0`, 27.50000, 21}
	ASharp0 Frequency = frequency{`B♭/A♯0`, 29.13524, 22}
	BFlat0  Frequency = ASharp0
	B0      Frequency = frequency{`B0`, 30.86771, 23}

	C1      Frequency = frequency{`C1`, 32.70320, 24}
	DFlat1  Frequency = frequency{`C♯/D♭1`, 34.64783, 25}
	CSharp1 Frequency = DFlat1
	D1      Frequency = frequency{`D1`, 36.70810, 26}
	DSharp1 Frequency = frequency{`E♭/D♯1`, 38.89087, 27}
	EFlat1  Frequency = DSharp1
	E1      Frequency = frequency{`E1`, 41.20344, 28}
	F1      Frequency = frequency{`F1`, 43.65353, 29}
	GFlat1  Frequency = frequency{`F♯/G♭1`, 46.24930, 30}
	FSharp1 Frequency = GFlat1
	G1      Frequency = frequency{`G1`, 48.99943, 31}
	GSharp1 Frequency = frequency{`A♭/G♯1`, 51.91309, 32}
	AFlat1  Frequency = GSharp1
	A1      Frequency = frequency{`A1`, 55.00000, 33}
	ASharp1 Frequency = frequency{`B♭/A♯1`, 58.27047, 34}
	BFlat1  Frequency = ASharp1
	B1      Frequency = frequency{`B1`, 61.73541, 35}

	C2      Frequency = frequency{`C2`, 65.40639, 36}
	DFlat2  Frequency = frequency{`C♯/D♭2`, 69.29566, 37}
	CSharp2 Frequency = DFlat2
	D2      Frequency = frequency{`D2`, 73.41619, 38}
	DSharp2 Frequency = frequency{`E♭/D♯2`, 77.78175, 39}
	EFlat2  Frequency = DSharp2
	E2      Frequency = frequency{`E2`, 82.40689, 40}
	F2      Frequency = frequency{`F2`, 87.30706, 41}
	GFlat2  Frequency = frequency{`F♯/G♭2`, 92.49861, 42}
	FSharp2 Frequency = GFlat2
	G2      Frequency = frequency{`G2`, 97.99886, 43}
	GSharp2 Frequency = frequency{`A♭/G♯2`, 103.8262, 44}
	AFlat2  Frequency = GSharp2
	A2      Frequency = frequency{`A2`, 110.0000, 45}
	ASharp2 Frequency = frequency{`B♭/A♯2`, 116.5409, 46}
	BFlat2  Frequency = ASharp2
	B2      Frequency = frequency{`B2`, 123.4708, 47}

	C3      Frequency = frequency{`C3`, 130.8128, 48}
	DFlat3  Frequency = frequency{`C♯/D♭3`, 138.5913, 49}
	CSharp3 Frequency = DFlat3
	D3      Frequency = frequency{`D3`, 146.8324, 50}
	DSharp3 Frequency = frequency{`E♭/D♯3`, 155.5635, 51}
	EFlat3  Frequency = DSharp3
	E3      Frequency = frequency{`E3`, 164.8138, 52}
	F3      Frequency = frequency{`F3`, 174.6141, 53}
	GFlat3  Frequency = frequency{`F♯/G♭3`, 184.9972, 54}
	FSharp3 Frequency = GFlat3
	G3      Frequency = frequency{`G3`, 195.9977, 55}
	GSharp3 Frequency = frequency{`A♭/G♯3`, 207.6523, 56}
	AFlat3  Frequency = GSharp3
	A3      Frequency = frequency{`A3`, 220.0000, 57}
	ASharp3 Frequency = frequency{`B♭/A♯3`, 233.0819, 58}
	BFlat3  Frequency = ASharp3
	B3      Frequency = frequency{`B3`, 246.9417, 59}

	C4      Frequency = frequency{`C4`, 261.6256, 60}
	DFlat4  Frequency = frequency{`C♯/D♭4`, 277.1826, 61}
	CSharp4 Frequency = DFlat4
	D4      Frequency = frequency{`D4`, 293.6648, 62}
	DSharp4 Frequency = frequency{`E♭/D♯4`, 311.1270, 63}
	EFlat4  Frequency = DSharp4
	E4      Frequency = frequency{`E4`, 329.6276, 64}
	F4      Frequency = frequency{`F4`, 349.2282, 65}
	GFlat4  Frequency = frequency{`F♯/G♭4`, 369.9944, 66}
	FSharp4 Frequency = GFlat4
	G4      Frequency = frequency{`G4`, 391.9954, 67}
	GSharp4 Frequency = frequency{`A♭/G♯4`, 415.3047, 68}
	AFlat4  Frequency = GSharp4
	A4      Frequency = frequency{`A4`, 440.0000, 69}
	ASharp4 Frequency = frequency{`B♭/A♯4`, 466.1638, 70}
	BFlat4  Frequency = ASharp4
	B4      Frequency = frequency{`B4`, 493.8833, 71}

	C5      Frequency = frequency{`C5`, 523.2511, 72}
	DFlat5  Frequency = frequency{`C♯/D♭5`, 554.3653, 73}
	CSharp5 Frequency = DFlat5
	D5      Frequency = frequency{`D5`, 587.3295, 74}
	DSharp5 Frequency = frequency{`E♭/D♯5`, 622.2540, 75}
	EFlat5  Frequency = DSharp5
	E5      Frequency = frequency{`E5`, 659.2551, 76}
	F5      Frequency = frequency{`F5`, 698.4565, 77}
	GFlat5  Frequency = frequency{`F♯/G♭5`, 739.9888, 78}
	FSharp5 Frequency = GFlat5
	G5      Frequency = frequency{`G5`, 783.9909, 79}
	GSharp5 Frequency = frequency{`A♭/G♯5`, 830.6094, 80}
	AFlat5  Frequency = GSharp5
	A5      Frequency = frequency{`A5`, 880.0000, 81}
	ASharp5 Frequency = frequency{`B♭/A♯5`, 932.3275, 82}
	BFlat5  Frequency = ASharp5
	B5      Frequency = frequency{`B5`, 987.7666, 83}

	C6      Frequency = frequency{`C6`, 1046.502, 84}
	DFlat6  Frequency = frequency{`C♯/D♭6`, 1108.731, 85}
	CSharp6 Frequency = DFlat6
	D6      Frequency = frequency{`D6`, 1174.659, 86}
	DSharp6 Frequency = frequency{`E♭/D♯6`, 1244.508, 87}
	EFlat6  Frequency = DSharp6
	E6      Frequency = frequency{`E6`, 1318.510, 88}
	F6      Frequency = frequency{`F6`, 1396.913, 89}
	GFlat6  Frequency = frequency{`F♯/G♭6`, 1479.978, 90}
	FSharp6 Frequency = GFlat6
	G6      Frequency = frequency{`G6`, 1567.982, 91}
	GSharp6 Frequency = frequency{`A♭/G♯6`, 1661.219, 92}
	AFlat6  Frequency = GSharp6
	A6      Frequency = frequency{`A6`, 1760.000, 93}
	ASharp6 Frequency = frequency{`B♭/A♯6`, 1864.655, 94}
	BFlat6  Frequency = ASharp6
	B6      Frequency = frequency{`B6`, 1975.533, 95}

	C7      Frequency = frequency{`C7`, 2093.005, 96}
	DFlat7  Frequency = frequency{`C♯/D♭7`, 2217.461, 97}
	CSharp7 Frequency = DFlat7
	D7      Frequency = frequency{`D7`, 2349.318, 98}
	DSharp7 Frequency = frequency{`E♭/D♯7`, 2489.016, 99}
	EFlat7  Frequency = DSharp7
	E7      Frequency = frequency{`E7`, 2637.020, 100}
	F7      Frequency = frequency{`F7`, 2793.826, 101}
	GFlat7  Frequency = frequency{`F♯/G♭7`, 2959.955, 102}
	FSharp7 Frequency = GFlat7
	G7      Frequency = frequency{`G7`, 3135.963, 103}
	GSharp7 Frequency = frequency{`A♭/G♯7`, 3322.438, 104}
	AFlat7  Frequency = GSharp7
	A7      Frequency = frequency{`A7`, 3520.000, 105}
	ASharp7 Frequency = frequency{`B♭/A♯7`, 3729.310, 106}
	BFlat7  Frequency = ASharp7
	B7      Frequency = frequency{`B7`, 3951.066, 107}

	C8      Frequency = frequency{`C8`, 4186.009, 108}
	DFlat8  Frequency = frequency{`C♯/D♭8`, 4434.922, 109}
	CSharp8 Frequency = DFlat8
	D8      Frequency = frequency{`D8`, 4698.636, 110}
	DSharp8 Frequency = frequency{`E♭/D♯8`, 4978.032, 111}
	EFlat8  Frequency = DSharp8
	E8      Frequency = frequency{`E8`, 5274.041, 112}
	F8      Frequency = frequency{`F8`, 5587.652, 113}
	GFlat8  Frequency = frequency{`F♯/G♭8`, 5919.911, 114}
	FSharp8 Frequency = GFlat8
	G8      Frequency = frequency{`G8`, 6271.927, 115}
	GSharp8 Frequency = frequency{`A♭/G♯8`, 6644.875, 116}
	AFlat8  Frequency = GSharp8
	A8      Frequency = frequency{`A8`, 7040.000, 117}
	ASharp8 Frequency = frequency{`B♭/A♯8`, 7458.620, 118}
	BFlat8  Frequency = ASharp8
	B8      Frequency = frequency{`B8`, 7902.133, 119}

	C9      Frequency = frequency{`C9`, 8372.018, 120}
	DFlat9  Frequency = frequency{`C♯/D♭9`, 8869.844, 121}
	CSharp9 Frequency = DFlat9
	D9      Frequency = frequency{`D9`, 9397.273, 122}
	DSharp9 Frequency = frequency{`E♭/D♯9`, 9956.063, 123}
	EFlat9  Frequency = DSharp9
	E9      Frequency = frequency{`E9`, 10548.08, 124}
	F9      Frequency = frequency{`F9`, 11175.30, 125}
	GFlat9  Frequency = frequency{`F♯/G♭9`, 11839.82, 126}
	FSharp9 Frequency = GFlat9
	G9      Frequency = frequency{`G9`, 12543.85, 127}
	GSharp9 Frequency = frequency{`A♭/G♯9`, 13289.75, unknownMidiId}
	AFlat9  Frequency = GSharp9
	A9      Frequency = frequency{`A9`, 14080.00, unknownMidiId}
	ASharp9 Frequency = frequency{`B♭/A♯9`, 14917.24, unknownMidiId}
	BFlat9  Frequency = ASharp9
	B9      Frequency = frequency{`B9`, 15804.27, unknownMidiId}

	C10      Frequency = frequency{`C10`, 16744.04, unknownMidiId}
	DFlat10  Frequency = frequency{`C♯/D♭10`, 17739.69, unknownMidiId}
	CSharp10 Frequency = DFlat10
	D10      Frequency = frequency{`D10`, 18794.55, unknownMidiId}
	DSharp10 Frequency = frequency{`E♭/D♯10`, 19912.13, unknownMidiId}
	EFlat10  Frequency = DSharp10
	E10      Frequency = frequency{`E10`, 21096.16, unknownMidiId}
	F10      Frequency = frequency{`F10`, 22350.61, unknownMidiId}
	GFlat10  Frequency = frequency{`F♯/G♭10`, 23679.64, unknownMidiId}
	FSharp10 Frequency = GFlat10
	G10      Frequency = frequency{`G10`, 25087.71, unknownMidiId}
	GSharp10 Frequency = frequency{`A♭/G♯10`, 26579.50, unknownMidiId}
	AFlat10  Frequency = GSharp10
	A10      Frequency = frequency{`A10`, 28160.00, unknownMidiId}
	ASharp10 Frequency = frequency{`B♭/A♯10`, 29834.48, unknownMidiId}
	BFlat10  Frequency = ASharp10
	B10      Frequency = frequency{`B10`, 31608.53, unknownMidiId}

	frequencies = []Frequency{
		CMinus1, DFlatMinus1, DMinus1, DSharpMinus1, EMinus1, FMinus1, GFlatMinus1, GMinus1, GSharpMinus1, AMinus1, ASharpMinus1, BMinus1,
		C0, DFlat0, D0, DSharp0, E0, F0, GFlat0, G0, GSharp0, A0, ASharp0, B0,
		C1, DFlat1, D1, DSharp1, E1, F1, GFlat1, G1, GSharp1, A1, ASharp1, B1,
		C2, DFlat2, D2, DSharp2, E2, F2, GFlat2, G2, GSharp2, A2, ASharp2, B2,
		C3, DFlat3, D3, DSharp3, E3, F3, GFlat3, G3, GSharp3, A3, ASharp3, B3,
		C4, DFlat4, D4, DSharp4, E4, F4, GFlat4, G4, GSharp4, A4, ASharp4, B4,
		C5, DFlat5, D5, DSharp5, E5, F5, GFlat5, G5, GSharp5, A5, ASharp5, B5,
		C6, DFlat6, D6, DSharp6, E6, F6, GFlat6, G6, GSharp6, A6, ASharp6, B6,
		C7, DFlat7, D7, DSharp7, E7, F7, GFlat7, G7, GSharp7, A7, ASharp7, B7,
		C8, DFlat8, D8, DSharp8, E8, F8, GFlat8, G8, GSharp8, A8, ASharp8, B8,
		C9, DFlat9, D9, DSharp9, E9, F9, GFlat9, G9, GSharp9, A9, ASharp9, B9,
		C10, DFlat10, D10, DSharp10, E10, F10, GFlat10, G10, GSharp10, A10, ASharp10, B10,
	}

	frequencyIndexes = initFrequencyIndexes()

	semitoneMultiplier = math.Pow(2, 1.0/12)
)

func initFrequencyIndexes() map[float64]int {
	frequencyIndexes := make(map[float64]int, len(frequencies))
	for i, freq := range frequencies {
		frequencyIndexes[freq.Frequency()] = i
	}
	return frequencyIndexes
}

type Frequency interface {
	Name() string
	Frequency() float64
	MidiID() int
	ShiftSemitone(semitonesDiff int) Frequency
	ShiftOctave(octavesDiff int) Frequency
	fmt.Stringer
}

func New(freq float64) Frequency {
	if i, ok := frequencyIndexes[freq]; ok {
		// return known frequency
		return frequencies[i]
	}

	return frequency{
		name:      fmt.Sprintf("%vHz", freq),
		frequency: freq,
		midiId:    unknownMidiId,
	}
}

func GetRange(from, to Frequency) []Frequency {
	fromIndex := slices.Index(frequencies, from)
	if fromIndex == -1 {
		return nil
	}
	toIndex := slices.Index(frequencies, to)
	if toIndex == -1 {
		return nil
	}
	reverse := false
	if fromIndex > toIndex {
		reverse = true
		fromIndex, toIndex = toIndex, fromIndex
	}

	src := frequencies[fromIndex : toIndex+1]
	dst := make([]Frequency, len(src))
	copy(dst, src)
	if reverse {
		slices.Reverse(dst)
	}
	return dst
}

type frequency struct {
	name      string
	frequency float64
	midiId    int
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
	i, ok := frequencyIndexes[f.frequency]
	if !ok {
		return shiftSemitone(f.frequency, semitonesDiff)
	}

	newIndex := i + semitonesDiff
	if newIndex < 0 || newIndex >= len(frequencies) {
		return shiftSemitone(f.frequency, semitonesDiff)
	}

	return frequencies[newIndex]
}

func (f frequency) ShiftOctave(octavesDiff int) Frequency {
	return f.ShiftSemitone(octavesDiff * 12)
}

func shiftSemitone(freq float64, semitonesDiff int) Frequency {
	return New(freq * math.Pow(semitoneMultiplier, float64(semitonesDiff)))
}
