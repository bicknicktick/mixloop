package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// AudioSequencer handles sequential audio processing with crossfades
type AudioSequencer struct {
	InputFiles        []string
	OutputFile        string
	CrossfadeDuration float64
	LoopCount         int
	TempDir           string
	Enhance           bool
	DolbyStereo       bool   // Dolby Stereo simulation
	OutputFormat      string // "mp3" or "wav"
	Quality           string // "320k" for mp3, "pcm_s24le" for wav
}

// NewAudioSequencer creates a new audio sequencer with default settings
func NewAudioSequencer(inputFiles []string, outputFile string, crossfadeDuration float64, loopCount int, tempDir string) *AudioSequencer {
	return NewAudioSequencerWithOptions(inputFiles, outputFile, crossfadeDuration, loopCount, tempDir, true, "mp3")
}

// NewAudioSequencerWithOptions creates a new audio sequencer with custom options
func NewAudioSequencerWithOptions(inputFiles []string, outputFile string, crossfadeDuration float64, loopCount int, tempDir string, enhance bool, format string) *AudioSequencer {
	return NewAudioSequencerWithStereoOptions(inputFiles, outputFile, crossfadeDuration, loopCount, tempDir, enhance, false, format)
}

// NewAudioSequencerWithStereoOptions creates a new audio sequencer with stereo options
func NewAudioSequencerWithStereoOptions(inputFiles []string, outputFile string, crossfadeDuration float64, loopCount int, tempDir string, enhance bool, dolbyStereo bool, format string) *AudioSequencer {
	// Determine quality based on format
	quality := "320k"
	if format == "wav" {
		quality = "pcm_s24le"
	}
	
	return &AudioSequencer{
		InputFiles:        inputFiles,
		OutputFile:        outputFile,
		CrossfadeDuration: crossfadeDuration,
		LoopCount:         loopCount,
		TempDir:           tempDir,
		Enhance:           enhance,
		DolbyStereo:       dolbyStereo,
		OutputFormat:      format,
		Quality:           quality,
	}
}

// Process creates a seamless sequence with crossfades and loops
func (as *AudioSequencer) Process() error {
	return as.ProcessWithProgress("", GlobalProgressTracker)
}

// ProcessWithProgress creates a seamless sequence with progress tracking
func (as *AudioSequencer) ProcessWithProgress(sessionID string, tracker *ProgressTracker) error {
	if len(as.InputFiles) == 0 {
		return fmt.Errorf("no input files provided")
	}

	totalSteps := 4
	if as.LoopCount > 1 {
		totalSteps = 5
	}
	if as.Enhance {
		totalSteps++
	}

	// Step 1: Validation
	if tracker != nil && sessionID != "" {
		tracker.UpdateProgress(sessionID, "validation", "Validating audio files...", 0, "", len(as.InputFiles))
	}

	// Step 2: Create sequence of all tracks with crossfades between them
	if tracker != nil && sessionID != "" {
		tracker.UpdateProgress(sessionID, "sequencing", "Creating audio sequence...", 20, "", len(as.InputFiles))
	}
	
	sequenceFile := filepath.Join(as.TempDir, "sequence.mp3")
	err := as.createSequenceWithCrossfades(sequenceFile)
	if err != nil {
		return fmt.Errorf("failed to create sequence: %v", err)
	}
	defer os.Remove(sequenceFile)

	// Step 3: Apply looping with crossfade at boundaries
	var finalFile string
	if as.LoopCount > 1 {
		if tracker != nil && sessionID != "" {
			tracker.UpdateProgress(sessionID, "looping", "Creating looped sequence...", 50, "", as.LoopCount)
		}
		err = as.createLoopedSequence(sequenceFile)
		if err != nil {
			return fmt.Errorf("failed to create looped sequence: %v", err)
		}
		finalFile = filepath.Join(as.TempDir, "looped.mp3")
	} else {
		finalFile = sequenceFile
	}

	// Step 4: Apply enhancement if requested
	if as.Enhance {
		if tracker != nil && sessionID != "" {
			tracker.UpdateProgress(sessionID, "enhancing", "Applying audio enhancement...", 75, "", 0)
		}
		enhancer := NewAudioEnhancer(as.TempDir)
		err = enhancer.ApplyEnhancement(finalFile, as.OutputFile, as.OutputFormat, as.Quality)
		if err != nil {
			return fmt.Errorf("failed to enhance audio: %v", err)
		}
	} else {
		if tracker != nil && sessionID != "" {
			tracker.UpdateProgress(sessionID, "finalizing", "Finalizing output...", 90, "", 0)
		}
		// Copy final file to output
		err = as.copyFile(finalFile, as.OutputFile)
		if err != nil {
			return fmt.Errorf("failed to copy final file: %v", err)
		}
	}

	// Step 5: Complete
	if tracker != nil && sessionID != "" {
		tracker.UpdateProgress(sessionID, "completed", "Audio processing completed!", 100, "", 0)
	}

	return nil
}

