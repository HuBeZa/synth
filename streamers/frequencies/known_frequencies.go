package frequencies

import "slices"

// source: https://en.wikipedia.org/wiki/Scientific_pitch_notation#Table_of_note_frequencies

const (
	unknownMidiId = -999
)

var (
	knownFrequencies = []Frequency{
		CMinus1(), DFlatMinus1(), DMinus1(), DSharpMinus1(), EMinus1(), FMinus1(), GFlatMinus1(), GMinus1(), GSharpMinus1(), AMinus1(), ASharpMinus1(), BMinus1(),
		C0(), DFlat0(), D0(), DSharp0(), E0(), F0(), GFlat0(), G0(), GSharp0(), A0(), ASharp0(), B0(),
		C1(), DFlat1(), D1(), DSharp1(), E1(), F1(), GFlat1(), G1(), GSharp1(), A1(), ASharp1(), B1(),
		C2(), DFlat2(), D2(), DSharp2(), E2(), F2(), GFlat2(), G2(), GSharp2(), A2(), ASharp2(), B2(),
		C3(), DFlat3(), D3(), DSharp3(), E3(), F3(), GFlat3(), G3(), GSharp3(), A3(), ASharp3(), B3(),
		C4(), DFlat4(), D4(), DSharp4(), E4(), F4(), GFlat4(), G4(), GSharp4(), A4(), ASharp4(), B4(),
		C5(), DFlat5(), D5(), DSharp5(), E5(), F5(), GFlat5(), G5(), GSharp5(), A5(), ASharp5(), B5(),
		C6(), DFlat6(), D6(), DSharp6(), E6(), F6(), GFlat6(), G6(), GSharp6(), A6(), ASharp6(), B6(),
		C7(), DFlat7(), D7(), DSharp7(), E7(), F7(), GFlat7(), G7(), GSharp7(), A7(), ASharp7(), B7(),
		C8(), DFlat8(), D8(), DSharp8(), E8(), F8(), GFlat8(), G8(), GSharp8(), A8(), ASharp8(), B8(),
		C9(), DFlat9(), D9(), DSharp9(), E9(), F9(), GFlat9(), G9(), GSharp9(), A9(), ASharp9(), B9(),
		C10(), DFlat10(), D10(), DSharp10(), E10(), F10(), GFlat10(), G10(), GSharp10(), A10(), ASharp10(), B10(),
	}

	knownFrequencyIndexes = initKnownFrequencyIndexes()
)

func initKnownFrequencyIndexes() map[float64]int {
	frequencyIndexes := make(map[float64]int, len(knownFrequencies))
	for i, freq := range knownFrequencies {
		frequencyIndexes[freq.Frequency()] = i
	}
	return frequencyIndexes
}

func GetRange(from, to Frequency) []Frequency {
	fromIndex := slices.Index(knownFrequencies, from)
	if fromIndex == -1 {
		return nil
	}
	toIndex := slices.Index(knownFrequencies, to)
	if toIndex == -1 {
		return nil
	}
	reverse := false
	if fromIndex > toIndex {
		reverse = true
		fromIndex, toIndex = toIndex, fromIndex
	}

	src := knownFrequencies[fromIndex : toIndex+1]
	dst := make([]Frequency, len(src))
	copy(dst, src)
	if reverse {
		slices.Reverse(dst)
	}
	return dst
}

