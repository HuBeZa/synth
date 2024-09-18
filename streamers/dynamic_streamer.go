package streamers

import (
	"fmt"
	"sync/atomic"
	"time"

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
	SetTremolo(duration time.Duration, startGain, endGain float64, pulsing bool) error
	SetTremoloOff() error
	Frequency() float64
	SetFrequency(freq float64) error
	Waveform() Waveform
	SetWaveform(waveform Waveform) error
	SetGenerator(streamerGenerator StreamerGeneratorFunc) error
}

type dynamicStreamer struct {
	silenced   atomic.Bool
	sampleRate beep.SampleRate
	pan        float64
	gain       float64
	frequency  float64
	waveform   Waveform
	generator  StreamerGeneratorFunc
	streamer   atomic.Pointer[beep.Streamer]

	tremolo struct {
		isOn      bool
		length    int
		startGain float64
		endGain   float64
		pulsing   bool
	}
}

func NewWaveformDynamicStreamer(sampleRate beep.SampleRate, freq, pan, gain float64, waveform Waveform) (DynamicStreamer, error) {
	generator, err := waveform.streamerGenerator()
	if err != nil {
		return nil, err
	}

	return NewDynamicStreamer(sampleRate, freq, pan, gain, generator)
}

func NewDynamicStreamer(sampleRate beep.SampleRate, freq, pan, gain float64, streamerGenerator StreamerGeneratorFunc) (DynamicStreamer, error) {
	s := &dynamicStreamer{
		sampleRate: sampleRate,
		generator:  streamerGenerator,
		frequency:  freq,
		pan:        pan,
		gain:       gain,
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
	return (*s.streamer.Load()).Stream(samples)
}

func (s *dynamicStreamer) Err() error {
	return (*s.streamer.Load()).Err()
}

func (s *dynamicStreamer) IsSilenced() bool {
	return s.silenced.Load()
}

func (s *dynamicStreamer) Silence() {
	s.silenced.Store(true)
}

func (s *dynamicStreamer) Unsilence() {
	// force restart of effects
	s.update()

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
	return s.pan
}

func (s *dynamicStreamer) SetPan(pan float64) error {
	if pan == s.pan {
		return nil
	}

	orig := s.pan
	s.pan = pan
	if err := s.update(); err != nil {
		s.pan = orig
		return err
	}

	return nil
}

func (s *dynamicStreamer) Gain() float64 {
	return s.gain
}

func (s *dynamicStreamer) SetGain(gain float64) error {
	if gain == s.gain {
		return nil
	}

	orig := s.gain
	s.gain = gain
	if err := s.update(); err != nil {
		s.gain = orig
		return err
	}

	return nil
}

func (s *dynamicStreamer) SetTremolo(duration time.Duration, startGain, endGain float64, pulsing bool) error {
	orig := s.tremolo
	s.tremolo.isOn = true
	s.tremolo.length = s.sampleRate.N(duration)
	s.tremolo.startGain = startGain
	s.tremolo.endGain = endGain
	s.tremolo.pulsing = pulsing

	if orig == s.tremolo {
		return nil
	}

	if err := s.update(); err != nil {
		s.tremolo = orig
		return err
	}

	return nil
}

func (s *dynamicStreamer) SetTremoloOff() error {
	if !s.tremolo.isOn {
		return nil
	}

	s.tremolo.isOn = false
	if err := s.update(); err != nil {
		s.tremolo.isOn = true
		return err
	}
	
	return nil
}

func (s *dynamicStreamer) Frequency() float64 {
	return s.frequency
}

func (s *dynamicStreamer) SetFrequency(freq float64) error {
	if freq == s.frequency {
		return nil
	}

	orig := s.frequency
	s.frequency = freq
	if err := s.update(); err != nil {
		s.frequency = orig
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
	orig := s.generator
	s.generator = streamerGenerator
	if err := s.update(); err != nil {
		s.generator = orig
		return err
	}
	s.waveform = waveform
	return nil
}

func (s *dynamicStreamer) update() error {
	if s.generator == nil {
		return fmt.Errorf("streamer generator is empty")
	}
	if s.pan < -1 || s.pan > 1 {
		return fmt.Errorf("pan should be between -1 (left channel) to 1 (right channel)")
	}

	streamer, err := s.generator(s.sampleRate, s.frequency)
	if err != nil {
		return err
	}

	if s.pan != 0 {
		streamer = &effects.Pan{
			Streamer: streamer,
			Pan:      s.pan,
		}
	}

	if s.gain != 0 {
		streamer = &effects.Gain{
			Streamer: streamer,
			Gain:     s.gain,
		}
	}

	// streamer = effects.Transition(streamer, s.sampleRate.N(time.Second/2), 2.5, 0, effects.TransitionLinear)
	// streamer = Tremolo(streamer, s.sampleRate.N(time.Second/2), 0.5, 1, false)

	if s.tremolo.isOn {
		streamer = Tremolo(streamer, s.tremolo.length, s.tremolo.startGain, s.tremolo.endGain, s.tremolo.pulsing)
	}

	s.streamer.Store(&streamer)
	return nil
}
