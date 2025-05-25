package sound

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log"
	"path"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/vinser/pacmantea/internal/embeddata"
)

// Sound names
const (
	BEGINNING    = "pacman_beginning.wav"
	CHOMP        = "pacman_chomp.wav"
	DEATH        = "pacman_death.wav"
	EATFRUIT     = "pacman_eatfruit.wav"
	EATGHOST     = "pacman_eatghost.wav"
	EXTRAPAC     = "pacman_extrapac.wav"
	INTERMISSION = "pacman_intermission.wav"
)

const commonSampleRate = 44100 // Common sample rate for normalization for all sounds

type Sound struct {
	Name   string
	Stream beep.StreamSeekCloser
	Format beep.Format
}

func LoadSamples() (map[string]Sound, error) {
	speaker.Init(beep.SampleRate(commonSampleRate), beep.SampleRate(commonSampleRate).N(time.Second/10))

	// Open the embedded ZIP archive
	soundsZip, err := embeddata.ReadSoundsZip()
	if err != nil {
		log.Fatalf("Failed to read embedded sounds.zip: %v", err)
	}
	reader, err := zip.NewReader(bytes.NewReader(soundsZip), int64(len(soundsZip)))
	if err != nil {
		fmt.Println("Error reading embedded ZIP:", err)
		return nil, err
	}

	// Map to store loaded sounds
	sounds := make(map[string]Sound)

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
		case BEGINNING, CHOMP, DEATH, EATFRUIT, EATGHOST, EXTRAPAC, INTERMISSION:
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
			sounds[fileName] = Sound{
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

func Play(sound Sound) {
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

func ClearSpeaker() {
	speaker.Clear()
}