func Silence() frequency      { return frequency{`Silence`, 0, unknownMidiId} }
func CMinus1() Frequency      { return frequency{`C-1`, 8.175799, 0} }
func DFlatMinus1() Frequency  { return frequency{`C♯/D♭-1`, 8.661957, 1} }
func CSharpMinus1() Frequency { return DFlatMinus1() }
func DMinus1() Frequency      { return frequency{`D-1`, 9.177024, 2} }
func DSharpMinus1() Frequency { return frequency{`E♭/D♯-1`, 9.722718, 3} }
func EFlatMinus1() Frequency  { return DSharpMinus1() }
func EMinus1() Frequency      { return frequency{`E-1`, 10.30086, 4} }
func FMinus1() Frequency      { return frequency{`F-1`, 10.91338, 5} }
func GFlatMinus1() Frequency  { return frequency{`F♯/G♭-1`, 11.56233, 6} }
func FSharpMinus1() Frequency { return GFlatMinus1() }
func GMinus1() Frequency      { return frequency{`G-1`, 12.24986, 7} }
func GSharpMinus1() Frequency { return frequency{`A♭/G♯-1`, 12.97827, 8} }
func AFlatMinus1() Frequency  { return GSharpMinus1() }
func AMinus1() Frequency      { return frequency{`A-1`, 13.75000, 9} }
func ASharpMinus1() Frequency { return frequency{`B♭/A♯-1`, 14.56762, 10} }
func BFlatMinus1() Frequency  { return ASharpMinus1() }
func BMinus1() Frequency      { return frequency{`B-1`, 15.43385, 11} }

func C0() Frequency      { return frequency{`C0`, 16.35160, 12} }
func DFlat0() Frequency  { return frequency{`C♯/D♭0`, 17.32391, 13} }
func CSharp0() Frequency { return DFlat0() }
func D0() Frequency      { return frequency{`D0`, 18.35405, 14} }
func DSharp0() Frequency { return frequency{`E♭/D♯0`, 19.44544, 15} }
func EFlat0() Frequency  { return DSharp0() }
func E0() Frequency      { return frequency{`E0`, 20.60172, 16} }
func F0() Frequency      { return frequency{`F0`, 21.82676, 17} }
func GFlat0() Frequency  { return frequency{`F♯/G♭0`, 23.12465, 18} }
func FSharp0() Frequency { return GFlat0() }
func G0() Frequency      { return frequency{`G0`, 24.49971, 19} }
func GSharp0() Frequency { return frequency{`A♭/G♯0`, 25.95654, 20} }
func AFlat0() Frequency  { return GSharp0() }
func A0() Frequency      { return frequency{`A0`, 27.50000, 21} }
func ASharp0() Frequency { return frequency{`B♭/A♯0`, 29.13524, 22} }
func BFlat0() Frequency  { return ASharp0() }
func B0() Frequency      { return frequency{`B0`, 30.86771, 23} }

func C1() Frequency      { return frequency{`C1`, 32.70320, 24} }
func DFlat1() Frequency  { return frequency{`C♯/D♭1`, 34.64783, 25} }
func CSharp1() Frequency { return DFlat1() }
func D1() Frequency      { return frequency{`D1`, 36.70810, 26} }
func DSharp1() Frequency { return frequency{`E♭/D♯1`, 38.89087, 27} }
func EFlat1() Frequency  { return DSharp1() }
func E1() Frequency      { return frequency{`E1`, 41.20344, 28} }
func F1() Frequency      { return frequency{`F1`, 43.65353, 29} }
func GFlat1() Frequency  { return frequency{`F♯/G♭1`, 46.24930, 30} }
func FSharp1() Frequency { return GFlat1() }
func G1() Frequency      { return frequency{`G1`, 48.99943, 31} }
func GSharp1() Frequency { return frequency{`A♭/G♯1`, 51.91309, 32} }
func AFlat1() Frequency  { return GSharp1() }
func A1() Frequency      { return frequency{`A1`, 55.00000, 33} }
func ASharp1() Frequency { return frequency{`B♭/A♯1`, 58.27047, 34} }
func BFlat1() Frequency  { return ASharp1() }
func B1() Frequency      { return frequency{`B1`, 61.73541, 35} }

