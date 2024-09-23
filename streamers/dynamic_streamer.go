package streamers

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/HuBeZa/synth/streamers/chords"
	"github.com/HuBeZa/synth/streamers/frequencies"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/generators"
)

var (
	silenceStreamer = generators.Silence(-1)
)

type StreamerGeneratorFunc func(sampleRate beep.SampleRate, freq float64) (beep.Streamer, error)

type DynamicStreamer interface {
	beep.Streamer
	IsSilenced() bool
	Silence()
	Unsilence()
	ToggleSilence()
	Pan() float64
	SetPan(pan float64) error
	Gain() float64
	SetGain(gain float64) error
	Frequency() frequencies.Frequency
	SetFrequency(freq frequencies.Frequency) error
	SetTremolo(duration time.Duration, startGain, endGain float64, pulsing bool) error
	SetTremoloOff() error
	SetEnvelop(attack time.Duration, decay time.Duration, sustain float64, release time.Duration) error
	SetChord(chord chords.ChordType, arpeggioDelay time.Duration) error
	SetChordOff() error
	SetOvertones(count int, gain float64) error
	Waveform() Waveform
	SetWaveform(waveform Waveform) error
	SetGenerator(streamerGenerator StreamerGeneratorFunc) error
	TriggerAttack()
	TriggerRelease()
}

type dynamicStreamer struct {
	streamerArgs streamerArgs
	waveform     Waveform
	silenced     atomic.Bool
	isReleased   bool
	streamer     atomic.Pointer[beep.Streamer]

	// additional tones effects:
	chordOptions struct {
		chord         chords.ChordType
		arpeggioDelay time.Duration
	}
	overtones struct {
		count int
		gain  float64
	}
}

type streamerArgs struct {
	sampleRate beep.SampleRate
	generator  StreamerGeneratorFunc
	frequency  frequencies.Frequency
	pan        float64
	gain       float64

	tremolo struct {
		isOn      bool
		length    int
		startGain float64
		endGain   float64
		pulsing   bool
	}

	envelop struct {
		isOn    bool
		attack  int
		decay   int
		sustain float64
		release int
	}
}

func NewWaveformDynamicStreamer(sampleRate beep.SampleRate, freq frequencies.Frequency, pan, gain float64, waveform Waveform) (DynamicStreamer, error) {
	generator, err := waveform.streamerGenerator()
	if err != nil {
		return nil, err
	}

	return NewDynamicStreamer(sampleRate, freq, pan, gain, generator)
}

