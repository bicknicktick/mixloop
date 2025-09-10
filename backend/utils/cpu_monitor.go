package utils

import (
	"runtime"
	"sync"
	"time"
)

// CPUMonitor tracks CPU usage and provides throttling capabilities
type CPUMonitor struct {
	mu              sync.RWMutex
	lastCPUTime     time.Duration
	lastSampleTime  time.Time
	currentLoad     float64
	maxAllowedLoad  float64 // Maximum allowed CPU load (0.0-1.0)
}

// NewCPUMonitor creates a new CPU monitor with specified max load
func NewCPUMonitor(maxLoad float64) *CPUMonitor {
	if maxLoad <= 0 || maxLoad > 1.0 {
		maxLoad = 0.7 // Default to 70% max CPU usage
	}
	
	return &CPUMonitor{
		maxAllowedLoad: maxLoad,
		lastSampleTime: time.Now(),
	}
}

// GetCurrentLoad returns the current CPU load (0.0-1.0)
func (cm *CPUMonitor) GetCurrentLoad() float64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.currentLoad
}

// ShouldThrottle returns true if CPU usage is too high
func (cm *CPUMonitor) ShouldThrottle() bool {
	cm.updateCPUUsage()
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.currentLoad > cm.maxAllowedLoad
}

// WaitForCPUCooldown waits until CPU usage drops below threshold
func (cm *CPUMonitor) WaitForCPUCooldown() {
	for cm.ShouldThrottle() {
		time.Sleep(time.Millisecond * 100) // Check every 100ms
	}
}

// GetThrottleDelay returns appropriate delay based on CPU load
func (cm *CPUMonitor) GetThrottleDelay() time.Duration {
	cm.updateCPUUsage()
	cm.mu.RLock()
	load := cm.currentLoad
	cm.mu.RUnlock()
	
	if load < 0.5 {
		return 0 // No delay needed
	} else if load < 0.7 {
		return time.Millisecond * 200 // Light throttling
	} else if load < 0.8 {
		return time.Millisecond * 500 // Medium throttling
	} else {
		return time.Second * 1 // Heavy throttling
	}
}

// updateCPUUsage updates the current CPU usage estimation
func (cm *CPUMonitor) updateCPUUsage() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	now := time.Now()
	if now.Sub(cm.lastSampleTime) < time.Millisecond*500 {
		return // Don't update too frequently
	}
	
	// Simple CPU usage estimation based on goroutines and system load
	numGoroutines := float64(runtime.NumGoroutine())
	numCPU := float64(runtime.NumCPU())
	
	// Estimate load based on goroutine count relative to CPU cores
	estimatedLoad := numGoroutines / (numCPU * 10) // Rough estimation
	if estimatedLoad > 1.0 {
		estimatedLoad = 1.0
	}
	
	// Smooth the load calculation
	if cm.currentLoad == 0 {
		cm.currentLoad = estimatedLoad
	} else {
		cm.currentLoad = (cm.currentLoad*0.7 + estimatedLoad*0.3) // Exponential smoothing
	}
	
	cm.lastSampleTime = now
}