func C2() Frequency      { return frequency{`C2`, 65.40639, 36} }
func DFlat2() Frequency  { return frequency{`C♯/D♭2`, 69.29566, 37} }
func CSharp2() Frequency { return DFlat2() }
func D2() Frequency      { return frequency{`D2`, 73.41619, 38} }
func DSharp2() Frequency { return frequency{`E♭/D♯2`, 77.78175, 39} }
func EFlat2() Frequency  { return DSharp2() }
func E2() Frequency      { return frequency{`E2`, 82.40689, 40} }
func F2() Frequency      { return frequency{`F2`, 87.30706, 41} }
func GFlat2() Frequency  { return frequency{`F♯/G♭2`, 92.49861, 42} }
func FSharp2() Frequency { return GFlat2() }
func G2() Frequency      { return frequency{`G2`, 97.99886, 43} }
func GSharp2() Frequency { return frequency{`A♭/G♯2`, 103.8262, 44} }
func AFlat2() Frequency  { return GSharp2() }
func A2() Frequency      { return frequency{`A2`, 110.0000, 45} }
func ASharp2() Frequency { return frequency{`B♭/A♯2`, 116.5409, 46} }
func BFlat2() Frequency  { return ASharp2() }
func B2() Frequency      { return frequency{`B2`, 123.4708, 47} }

func C3() Frequency      { return frequency{`C3`, 130.8128, 48} }
func DFlat3() Frequency  { return frequency{`C♯/D♭3`, 138.5913, 49} }
func CSharp3() Frequency { return DFlat3() }
func D3() Frequency      { return frequency{`D3`, 146.8324, 50} }
func DSharp3() Frequency { return frequency{`E♭/D♯3`, 155.5635, 51} }
func EFlat3() Frequency  { return DSharp3() }
func E3() Frequency      { return frequency{`E3`, 164.8138, 52} }
func F3() Frequency      { return frequency{`F3`, 174.6141, 53} }
func GFlat3() Frequency  { return frequency{`F♯/G♭3`, 184.9972, 54} }
func FSharp3() Frequency { return GFlat3() }
func G3() Frequency      { return frequency{`G3`, 195.9977, 55} }
func GSharp3() Frequency { return frequency{`A♭/G♯3`, 207.6523, 56} }
func AFlat3() Frequency  { return GSharp3() }
func A3() Frequency      { return frequency{`A3`, 220.0000, 57} }
func ASharp3() Frequency { return frequency{`B♭/A♯3`, 233.0819, 58} }
func BFlat3() Frequency  { return ASharp3() }
func B3() Frequency      { return frequency{`B3`, 246.9417, 59} }

func C4() Frequency      { return frequency{`C4`, 261.6256, 60} }
func DFlat4() Frequency  { return frequency{`C♯/D♭4`, 277.1826, 61} }
func CSharp4() Frequency { return DFlat4() }
func D4() Frequency      { return frequency{`D4`, 293.6648, 62} }
func DSharp4() Frequency { return frequency{`E♭/D♯4`, 311.1270, 63} }
func EFlat4() Frequency  { return DSharp4() }
func E4() Frequency      { return frequency{`E4`, 329.6276, 64} }
func F4() Frequency      { return frequency{`F4`, 349.2282, 65} }
func GFlat4() Frequency  { return frequency{`F♯/G♭4`, 369.9944, 66} }
func FSharp4() Frequency { return GFlat4() }
func G4() Frequency      { return frequency{`G4`, 391.9954, 67} }
func GSharp4() Frequency { return frequency{`A♭/G♯4`, 415.3047, 68} }
func AFlat4() Frequency  { return GSharp4() }
func A4() Frequency      { return frequency{`A4`, 440.0000, 69} }
func ASharp4() Frequency { return frequency{`B♭/A♯4`, 466.1638, 70} }
func BFlat4() Frequency  { return ASharp4() }
func B4() Frequency      { return frequency{`B4`, 493.8833, 71} }