// createSequenceWithCrossfades concatenates all input files with crossfades between them
func (as *AudioSequencer) createSequenceWithCrossfades(outputFile string) error {
	if len(as.InputFiles) == 1 {
		// Single file, just copy it
		return as.copyFile(as.InputFiles[0], outputFile)
	}

	if as.CrossfadeDuration <= 0 {
		// No crossfade, use simple concatenation
		return as.concatenateFiles(outputFile)
	}

	// Use crossfade concatenation
	return as.concatenateWithCrossfade(outputFile)
}

// concatenateFiles concatenates files without crossfade using concat demuxer
func (as *AudioSequencer) concatenateFiles(outputFile string) error {
	// Create concat file list
	concatFile := filepath.Join(as.TempDir, "concat_list.txt")
	defer os.Remove(concatFile)

	var concatContent strings.Builder
	for _, file := range as.InputFiles {
		concatContent.WriteString(fmt.Sprintf("file '%s'\n", file))
	}

	err := os.WriteFile(concatFile, []byte(concatContent.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to create concat file: %v", err)
	}

	// Use concat demuxer for perfect concatenation
	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", concatFile,
		"-c", "copy",
		"-y", outputFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg concat error: %v\nOutput: %s", err, output)
	}

	return nil
}

// concatenateWithCrossfade concatenates files with crossfade transitions
func (as *AudioSequencer) concatenateWithCrossfade(outputFile string) error {
	// Start with first file
	currentFile := as.InputFiles[0]
	
	for i := 1; i < len(as.InputFiles); i++ {
		nextFile := as.InputFiles[i]
		tempOutput := filepath.Join(as.TempDir, fmt.Sprintf("temp_concat_%d.mp3", i))
		
		// Crossfade current with next
		cmd := exec.Command("ffmpeg",
			"-i", currentFile,
			"-i", nextFile,
			"-filter_complex",
			fmt.Sprintf("[0][1]acrossfade=d=%.1f:c1=tri:c2=tri", as.CrossfadeDuration),
			"-metadata", "artist=e.bitzy.id",
			"-metadata", "author=e.bitzy.id",
			"-metadata", "composer=e.bitzy.id",
			"-metadata", "comment=Mixed with MixLoop by BITZY.ID",
			"-acodec", "libmp3lame",
			"-y", tempOutput)
		
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("ffmpeg crossfade concat error at step %d: %v\nOutput: %s", i, err, output)
		}
		
		// Clean up previous temp file if it's not an original input
		if i > 1 {
			os.Remove(currentFile)
		}
		
		currentFile = tempOutput
		defer os.Remove(tempOutput)
	}
	
	// Copy final result to output
	return as.copyFile(currentFile, outputFile)
}

