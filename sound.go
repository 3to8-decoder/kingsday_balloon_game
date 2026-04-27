package main

import (
	"encoding/binary"
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	sampleRate = 44100
)

// generatePop synthesizes a short "pop" sound as 16-bit PCM.
// Frequency rises from baseFreq, amplitude decays exponentially.
func generatePop(baseFreq float64, duration float64) []byte {
	numSamples := int(sampleRate * duration)
	data := make([]byte, numSamples*4) // 2 channels × 2 bytes

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)

		// Frequency sweep up
		freq := baseFreq + 400*t
		// Exponential decay
		amp := math.Exp(-t * 15)

		val := math.Sin(2 * math.Pi * freq * t) * amp

		// Convert to 16-bit integer
		sample := int16(val * 0.5 * 32767)
		if sample > 32767 {
			sample = 32767
		} else if sample < -32767 {
			sample = -32767
		}

		// Write to both channels (stereo, little-endian)
		offset := i * 4
		binary.LittleEndian.PutUint16(data[offset:], uint16(sample))
		binary.LittleEndian.PutUint16(data[offset+2:], uint16(sample))
	}

	return data
}

// generateCombo synthesizes a pleasant "ding" for combo hits.
func generateCombo() []byte {
	numSamples := int(sampleRate * 0.25)
	data := make([]byte, numSamples*4)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)

		// Two harmonics for a nice ding
		fundamental := math.Sin(2*math.Pi*880*t) * math.Exp(-t*8)
		harmonic := math.Sin(2*math.Pi*1320*t) * math.Exp(-t*12)
		val := fundamental*0.6 + harmonic*0.4

		sample := int16(val * 0.4 * 32767)

		offset := i * 4
		binary.LittleEndian.PutUint16(data[offset:], uint16(sample))
		binary.LittleEndian.PutUint16(data[offset+2:], uint16(sample))
	}

	return data
}

// generateGameOver synthesizes a descending "wah wah" effect.
func generateGameOver() []byte {
	numSamples := int(sampleRate * 0.6)
	data := make([]byte, numSamples*4)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)

		// Frequency sweep down
		freq := 400 - 300*t
		amp := math.Exp(-t * 4)
		val := math.Sin(2*math.Pi*freq*t) * amp

		sample := int16(val * 0.5 * 32767)

		offset := i * 4
		binary.LittleEndian.PutUint16(data[offset:], uint16(sample))
		binary.LittleEndian.PutUint16(data[offset+2:], uint16(sample))
	}

	return data
}

type SoundPlayer struct {
	ctx        *audio.Context
	popBytes   []byte
	comboBytes []byte
	overBytes  []byte
}

func NewSoundPlayer() *SoundPlayer {
	ctx := audio.NewContext(sampleRate)
	return &SoundPlayer{
		ctx:        ctx,
		popBytes:   generatePop(600, 0.1),
		comboBytes: generateCombo(),
		overBytes:  generateGameOver(),
	}
}

func (s *SoundPlayer) PlayPop() {
	p := s.ctx.NewPlayerFromBytes(s.popBytes)
	p.Play()
}

func (s *SoundPlayer) PlayCombo() {
	p := s.ctx.NewPlayerFromBytes(s.comboBytes)
	p.Play()
}

func (s *SoundPlayer) PlayGameOver() {
	p := s.ctx.NewPlayerFromBytes(s.overBytes)
	p.Play()
}