func C5() Frequency      { return frequency{`C5`, 523.2511, 72} }
func DFlat5() Frequency  { return frequency{`C♯/D♭5`, 554.3653, 73} }
func CSharp5() Frequency { return DFlat5() }
func D5() Frequency      { return frequency{`D5`, 587.3295, 74} }
func DSharp5() Frequency { return frequency{`E♭/D♯5`, 622.2540, 75} }
func EFlat5() Frequency  { return DSharp5() }
func E5() Frequency      { return frequency{`E5`, 659.2551, 76} }
func F5() Frequency      { return frequency{`F5`, 698.4565, 77} }
func GFlat5() Frequency  { return frequency{`F♯/G♭5`, 739.9888, 78} }
func FSharp5() Frequency { return GFlat5() }
func G5() Frequency      { return frequency{`G5`, 783.9909, 79} }
func GSharp5() Frequency { return frequency{`A♭/G♯5`, 830.6094, 80} }
func AFlat5() Frequency  { return GSharp5() }
func A5() Frequency      { return frequency{`A5`, 880.0000, 81} }
func ASharp5() Frequency { return frequency{`B♭/A♯5`, 932.3275, 82} }
func BFlat5() Frequency  { return ASharp5() }
func B5() Frequency      { return frequency{`B5`, 987.7666, 83} }

func C6() Frequency      { return frequency{`C6`, 1046.502, 84} }
func DFlat6() Frequency  { return frequency{`C♯/D♭6`, 1108.731, 85} }
func CSharp6() Frequency { return DFlat6() }
func D6() Frequency      { return frequency{`D6`, 1174.659, 86} }
func DSharp6() Frequency { return frequency{`E♭/D♯6`, 1244.508, 87} }
func EFlat6() Frequency  { return DSharp6() }
func E6() Frequency      { return frequency{`E6`, 1318.510, 88} }
func F6() Frequency      { return frequency{`F6`, 1396.913, 89} }
func GFlat6() Frequency  { return frequency{`F♯/G♭6`, 1479.978, 90} }
func FSharp6() Frequency { return GFlat6() }
func G6() Frequency      { return frequency{`G6`, 1567.982, 91} }
func GSharp6() Frequency { return frequency{`A♭/G♯6`, 1661.219, 92} }
func AFlat6() Frequency  { return GSharp6() }
func A6() Frequency      { return frequency{`A6`, 1760.000, 93} }
func ASharp6() Frequency { return frequency{`B♭/A♯6`, 1864.655, 94} }
func BFlat6() Frequency  { return ASharp6() }
func B6() Frequency      { return frequency{`B6`, 1975.533, 95} }

func C7() Frequency      { return frequency{`C7`, 2093.005, 96} }
func DFlat7() Frequency  { return frequency{`C♯/D♭7`, 2217.461, 97} }
func CSharp7() Frequency { return DFlat7() }
func D7() Frequency      { return frequency{`D7`, 2349.318, 98} }
func DSharp7() Frequency { return frequency{`E♭/D♯7`, 2489.016, 99} }
func EFlat7() Frequency  { return DSharp7() }
func E7() Frequency      { return frequency{`E7`, 2637.020, 100} }
func F7() Frequency      { return frequency{`F7`, 2793.826, 101} }
func GFlat7() Frequency  { return frequency{`F♯/G♭7`, 2959.955, 102} }
func FSharp7() Frequency { return GFlat7() }
func G7() Frequency      { return frequency{`G7`, 3135.963, 103} }
func GSharp7() Frequency { return frequency{`A♭/G♯7`, 3322.438, 104} }
func AFlat7() Frequency  { return GSharp7() }
func A7() Frequency      { return frequency{`A7`, 3520.000, 105} }
func ASharp7() Frequency { return frequency{`B♭/A♯7`, 3729.310, 106} }
func BFlat7() Frequency  { return ASharp7() }
func B7() Frequency      { return frequency{`B7`, 3951.066, 107} }

