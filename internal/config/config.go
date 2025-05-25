package config

import (
	"log"
	"os"
	"path"

	"github.com/vinser/pacmantea/internal/embeddata"
	"gopkg.in/yaml.v3"
)

type Level struct {
	Name           string   `yaml:"name"`
	DifficultyName string   `yaml:"difficulty"`
	Maze           []string `yaml:"maze"`
	PacmanBadge    string   `yaml:"pacman_badge"` // Badge style for Pac-Man
	GhostBadges    string   `yaml:"ghost_badges"` // Badge style for ghosts
}

type Difficulty struct {
	GhostSpeed       int `yaml:"ghost_speed"`
	RampantDuration  int `yaml:"rampant_duration"`
	CooldownDuration int `yaml:"cooldown_duration"`
	RevivalTimer     int `yaml:"revival_timer"`
	SpeedBonus       int `yaml:"speed_bonus"` // base points for each second in formula speedBonus * wonGames * seconds
}

type Badges struct {
	Pacman map[string]map[string]string `yaml:"pacman"` // Badge styles for Pac-Man
	Ghosts map[string]map[string]string `yaml:"ghosts"` // Badge styles for ghosts
}

type Config struct {
	Badges       Badges                `yaml:"badges"`
	Difficulties map[string]Difficulty `yaml:"difficulties"`
	Levels       []Level               `yaml:"levels"`
}

func WriteDefaultConfig() error {
	// Write the embedded CONFIG_DATA to a file
	if err := os.MkdirAll("config", os.ModePerm); err != nil {
		return err
	}
	file, err := os.Create(path.Join("config", "config.yml"))
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := embeddata.ReadConfig()
	if err != nil {
		log.Fatalf("Failed to read embedded config: %v", err)
	}

	_, err = file.Write(data)
	return err
}

func Load() Config {
	var config Config
	// Determine the config path based on the environment variable
	configPath := os.Getenv("PACMANTEA_CONFIG_PATH")
	if configPath == "" {
		configPath = path.Join("config", "config.yml")
	}

	// Load the configuration from the file or use the embedded one as a fallback
	var data []byte
	var err error
	if _, err = os.Stat(configPath); err == nil {
		data, err = os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("Failed to read %s: %v", configPath, err)
		}
	} else {
		data, err = embeddata.ReadConfig()
		if err != nil {
			log.Fatalf("Failed to read embedded config: %v", err)
		}
	}

	// Unmarshal the configuration
	if err = yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to parse %s: %v", configPath, err)
	}
	return config
}
