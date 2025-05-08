package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

//go:embed sounds.zip
var soundsZip []byte

const (
	SOUND_BEGINNING    = "pacman_beginning.wav"
	SOUND_CHOMP        = "pacman_chomp.wav"
	SOUND_DEATH        = "pacman_death.wav"
	SOUND_EATFRUIT     = "pacman_eatfruit.wav"
	SOUND_EATGHOST     = "pacman_eatghost.wav"
	SOUND_EXTRAPAC     = "pacman_extrapac.wav"
	SOUND_INTERMISSION = "pacman_intermission.wav"
)

const commonSampleRate = 44100 // Common sample rate for normalization

type sound struct {
	Name   string
	Stream beep.StreamSeekCloser
	Format beep.Format
}

// Updated loadSoundSamples function
func loadSoundSamples() (map[string]sound, error) {
	// Common sample rate for all sounds
	const commonSampleRate = 44100
	speaker.Init(beep.SampleRate(commonSampleRate), beep.SampleRate(commonSampleRate).N(time.Second/10))

	// Open the embedded ZIP archive
	reader, err := zip.NewReader(bytes.NewReader(soundsZip), int64(len(soundsZip)))
	if err != nil {
		fmt.Println("Error reading embedded ZIP:", err)
		return nil, err
	}

	// Map to store loaded sounds
	sounds := make(map[string]sound)

	// Iterate through the files in the ZIP archive
	for _, file := range reader.File {
		// Open the file
		rc, err := file.Open()
		if err != nil {
			fmt.Println("Error opening file in ZIP:", err)
			return nil, err
		}
		defer rc.Close()

		// Check if the file matches one of the sound constants
		fileName := path.Base(file.Name)
		switch fileName {
		case SOUND_BEGINNING, SOUND_CHOMP, SOUND_DEATH, SOUND_EATFRUIT, SOUND_EATGHOST, SOUND_EXTRAPAC, SOUND_INTERMISSION:
			// Load the file into memory to enable seeking
			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, rc)
			if err != nil {
				return nil, fmt.Errorf("error reading file into memory: %w", err)
			}

			// Decode the WAV file
			stream, format, err := wav.Decode(bytes.NewReader(buf.Bytes()))
			if err != nil {
				fmt.Println("Error decoding WAV file:", err)
				return nil, err
			}

			// Add the sound to the map
			sounds[fileName] = sound{
				Name:   fileName,
				Stream: stream,
				Format: beep.Format{
					SampleRate:  format.SampleRate,
					NumChannels: format.NumChannels,
					Precision:   format.Precision,
				},
			}
		}
	}

	return sounds, nil
}

func (m *model) playSound(name string) {
	if sound, ok := m.sounds[name]; ok && !m.mute {
		// Reset the stream to the beginning
		err := sound.Stream.Seek(0)
		if err != nil {
			fmt.Println("Error seeking stream:", err)
			return
		}

		// Play the sound
		done := make(chan bool)
		speaker.Play(beep.Seq(beep.Resample(1, sound.Format.SampleRate, beep.SampleRate(commonSampleRate), sound.Stream), beep.Callback(func() {
			done <- true
		})))

		// Wait for the sound to finish
		<-done
	}
}