func C8() Frequency      { return frequency{`C8`, 4186.009, 108} }
func DFlat8() Frequency  { return frequency{`C♯/D♭8`, 4434.922, 109} }
func CSharp8() Frequency { return DFlat8() }
func D8() Frequency      { return frequency{`D8`, 4698.636, 110} }
func DSharp8() Frequency { return frequency{`E♭/D♯8`, 4978.032, 111} }
func EFlat8() Frequency  { return DSharp8() }
func E8() Frequency      { return frequency{`E8`, 5274.041, 112} }
func F8() Frequency      { return frequency{`F8`, 5587.652, 113} }
func GFlat8() Frequency  { return frequency{`F♯/G♭8`, 5919.911, 114} }
func FSharp8() Frequency { return GFlat8() }
func G8() Frequency      { return frequency{`G8`, 6271.927, 115} }
func GSharp8() Frequency { return frequency{`A♭/G♯8`, 6644.875, 116} }
func AFlat8() Frequency  { return GSharp8() }
func A8() Frequency      { return frequency{`A8`, 7040.000, 117} }
func ASharp8() Frequency { return frequency{`B♭/A♯8`, 7458.620, 118} }
func BFlat8() Frequency  { return ASharp8() }
func B8() Frequency      { return frequency{`B8`, 7902.133, 119} }

func C9() Frequency      { return frequency{`C9`, 8372.018, 120} }
func DFlat9() Frequency  { return frequency{`C♯/D♭9`, 8869.844, 121} }
func CSharp9() Frequency { return DFlat9() }
func D9() Frequency      { return frequency{`D9`, 9397.273, 122} }
func DSharp9() Frequency { return frequency{`E♭/D♯9`, 9956.063, 123} }
func EFlat9() Frequency  { return DSharp9() }
func E9() Frequency      { return frequency{`E9`, 10548.08, 124} }
func F9() Frequency      { return frequency{`F9`, 11175.30, 125} }
func GFlat9() Frequency  { return frequency{`F♯/G♭9`, 11839.82, 126} }
func FSharp9() Frequency { return GFlat9() }
func G9() Frequency      { return frequency{`G9`, 12543.85, 127} }
func GSharp9() Frequency { return frequency{`A♭/G♯9`, 13289.75, unknownMidiId} }
func AFlat9() Frequency  { return GSharp9() }
func A9() Frequency      { return frequency{`A9`, 14080.00, unknownMidiId} }
func ASharp9() Frequency { return frequency{`B♭/A♯9`, 14917.24, unknownMidiId} }
func BFlat9() Frequency  { return ASharp9() }
func B9() Frequency      { return frequency{`B9`, 15804.27, unknownMidiId} }

func C10() Frequency      { return frequency{`C10`, 16744.04, unknownMidiId} }
func DFlat10() Frequency  { return frequency{`C♯/D♭10`, 17739.69, unknownMidiId} }
func CSharp10() Frequency { return DFlat10() }
func D10() Frequency      { return frequency{`D10`, 18794.55, unknownMidiId} }
func DSharp10() Frequency { return frequency{`E♭/D♯10`, 19912.13, unknownMidiId} }
func EFlat10() Frequency  { return DSharp10() }
func E10() Frequency      { return frequency{`E10`, 21096.16, unknownMidiId} }
func F10() Frequency      { return frequency{`F10`, 22350.61, unknownMidiId} }
func GFlat10() Frequency  { return frequency{`F♯/G♭10`, 23679.64, unknownMidiId} }
func FSharp10() Frequency { return GFlat10() }
func G10() Frequency      { return frequency{`G10`, 25087.71, unknownMidiId} }
func GSharp10() Frequency { return frequency{`A♭/G♯10`, 26579.50, unknownMidiId} }
func AFlat10() Frequency  { return GSharp10() }
func A10() Frequency      { return frequency{`A10`, 28160.00, unknownMidiId} }
func ASharp10() Frequency { return frequency{`B♭/A♯10`, 29834.48, unknownMidiId} }
func BFlat10() Frequency  { return ASharp10() }
func B10() Frequency      { return frequency{`B10`, 31608.53, unknownMidiId} }
