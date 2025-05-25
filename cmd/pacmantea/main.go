package main

import (
	"flag"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinser/pacmantea/internal/config"
	"github.com/vinser/pacmantea/internal/model"
	"github.com/vinser/pacmantea/internal/sound"
)

func main() {
	// Define the -config flag
	configFlag := flag.Bool("config", false, "Generate a default config.yml file in the config directory")
	flag.Parse()

	// If -config flag is set, write the default config.yml and exit
	if *configFlag {
		err := config.WriteDefaultConfig()
		if err != nil {
			log.Fatalf("Failed to write default config.yml: %v", err)
		}
		fmt.Println("Default config.yml has been written to the config directory.")
		return
	}

	// Run the game
	model := model.New()
	model.PlaySound(sound.BEGINNING)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
