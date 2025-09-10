package utils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// AudioEnhancer handles audio enhancement filters
type AudioEnhancer struct {
	TempDir string
}

// NewAudioEnhancer creates a new audio enhancer
func NewAudioEnhancer(tempDir string) *AudioEnhancer {
	return &AudioEnhancer{
		TempDir: tempDir,
	}
}

// ApplyEnhancement applies audio enhancement filters to improve quality
func (ae *AudioEnhancer) ApplyEnhancement(inputFile, outputFile string, outputFormat, quality string) error {
	// Build enhancement filter chain
	filterChain := ae.buildEnhancementFilters()
	
	// Build FFmpeg command based on output format
	var args []string
	args = append(args, "-i", inputFile)
	args = append(args, "-af", filterChain)
	
	// Add metadata
	args = append(args, "-metadata", "artist=e.bitzy.id")
	args = append(args, "-metadata", "author=e.bitzy.id")
	args = append(args, "-metadata", "composer=e.bitzy.id")
	args = append(args, "-metadata", "comment=Mixed with MixLoop by BITZY.ID")
	
	if outputFormat == "wav" {
		args = append(args, "-c:a", quality, "-ar", "48000")
	} else {
		// MP3 format
		args = append(args, "-c:a", "libmp3lame", "-b:a", quality, "-ar", "48000")
	}
	
	args = append(args, "-y", outputFile)
	
	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg enhancement error: %v\nOutput: %s", err, output)
	}
	
	return nil
}

// buildEnhancementFilters creates the audio enhancement filter chain
func (ae *AudioEnhancer) buildEnhancementFilters() string {
	filters := []string{
		"highpass=f=80",                           // Remove low-frequency hum
		"lowpass=f=16000",                         // Remove ultrasonic noise
		"acompressor=threshold=-20dB:ratio=3",     // Stabilize dynamics
		"loudnorm=I=-14:TP=-2:LRA=11",            // Normalize loudness
	}
	
	return strings.Join(filters, ",")
}

// ApplyEnhancementToFile is a convenience method for single file enhancement
func (ae *AudioEnhancer) ApplyEnhancementToFile(inputFile string, outputFormat, quality string) (string, error) {
	// Generate output filename
	baseName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	var outputFile string
	
	if outputFormat == "wav" {
		outputFile = filepath.Join(ae.TempDir, baseName+"_enhanced.wav")
	} else {
		outputFile = filepath.Join(ae.TempDir, baseName+"_enhanced.mp3")
	}
	
	err := ae.ApplyEnhancement(inputFile, outputFile, outputFormat, quality)
	if err != nil {
		return "", err
	}
	
	return outputFile, nil
}
