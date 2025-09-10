package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AudioManager coordinates all audio processing operations
type AudioManager struct {
	Validator *AudioValidator
	Sequencer *AudioSequencer
	TempDir   string
}

// NewAudioManager creates a new audio manager
func NewAudioManager(tempDir string) *AudioManager {
	return &AudioManager{
		Validator: NewAudioValidator(),
		TempDir:   tempDir,
	}
}

// ProcessAudioSequence handles the complete audio processing pipeline
func (am *AudioManager) ProcessAudioSequence(inputFiles []string, outputFile string, loops int, crossfadeDuration float64) error {
	return am.ProcessAudioSequenceWithOptions(inputFiles, outputFile, loops, crossfadeDuration, true, "mp3")
}

// ProcessAudioSequenceWithOptions handles audio processing with enhancement and format options
func (am *AudioManager) ProcessAudioSequenceWithOptions(inputFiles []string, outputFile string, loops int, crossfadeDuration float64, enhance bool, format string) error {
	// Step 1: Validate all input files
	if err := am.Validator.ValidateFiles(inputFiles); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

	// Step 2: Create unique session directory
	sessionID := fmt.Sprintf("%d", time.Now().UnixNano())
	sessionDir := filepath.Join(am.TempDir, sessionID)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return fmt.Errorf("failed to create session directory: %v", err)
	}
	defer os.RemoveAll(sessionDir)

	// Step 3: Initialize sequencer with options
	am.Sequencer = NewAudioSequencerWithOptions(inputFiles, outputFile, crossfadeDuration, loops, sessionDir, enhance, format)

	// Step 4: Process the sequence
	if err := am.Sequencer.Process(); err != nil {
		return fmt.Errorf("sequencing failed: %v", err)
	}

	return nil
}

// GetAudioFileInfo returns information about an audio file
func (am *AudioManager) GetAudioFileInfo(filePath string) (map[string]string, error) {
	return am.Validator.GetAudioInfo(filePath)
}

// ValidateAudioFile validates a single audio file
func (am *AudioManager) ValidateAudioFile(filePath string) error {
	return am.Validator.ValidateFile(filePath)
}
