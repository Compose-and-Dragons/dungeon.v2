package ui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ThinkingController manages the thinking animation
type ThinkingController struct {
	stopChan  chan bool
	pauseChan chan bool
	doneChan  chan bool
	stopped   bool
	paused    bool
	message   string
	color     string
	mutex     sync.RWMutex
}

// NewThinkingController creates a new thinking animation controller with initialized channels
func NewThinkingController() *ThinkingController {
	return &ThinkingController{
		stopChan:  make(chan bool),
		pauseChan: make(chan bool),
		doneChan:  make(chan bool),
	}
}

// Start begins the thinking animation with the specified color and message
func (tc *ThinkingController) Start(color string, message string) {
	tc.mutex.Lock()
	tc.message = message
	tc.color = color
	tc.stopped = false
	tc.paused = false
	tc.mutex.Unlock()
	
	go func() {
		defer close(tc.doneChan)

		animationChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		index := 0

		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))

	resumeAnimation:
		for {
			select {
			case <-tc.stopChan:
				// Clear the current line
				tc.mutex.RLock()
				currentMessage := tc.message
				tc.mutex.RUnlock()
				fmt.Print("\r" + strings.Repeat(" ", len(currentMessage)+5) + "\r")
				return
			case <-tc.pauseChan:
				// Animation is paused, wait for resume or stop
				tc.mutex.Lock()
				tc.paused = true
				tc.mutex.Unlock()
				for {
					select {
					case <-tc.stopChan:
						tc.mutex.RLock()
						currentMessage := tc.message
						tc.mutex.RUnlock()
						fmt.Print("\r" + strings.Repeat(" ", len(currentMessage)+5) + "\r")
						return
					case <-tc.pauseChan:
						// Resume animation
						tc.mutex.Lock()
						tc.paused = false
						tc.mutex.Unlock()
						goto resumeAnimation
					}
				}
			case <-ticker.C:
				tc.mutex.RLock()
				isPaused := tc.paused
				currentMessage := tc.message
				tc.mutex.RUnlock()
				
				if !isPaused {
					// Clear current line and print new animation frame
					animatedMessage := fmt.Sprintf("\r%s %s", animationChars[index], currentMessage)
					fmt.Print(textStyle.Render(animatedMessage))
					index = (index + 1) % len(animationChars)
				}
			}
		}
	}()
}

// UpdateMessage safely updates the message displayed in the thinking animation
func (tc *ThinkingController) UpdateMessage(message string) {
	tc.mutex.Lock()
	tc.message = message
	tc.mutex.Unlock()
}

// Pause pauses the thinking animation if it's currently running
func (tc *ThinkingController) Pause() {
	tc.mutex.RLock()
	if !tc.stopped && !tc.paused {
		tc.mutex.RUnlock()
		select {
		case tc.pauseChan <- true:
		default:
		}
	} else {
		tc.mutex.RUnlock()
	}
}

// Resume resumes the thinking animation if it's currently paused
func (tc *ThinkingController) Resume() {
	tc.mutex.RLock()
	if !tc.stopped && tc.paused {
		tc.mutex.RUnlock()
		select {
		case tc.pauseChan <- true:
		default:
		}
	} else {
		tc.mutex.RUnlock()
	}
}

// IsPaused returns true if the thinking animation is currently paused
func (tc *ThinkingController) IsPaused() bool {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()
	return tc.paused
}

// Stop stops the thinking animation, clears the line, and waits for the goroutine to finish
func (tc *ThinkingController) Stop() {
	if !tc.stopped {
		tc.stopped = true
		close(tc.stopChan)
		<-tc.doneChan
	}
}

// IsStarted returns true if the thinking animation is currently running
func (tc *ThinkingController) IsStarted() bool {
	select {
	case <-tc.doneChan:
		return false
	default:
		return true
	}
}