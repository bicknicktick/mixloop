package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// BatchProcessor handles large-scale audio processing with memory optimization
type BatchProcessor struct {
	MaxConcurrent int
	TempDir       string
	ChunkSize     int
	CPUMonitor    *CPUMonitor
}

// NewBatchProcessor creates a new batch processor optimized for system resources
func NewBatchProcessor(tempDir string) *BatchProcessor {
	// Conservative CPU usage to prevent high load (max 70% CPU usage)
	maxConcurrent := runtime.NumCPU() / 2
	if maxConcurrent < 1 {
		maxConcurrent = 1
	}
	if maxConcurrent > 4 {
		maxConcurrent = 4 // Cap at 4 to prevent CPU overload
	}

	return &BatchProcessor{
		MaxConcurrent: maxConcurrent,
		TempDir:       tempDir,
		ChunkSize:     15, // Smaller chunks to reduce CPU load
		CPUMonitor:    NewCPUMonitor(0.7), // Max 70% CPU usage
	}
}

// ProcessLargeAudioSet processes 100+ audio files efficiently
func (bp *BatchProcessor) ProcessLargeAudioSet(inputFiles []string, outputFile string, loops int, crossfade float64, enhance bool, format, sessionID string) error {
	return bp.ProcessLargeAudioSetWithStereo(inputFiles, outputFile, loops, crossfade, enhance, false, format, sessionID)
}

// ProcessLargeAudioSetWithStereo processes 100+ audio files efficiently with stereo options
func (bp *BatchProcessor) ProcessLargeAudioSetWithStereo(inputFiles []string, outputFile string, loops int, crossfade float64, enhance bool, dolbyStereo bool, format, sessionID string) error {
	if len(inputFiles) <= 20 {
		// Use regular processing for smaller sets
		manager := NewAudioManager(bp.TempDir)
		return manager.ProcessAudioSequenceWithProgress(inputFiles, outputFile, loops, crossfade, enhance, format, sessionID)
	}

	// Update progress
	if sessionID != "" {
		GlobalProgressTracker.UpdateProgress(sessionID, "preparation", "Preparing batch processing...", 5, "", len(inputFiles))
	}

	// Process in chunks to manage memory
	chunks := bp.chunkFiles(inputFiles)
	chunkOutputs := make([]string, len(chunks))
	
	// Process chunks with limited concurrency
	semaphore := make(chan struct{}, bp.MaxConcurrent)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var processingError error

	for i, chunk := range chunks {
		wg.Add(1)
		go func(chunkIndex int, files []string) {
			defer wg.Done()
			semaphore <- struct{}{} // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			// CPU monitoring and throttling
			if chunkIndex > 0 {
				// Wait for CPU to cool down if needed
				bp.CPUMonitor.WaitForCPUCooldown()
				
				// Apply dynamic delay based on CPU load
				delay := bp.CPUMonitor.GetThrottleDelay()
				if delay > 0 {
					time.Sleep(delay)
				}
			}

			// Create chunk-specific temp directory
			chunkDir := filepath.Join(bp.TempDir, fmt.Sprintf("chunk_%d_%s", chunkIndex, sessionID))
			os.MkdirAll(chunkDir, 0755)
			defer os.RemoveAll(chunkDir)

			// Process chunk
			chunkOutput := filepath.Join(chunkDir, fmt.Sprintf("chunk_%d.mp3", chunkIndex))
			manager := NewAudioManager(chunkDir)
			
			// Update progress for this chunk
			if sessionID != "" {
				progress := 10 + float64(chunkIndex)*60/float64(len(chunks))
				GlobalProgressTracker.UpdateProgress(sessionID, "processing", 
					fmt.Sprintf("Processing chunk %d/%d", chunkIndex+1, len(chunks)), 
					progress, "", len(files))
			}

			err := manager.ProcessAudioSequenceWithOptions(files, chunkOutput, 1, crossfade, false, "mp3")
			
			mu.Lock()
			if err != nil && processingError == nil {
				processingError = fmt.Errorf("chunk %d processing failed: %v", chunkIndex, err)
			} else {
				chunkOutputs[chunkIndex] = chunkOutput
			}
			mu.Unlock()
		}(i, chunk)
	}

	wg.Wait()

	if processingError != nil {
		return processingError
	}

	// Update progress
	if sessionID != "" {
		GlobalProgressTracker.UpdateProgress(sessionID, "merging", "Merging processed chunks...", 75, "", len(chunks))
	}

	// Merge all chunks into final output
	err := bp.mergeChunksWithStereo(chunkOutputs, outputFile, loops, enhance, dolbyStereo, format, sessionID)
	if err != nil {
		return fmt.Errorf("failed to merge chunks: %v", err)
	}

	// Final progress update
	if sessionID != "" {
		GlobalProgressTracker.UpdateProgress(sessionID, "completed", "Large audio set processing completed!", 100, "", 0)
	}

	return nil
}

// chunkFiles splits input files into manageable chunks
func (bp *BatchProcessor) chunkFiles(files []string) [][]string {
	var chunks [][]string
	for i := 0; i < len(files); i += bp.ChunkSize {
		end := i + bp.ChunkSize
		if end > len(files) {
			end = len(files)
		}
		chunks = append(chunks, files[i:end])
	}
	return chunks
}

// mergeChunks combines processed chunks into final output
func (bp *BatchProcessor) mergeChunks(chunkFiles []string, outputFile string, loops int, enhance bool, format, sessionID string) error {
	return bp.mergeChunksWithStereo(chunkFiles, outputFile, loops, enhance, false, format, sessionID)
}

// mergeChunksWithStereo combines processed chunks into final output with stereo options
func (bp *BatchProcessor) mergeChunksWithStereo(chunkFiles []string, outputFile string, loops int, enhance bool, dolbyStereo bool, format, sessionID string) error {
	// Filter out empty chunk files
	validChunks := make([]string, 0, len(chunkFiles))
	for _, chunk := range chunkFiles {
		if chunk != "" {
			if _, err := os.Stat(chunk); err == nil {
				validChunks = append(validChunks, chunk)
			}
		}
	}

	if len(validChunks) == 0 {
		return fmt.Errorf("no valid chunks to merge")
	}

	// Create final merge directory
	mergeDir := filepath.Join(bp.TempDir, fmt.Sprintf("merge_%s", sessionID))
	os.MkdirAll(mergeDir, 0755)
	defer os.RemoveAll(mergeDir)

	// Use AudioManager to merge chunks with stereo options
	manager := NewAudioManager(mergeDir)
	return manager.ProcessAudioSequenceWithProgressAndStereo(validChunks, outputFile, loops, 2.0, enhance, dolbyStereo, format, sessionID)
}

// OptimizeForLargeFiles adjusts settings for processing many files
func (bp *BatchProcessor) OptimizeForLargeFiles(fileCount int) {
	if fileCount > 100 {
		bp.ChunkSize = 10 // Much smaller chunks for very large sets
		bp.MaxConcurrent = min(bp.MaxConcurrent, 2) // Very conservative concurrency
	} else if fileCount > 50 {
		bp.ChunkSize = 12
		bp.MaxConcurrent = min(bp.MaxConcurrent, 3) // Reduced concurrency
	} else if fileCount > 20 {
		bp.ChunkSize = 15
		bp.MaxConcurrent = min(bp.MaxConcurrent, 4)
	}
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