// createLoopedSequence takes a sequence and loops it with crossfade at boundaries
func (as *AudioSequencer) createLoopedSequence(sequenceFile string) error {
	if as.LoopCount <= 1 {
		return as.copyFile(sequenceFile, as.OutputFile)
	}

	// Get duration of the sequence
	duration, err := GetAudioDuration(sequenceFile)
	if err != nil {
		return fmt.Errorf("failed to get sequence duration: %v", err)
	}

	// Adjust crossfade if it's too long
	crossfade := as.CrossfadeDuration
	if crossfade > duration/2 {
		crossfade = duration / 2
	}

	if as.LoopCount == 2 {
		// Simple case: two loops with crossfade
		loopedFile := filepath.Join(as.TempDir, "looped.mp3")
		cmd := exec.Command("ffmpeg",
			"-i", sequenceFile,
			"-i", sequenceFile,
			"-filter_complex",
			fmt.Sprintf("[0][1]acrossfade=d=%.1f:c1=tri:c2=tri", crossfade),
			"-metadata", "artist=e.bitzy.id",
			"-metadata", "author=e.bitzy.id",
			"-metadata", "composer=e.bitzy.id",
			"-metadata", "comment=Mixed with MixLoop by BITZY.ID",
			"-acodec", "libmp3lame",
			"-y", loopedFile)
		
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("ffmpeg loop crossfade error: %v\nOutput: %s", err, output)
		}
		return nil
	}

	// Multiple loops: create chain of crossfades
	loopedFile := filepath.Join(as.TempDir, "looped.mp3")
	err = as.createMultipleLoops(sequenceFile, crossfade, loopedFile)
	if err != nil {
		return err
	}
	return as.copyFile(loopedFile, as.OutputFile)
}

// createMultipleLoops handles more than 2 loops with crossfades
func (as *AudioSequencer) createMultipleLoops(sequenceFile string, crossfade float64, outputFile string) error {
	// Create temporary copies for each loop
	var tempFiles []string
	var inputs []string
	
	for i := 0; i < as.LoopCount; i++ {
		tempFile := filepath.Join(as.TempDir, fmt.Sprintf("loop_%d.mp3", i))
		err := as.copyFile(sequenceFile, tempFile)
		if err != nil {
			return fmt.Errorf("failed to create loop copy %d: %v", i, err)
		}
		tempFiles = append(tempFiles, tempFile)
		inputs = append(inputs, "-i", tempFile)
		defer os.Remove(tempFile)
	}

	// Build crossfade chain
	var filterParts []string
	currentLabel := "0"
	
	for i := 1; i < as.LoopCount; i++ {
		if i == as.LoopCount-1 {
			// Last crossfade
			filterParts = append(filterParts, fmt.Sprintf("[%s][%d]acrossfade=d=%.1f:c1=tri:c2=tri", currentLabel, i, crossfade))
		} else {
			nextLabel := fmt.Sprintf("cf%d", i)
			filterParts = append(filterParts, fmt.Sprintf("[%s][%d]acrossfade=d=%.1f:c1=tri:c2=tri[%s]", currentLabel, i, crossfade, nextLabel))
			currentLabel = nextLabel
		}
	}

	filterComplex := strings.Join(filterParts, ";")

	// Execute FFmpeg command with metadata
	args := inputs
	args = append(args, "-filter_complex", filterComplex)
	args = append(args, "-metadata", "artist=e.bitzy.id")
	args = append(args, "-metadata", "author=e.bitzy.id")
	args = append(args, "-metadata", "composer=e.bitzy.id")
	args = append(args, "-metadata", "comment=Mixed with MixLoop by BITZY.ID")
	args = append(args, "-acodec", "libmp3lame", "-y", outputFile)

	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg multiple loops error: %v\nOutput: %s", err, output)
	}

	return nil
}

// copyFile copies a file from src to dst with proper format handling
func (as *AudioSequencer) copyFile(src, dst string) error {
	var args []string
	args = append(args, "-i", src)
	
	// Apply Dolby Stereo simulation if enabled
	if as.DolbyStereo {
		args = append(args, "-af", "stereotools=mlev=1.2")
		args = append(args, "-ac", "2") // Force stereo output
	}
	
	if as.OutputFormat == "wav" {
		args = append(args, "-c:a", as.Quality, "-ar", "48000")
	} else {
		args = append(args, "-c:a", "libmp3lame", "-b:a", as.Quality, "-ar", "48000")
	}
	
	args = append(args, "-y", dst)
	
	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg copy error: %v\nOutput: %s", err, output)
	}
	return nil
}

// GetAudioDuration returns the duration of an audio file in seconds
func GetAudioDuration(file string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		file)
	
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe error: %v", err)
	}

	duration := 0.0
	_, err = fmt.Sscanf(strings.TrimSpace(string(output)), "%f", &duration)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %v", err)
	}
	
	return duration, nil
}
