package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
)

var encryptionKey = getEncryptionKey()

// --- Game State  ---
type State struct {
	Mute        bool           `json:"mute"`         // Disable sound effects
	LevelName   string         `json:"level_name"`   // Current level
	GamesWon    int            `json:"games won"`    // Total number of games won
	HighScore   int            `json:"high_score"`   // Global high score
	ElapsedTime map[string]int `json:"elapsed_time"` // Per-level elapsed time records in seconds by level name
}

// ========================
// 🛡️ Security Functions
// ========================

// Generate encryption key from system-specific data
func getEncryptionKey() []byte {
	user, _ := os.UserHomeDir()
	hash := sha256.Sum256([]byte(user + "packman-secret-salt!"))
	return hash[:] // 32 bytes for AES-256
}

// Encrypt data with AES-GCM
func encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// Decrypt data with AES-GCM
func decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("invalid ciphertext")
	}

	return gcm.Open(nil, ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():], nil)
}

// ========================
// 💾 Save/Load Functions
// ========================

func saveState(s State) error {
	// Get save file path
	filename, err := getSavePath()
	if err != nil {
		return err
	}
	// Serialize game data to JSON
	jsonData, err := json.Marshal(s)
	if err != nil {
		return err
	}

	// Calculate CRC32 checksum
	crc := crc32.ChecksumIEEE(jsonData)

	// Create buffer: [4-byte CRC32][JSON data]
	buf := make([]byte, 4+len(jsonData))
	binary.LittleEndian.PutUint32(buf[:4], crc)
	copy(buf[4:], jsonData)

	// Encrypt and write to file
	encrypted, err := encrypt(buf)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, encrypted, 0644)
}

func loadState() State {
	var state State
	state.ElapsedTime = make(map[string]int)
	// Get save file path
	filename, err := getSavePath()
	if err != nil {
		return state
	}
	// Read encrypted file
	encrypted, err := os.ReadFile(filename)
	if err != nil {
		return state
	}

	// Decrypt data
	decrypted, err := decrypt(encrypted)
	if err != nil {
		return state
	}

	// Verify minimum length (4 bytes CRC32 + at least 1 byte data)
	if len(decrypted) < 5 {
		return state
	}

	// Extract CRC32 (first 4 bytes)
	loadedCRC := binary.LittleEndian.Uint32(decrypted[:4])
	jsonData := decrypted[4:]

	// Verify checksum
	if crc32.ChecksumIEEE(jsonData) != loadedCRC {
		return state
	}

	// Deserialize JSON
	json.Unmarshal(jsonData, &state)

	return state
}

func getSavePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	// Create save directory if it doesn't exist
	saveDir := filepath.Join(configDir, "pacmantea")
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return "", err
	}

	// Path to save binary file
	return filepath.Join(saveDir, "savegame.dat"), nil
}