func NewDynamicStreamer(sampleRate beep.SampleRate, freq frequencies.Frequency, pan, gain float64, streamerGenerator StreamerGeneratorFunc) (DynamicStreamer, error) {
	s := &dynamicStreamer{
		streamerArgs: streamerArgs{
			sampleRate: sampleRate,
			generator:  streamerGenerator,
			frequency:  freq,
			pan:        pan,
			gain:       gain,
		},
	}

	if err := s.update(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *dynamicStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	if s.IsSilenced() {
		return silenceStreamer.Stream(samples)
	}
	return s.getStreamer().Stream(samples)
}

func (s *dynamicStreamer) Err() error {
	return s.getStreamer().Err()
}

func (s *dynamicStreamer) IsSilenced() bool {
	return s.silenced.Load()
}

func (s *dynamicStreamer) Silence() {
	s.silenced.Store(true)
}

func (s *dynamicStreamer) Unsilence() {
	s.silenced.Store(false)
}

func (s *dynamicStreamer) ToggleSilence() {
	if s.IsSilenced() {
		s.Unsilence()
	} else {
		s.Silence()
	}
}

func (s *dynamicStreamer) Pan() float64 {
	return s.streamerArgs.pan
}

func (s *dynamicStreamer) SetPan(pan float64) error {
	if pan == s.streamerArgs.pan {
		return nil
	}

	orig := s.streamerArgs.pan
	s.streamerArgs.pan = pan
	if err := s.update(); err != nil {
		s.streamerArgs.pan = orig
		return err
	}

	return nil
}

func (s *dynamicStreamer) Gain() float64 {
	return s.streamerArgs.gain
}

func (s *dynamicStreamer) SetGain(gain float64) error {
	if gain == s.streamerArgs.gain {
		return nil
	}

	orig := s.streamerArgs.gain
	s.streamerArgs.gain = gain
	if err := s.update(); err != nil {
		s.streamerArgs.gain = orig
		return err
	}

	return nil
}

func (s *dynamicStreamer) SetTremolo(duration time.Duration, startGain, endGain float64, pulsing bool) error {
	orig := s.streamerArgs.tremolo
	s.streamerArgs.tremolo.isOn = true
	s.streamerArgs.tremolo.length = s.streamerArgs.sampleRate.N(duration)
	s.streamerArgs.tremolo.startGain = startGain
	s.streamerArgs.tremolo.endGain = endGain
	s.streamerArgs.tremolo.pulsing = pulsing

	if orig == s.streamerArgs.tremolo {
		return nil
	}

	if err := s.update(); err != nil {
		s.streamerArgs.tremolo = orig
		return err
	}

	return nil
}

func (s *dynamicStreamer) SetTremoloOff() error {
	if !s.streamerArgs.tremolo.isOn {
		return nil
	}

	s.streamerArgs.tremolo.isOn = false
	if err := s.update(); err != nil {
		s.streamerArgs.tremolo.isOn = true
		return err
	}

	return nil
}

func (s *dynamicStreamer) SetEnvelop(attack time.Duration, decay time.Duration, sustain float64, release time.Duration) error {
	orig := s.streamerArgs.envelop
	s.streamerArgs.envelop.isOn = true
	s.streamerArgs.envelop.attack = s.streamerArgs.sampleRate.N(attack)
	s.streamerArgs.envelop.decay = s.streamerArgs.sampleRate.N(decay)
	s.streamerArgs.envelop.sustain = sustain
	s.streamerArgs.envelop.release = s.streamerArgs.sampleRate.N(release)

	if orig == s.streamerArgs.envelop {
		return nil
	}

	if err := s.update(); err != nil {
		s.streamerArgs.envelop = orig
		return err
	}

	return nil
}

func (s *dynamicStreamer) SetChord(chord chords.ChordType, arpeggioDelay time.Duration) error {
	if chords.Equals(s.chordOptions.chord, chord) && s.chordOptions.arpeggioDelay == arpeggioDelay {
		return nil
	}

	orig := s.chordOptions
	s.chordOptions.chord = chord
	s.chordOptions.arpeggioDelay = arpeggioDelay
	if err := s.update(); err != nil {
		s.chordOptions = orig
		return err
	}

	return nil
}

func (s *dynamicStreamer) SetChordOff() error {
	if s.chordOptions.chord == nil {
		return nil
	}

	orig := s.chordOptions
	s.chordOptions.chord = nil
	if err := s.update(); err != nil {
		s.chordOptions = orig
		return err
	}

	return nil
}

func (s *dynamicStreamer) SetOvertones(count int, gain float64) error {
	if s.overtones.count == count && s.overtones.gain == gain {
		return nil
	}

	orig := s.overtones
	s.overtones.count = count
	s.overtones.gain = gain

	if err := s.update(); err != nil {
		s.overtones = orig
		return err
	}

	return nil
}

func (s *dynamicStreamer) Frequency() frequencies.Frequency {
	return s.streamerArgs.frequency
}

func (s *dynamicStreamer) SetFrequency(freq frequencies.Frequency) error {
	if freq.Frequency() == s.streamerArgs.frequency.Frequency() {
		return nil
	}

	orig := s.streamerArgs.frequency
	s.streamerArgs.frequency = freq
	if err := s.update(); err != nil {
		s.streamerArgs.frequency = orig
		return err
	}

	return nil
}

func (s *dynamicStreamer) Waveform() Waveform {
	return s.waveform
}

func (s *dynamicStreamer) SetWaveform(waveform Waveform) error {
	if waveform == s.waveform {
		return nil
	}

	generator, err := waveform.streamerGenerator()
	if err != nil {
		return err
	}

	return s.setGenerator(waveform, generator)
}

func (s *dynamicStreamer) SetGenerator(streamerGenerator StreamerGeneratorFunc) error {
	return s.setGenerator(Unknown, streamerGenerator)
}

func (s *dynamicStreamer) setGenerator(waveform Waveform, streamerGenerator StreamerGeneratorFunc) error {
	orig := s.streamerArgs.generator
	s.streamerArgs.generator = streamerGenerator
	if err := s.update(); err != nil {
		s.streamerArgs.generator = orig
		return err
	}
	s.waveform = waveform
	return nil
}

func (s *dynamicStreamer) TriggerAttack() {
	s.isReleased = false
	// update() forces restart of time dependent effects
	s.update()
	s.Unsilence()
}

func (s *dynamicStreamer) TriggerRelease() {
	s.isReleased = true
	streamer := SetRelease(s.getStreamer(), s.streamerArgs.envelop.sustain, s.streamerArgs.envelop.release)
	s.streamer.Store(&streamer)
}

func (s *dynamicStreamer) update() error {
	streamer, err := createStreamer(s.streamerArgs)
	if err != nil {
		return err
	}

	if s.chordOptions.chord != nil {
		streamer = s.addChord(streamer)
	}

	if s.overtones.count > 0 {
		streamer = s.addOvertones(streamer)
	}

	if s.isReleased {
		s.Silence()
	}

	s.streamer.Store(&streamer)
	return nil
}

func (s *dynamicStreamer) addChord(rootStreamer beep.Streamer) beep.Streamer {
	mixer := &beep.Mixer{}
	for i, semitone := range s.chordOptions.chord.Semitones() {
		if semitone == 0 {
			mixer.Add(rootStreamer)
			continue
		}

		argsCopy := s.streamerArgs
		argsCopy.frequency = s.streamerArgs.frequency.ShiftSemitone(semitone)
		if semitoneStreamer, err := createStreamer(argsCopy); err == nil {
			if s.chordOptions.arpeggioDelay > 0 {
				// Delay each tone, up to 2 delays. After that play all remaining tones together.
				delay := time.Duration(min(i, 2)) * s.chordOptions.arpeggioDelay
				time.AfterFunc(delay, func() { mixer.Add(semitoneStreamer) })
			} else {
				mixer.Add(semitoneStreamer)
			}
		}
	}

	return mixer
}

func (s *dynamicStreamer) addOvertones(rootStreamer beep.Streamer) beep.Streamer {
	mixer := &beep.Mixer{}
	mixer.Add(rootStreamer)

	for i := 1; i <= s.overtones.count; i++ {
		argsCopy := s.streamerArgs
		argsCopy.frequency = s.streamerArgs.frequency.ShiftOctave(i)
		argsCopy.gain *= s.overtones.gain

		// note that some overtones may not be created because they will overpass sampleRate/2
		if overtone, err := createStreamer(argsCopy); err == nil {
			mixer.Add(overtone)
		}
	}

	return mixer
}

func createStreamer(args streamerArgs) (beep.Streamer, error) {
	if args.generator == nil {
		return nil, fmt.Errorf("streamer generator is empty")
	}
	if args.pan < -1 || args.pan > 1 {
		return nil, fmt.Errorf("pan should be between -1 (left channel) to 1 (right channel)")
	}

	streamer, err := args.generator(args.sampleRate, args.frequency.Frequency())
	if err != nil {
		return nil, err
	}

	if args.pan != 0 {
		streamer = &effects.Pan{
			Streamer: streamer,
			Pan:      args.pan,
		}
	}

	if args.gain != 1 {
		streamer = &effects.Gain{
			Streamer: streamer,
			Gain:     args.gain - 1,
		}
	}

	// streamer = effects.Transition(streamer, s.sampleRate.N(time.Second/2), 2.5, 0, effects.TransitionLinear)

	if args.tremolo.isOn {
		streamer = Tremolo(streamer, args.tremolo.length, args.tremolo.startGain, args.tremolo.endGain, args.tremolo.pulsing)
	}

	if args.envelop.isOn {
		streamer = SetAttackDecaySustain(streamer, args.envelop.attack, args.envelop.decay, args.envelop.sustain)
	}

	return streamer, nil
}

func (s *dynamicStreamer) getStreamer() beep.Streamer {
	return *s.streamer.Load()
}
