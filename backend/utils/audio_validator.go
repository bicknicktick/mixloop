package utils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// AudioValidator handles validation of audio files
type AudioValidator struct{}

// NewAudioValidator creates a new audio validator
func NewAudioValidator() *AudioValidator {
	return &AudioValidator{}
}

// ValidateFile checks if a file is a valid audio file
func (av *AudioValidator) ValidateFile(filePath string) error {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".mp3" && ext != ".wav" {
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	// Use ffprobe to validate the file
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "a:0", "-show_entries", "stream=codec_type", "-of", "csv=p=0", filePath)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("invalid audio file: %v", err)
	}

	if strings.TrimSpace(string(output)) != "audio" {
		return fmt.Errorf("file does not contain valid audio stream")
	}

	return nil
}

// ValidateFiles validates multiple audio files
func (av *AudioValidator) ValidateFiles(filePaths []string) error {
	for i, filePath := range filePaths {
		if err := av.ValidateFile(filePath); err != nil {
			return fmt.Errorf("file %d (%s): %v", i+1, filepath.Base(filePath), err)
		}
	}
	return nil
}

// GetAudioInfo returns basic information about an audio file
func (av *AudioValidator) GetAudioInfo(filePath string) (map[string]string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "a:0", 
		"-show_entries", "stream=codec_name,sample_rate,channels", 
		"-show_entries", "format=duration", 
		"-of", "csv=p=0", filePath)
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get audio info: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("unexpected ffprobe output format")
	}

	// Parse stream info (codec, sample_rate, channels)
	streamInfo := strings.Split(lines[0], ",")
	if len(streamInfo) < 3 {
		return nil, fmt.Errorf("incomplete stream information")
	}

	info := map[string]string{
		"codec":       streamInfo[0],
		"sample_rate": streamInfo[1],
		"channels":    streamInfo[2],
		"duration":    lines[1],
	}

	return info, nil
}
