package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"mixloop/utils"
)

func MixAudioHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse form parameters
	loopsStr := r.FormValue("loops")
	crossfadeStr := r.FormValue("crossfade")
	enhanceStr := r.FormValue("enhance")
	dolbyStereoStr := r.FormValue("dolby_stereo")
	formatStr := r.FormValue("format")

	loops := 1
	if loopsStr != "" {
		if parsedLoops, err := strconv.Atoi(loopsStr); err == nil && parsedLoops > 0 {
			loops = parsedLoops
		}
	}

	crossfade := 2.0
	if crossfadeStr != "" {
		if parsedCrossfade, err := strconv.ParseFloat(crossfadeStr, 64); err == nil && parsedCrossfade >= 0 {
			crossfade = parsedCrossfade
		}
	}

	enhance := true // Default to true
	if enhanceStr != "" {
		enhance = enhanceStr == "true"
	}

	dolbyStereo := false // Default to false
	if dolbyStereoStr != "" {
		dolbyStereo = dolbyStereoStr == "true"
	}

	format := "mp3" // Default format
	if formatStr == "wav" {
		format = "wav"
	}

	// Get uploaded files
	files := r.MultipartForm.File["audio_files"]
	if len(files) == 0 {
		http.Error(w, "No audio files provided", http.StatusBadRequest)
		return
	}

	// Create unique session directory
	sessionID := fmt.Sprintf("%d", time.Now().UnixNano())
	sessionDir := filepath.Join("uploads", sessionID)
	os.MkdirAll(sessionDir, 0755)
	defer os.RemoveAll(sessionDir) // Clean up after processing

	// Save uploaded files
	var savedFiles []string
	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to open uploaded file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Save file to disk
		filename := fmt.Sprintf("input_%d%s", i, filepath.Ext(fileHeader.Filename))
		filePath := filepath.Join(sessionDir, filename)
		
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			return
		}

		savedFiles = append(savedFiles, filePath)
	}

	// Generate output filename with proper extension
	var outputFile string
	if format == "wav" {
		outputFile = filepath.Join("output", fmt.Sprintf("mix_%s.wav", sessionID))
	} else {
		outputFile = filepath.Join("output", fmt.Sprintf("mix_%s.mp3", sessionID))
	}

	// Use existing session ID for progress tracking
	
	// Choose processing method based on file count
	if len(savedFiles) > 20 {
		// Use batch processor for large sets
		batchProcessor := utils.NewBatchProcessor(sessionDir)
		batchProcessor.OptimizeForLargeFiles(len(savedFiles))
		err = batchProcessor.ProcessLargeAudioSetWithStereo(savedFiles, outputFile, loops, crossfade, enhance, dolbyStereo, format, sessionID)
	} else {
		// Use regular processing for smaller sets
		manager := utils.NewAudioManager(sessionDir)
		err = manager.ProcessAudioSequenceWithProgressAndStereo(savedFiles, outputFile, loops, crossfade, enhance, dolbyStereo, format, sessionID)
	}
	
	if err != nil {
		fmt.Printf("Audio processing error: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to process audio: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.Remove(outputFile) // Clean up output file after sending

	// Send result file with proper headers
	if format == "wav" {
		w.Header().Set("Content-Type", "audio/wav")
		w.Header().Set("Content-Disposition", "attachment; filename=mixloop_output.wav")
	} else {
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Header().Set("Content-Disposition", "attachment; filename=mixloop_output.mp3")
	}

	outputData, err := os.ReadFile(outputFile)
	if err != nil {
		http.Error(w, "Failed to read output file", http.StatusInternalServerError)
		return
	}

	w.Write(outputData)
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
