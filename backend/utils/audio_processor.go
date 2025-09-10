package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// ProcessAudio handles the audio processing using the new modular approach
func ProcessAudio(inputFiles []string, outputFile string, loops int, crossfadeDuration float64) error {
	return ProcessAudioWithOptions(inputFiles, outputFile, loops, crossfadeDuration, true, "mp3")
}

// ProcessAudioWithOptions handles audio processing with enhancement and format options
func ProcessAudioWithOptions(inputFiles []string, outputFile string, loops int, crossfadeDuration float64, enhance bool, format string) error {
	if len(inputFiles) == 0 {
		return fmt.Errorf("no input files provided")
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputFile)
	os.MkdirAll(outputDir, 0755)

	// Create audio manager with temp directory
	tempDir := filepath.Join(outputDir, "temp")
	audioManager := NewAudioManager(tempDir)

	// Process audio sequence with crossfades, loops, and enhancement
	return audioManager.ProcessAudioSequenceWithOptions(inputFiles, outputFile, loops, crossfadeDuration, enhance, format)
}
