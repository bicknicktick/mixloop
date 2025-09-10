package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ProgressUpdate represents a progress update message
type ProgressUpdate struct {
	SessionID   string  `json:"session_id"`
	Stage       string  `json:"stage"`
	Progress    float64 `json:"progress"`
	Message     string  `json:"message"`
	CurrentFile string  `json:"current_file,omitempty"`
	TotalFiles  int     `json:"total_files,omitempty"`
	Timestamp   int64   `json:"timestamp"`
}

// ProgressTracker manages progress updates for audio processing
type ProgressTracker struct {
	connections map[string]*websocket.Conn
	updates     map[string]*ProgressUpdate
	mutex       sync.RWMutex
	upgrader    websocket.Upgrader
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker() *ProgressTracker {
	return &ProgressTracker{
		connections: make(map[string]*websocket.Conn),
		updates:     make(map[string]*ProgressUpdate),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
	}
}

// Global progress tracker instance
var GlobalProgressTracker = NewProgressTracker()

// HandleWebSocket handles WebSocket connections for progress updates
func (pt *ProgressTracker) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := pt.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		log.Println("No session ID provided")
		return
	}

	pt.mutex.Lock()
	pt.connections[sessionID] = conn
	pt.mutex.Unlock()

	// Send existing progress if available
	pt.mutex.RLock()
	if update, exists := pt.updates[sessionID]; exists {
		conn.WriteJSON(update)
	}
	pt.mutex.RUnlock()

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	// Clean up connection
	pt.mutex.Lock()
	delete(pt.connections, sessionID)
	pt.mutex.Unlock()
}

// UpdateProgress sends a progress update
func (pt *ProgressTracker) UpdateProgress(sessionID, stage, message string, progress float64, currentFile string, totalFiles int) {
	update := &ProgressUpdate{
		SessionID:   sessionID,
		Stage:       stage,
		Progress:    progress,
		Message:     message,
		CurrentFile: currentFile,
		TotalFiles:  totalFiles,
		Timestamp:   time.Now().UnixMilli(),
	}

	pt.mutex.Lock()
	pt.updates[sessionID] = update
	conn, exists := pt.connections[sessionID]
	pt.mutex.Unlock()

	if exists && conn != nil {
		if err := conn.WriteJSON(update); err != nil {
			log.Printf("Error sending progress update: %v", err)
		}
	}
}

// GetProgress returns the current progress for a session
func (pt *ProgressTracker) GetProgress(sessionID string) (*ProgressUpdate, bool) {
	pt.mutex.RLock()
	defer pt.mutex.RUnlock()
	update, exists := pt.updates[sessionID]
	return update, exists
}

// CleanupSession removes progress data for a completed session
func (pt *ProgressTracker) CleanupSession(sessionID string) {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()
	delete(pt.updates, sessionID)
	if conn, exists := pt.connections[sessionID]; exists {
		conn.Close()
		delete(pt.connections, sessionID)
	}
}

// ProgressHandler handles HTTP requests for progress updates
func ProgressHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	update, exists := GlobalProgressTracker.GetProgress(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(update)
}

// WebSocketHandler handles WebSocket connections
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	GlobalProgressTracker.HandleWebSocket(w, r)
}
