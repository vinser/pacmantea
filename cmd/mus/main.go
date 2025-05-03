package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"

	_ "embed"
)

//go:embed sounds.zip
var soundsZip []byte // Embed the sounds.zip file

const commonSampleRate = 44100 // Common sample rate for normalization

type sound struct {
	Name   string
	Stream beep.StreamSeekCloser
	Format beep.Format
}

func loadSoundSamples() ([]sound, error) {
	// Open the embedded ZIP archive
	reader, err := zip.NewReader(bytes.NewReader(soundsZip), int64(len(soundsZip)))
	if err != nil {
		return nil, fmt.Errorf("error reading embedded ZIP: %w", err)
	}

	var sounds []sound

	// Iterate through the files in the ZIP archive
	for _, file := range reader.File {
		// Open the file
		rc, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("error opening file in ZIP: %w", err)
		}
		defer rc.Close()

		// Load the file into memory to enable seeking
		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, rc)
		if err != nil {
			return nil, fmt.Errorf("error reading file into memory: %w", err)
		}
		switch path.Ext(file.Name) {
		case ".wav":
			// Decode the WAV file
			streamer, format, err := wav.Decode(bytes.NewReader(buf.Bytes()))
			if err != nil {
				return nil, fmt.Errorf("error decoding WAV: %w", err)
			}

			// Add the sound to the list
			sounds = append(sounds, sound{
				Name:   file.Name,
				Stream: streamer,
				Format: beep.Format{
					SampleRate:  format.SampleRate,
					NumChannels: format.NumChannels,
					Precision:   format.Precision,
				},
			})
		}
	}

	return sounds, nil
}

func playSounds(sounds []sound) {
	// Initialize the speaker
	speaker.Init(beep.SampleRate(commonSampleRate), beep.SampleRate(commonSampleRate).N(time.Second/10))

	for _, sound := range sounds {
		playSound(sound)

		// Pause for 2 seconds before playing the next sound
		time.Sleep(2 * time.Second)
	}
}

func playSound(sound sound) {
	fmt.Println("Playing sound:", sound.Name)

	// Reset the stream to the beginning
	err := sound.Stream.Seek(0)
	if err != nil {
		fmt.Println("Error seeking stream:", err)
		return
	}

	// Play the sound
	done := make(chan bool)
	speaker.Play(beep.Seq(beep.Resample(4, sound.Format.SampleRate, beep.SampleRate(commonSampleRate), sound.Stream), beep.Callback(func() {
		done <- true
	})))

	// Wait for the sound to finish
	<-done
}

func main() {
	// Load sound samples
	sounds, err := loadSoundSamples()
	if err != nil {
		fmt.Println("Error loading sounds:", err)
		return
	}

	// Play all sounds one by one
	playSounds(sounds)
}
