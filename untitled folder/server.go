package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

const port = 8000

// OverlayState represents the shared state for the overlay
type OverlayState struct {
	Label             string  `json:"label"`
	Likes             int     `json:"likes"`
	Goal              int     `json:"goal"`
	RotationDuration  float64 `json:"rotationDuration"`
	RotationSpeed     float64 `json:"rotationSpeed"`
}

var (
	state = OverlayState{
		Label:            "LIKE GOAL",
		Likes:            0,
		Goal:             1000000,
		RotationDuration: 4.0,
		RotationSpeed:    1.0,
	}
	stateMutex sync.RWMutex
)

func main() {
	// Print welcome banner
	printWelcomeBanner()

	// Setup routes
	http.HandleFunc("/api/state", handleState)
	http.Handle("/", http.FileServer(http.Dir(".")))

	// Start server in a goroutine
	serverAddr := fmt.Sprintf(":%d", port)
	go func() {
		log.Printf("Server starting on port %d...\n", port)
		if err := http.ListenAndServe(serverAddr, nil); err != nil {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Wait a moment for server to start
	time.Sleep(500 * time.Millisecond)

	// Print instructions
	printInstructions()

	// Try to open browser (optional, doesn't fail if it can't)
	go func() {
		time.Sleep(1 * time.Second)
		openBrowser(fmt.Sprintf("http://localhost:%d/control.html", port))
	}()

	// Keep server running
	select {}
}

func printWelcomeBanner() {
	banner := `

                                                               
$$\      $$\ $$$$$$$$\ $$\       $$$$$$\   $$$$$$\  $$\      $$\ $$$$$$$$\                                            
$$ | $\  $$ |$$  _____|$$ |     $$  __$$\ $$  __$$\ $$$\    $$$ |$$  _____|                                           
$$ |$$$\ $$ |$$ |      $$ |     $$ /  \__|$$ /  $$ |$$$$\  $$$$ |$$ |                                                 
$$ $$ $$\$$ |$$$$$\    $$ |     $$ |      $$ |  $$ |$$\$$\$$ $$ |$$$$$\                                               
$$$$  _$$$$ |$$  __|   $$ |     $$ |      $$ |  $$ |$$ \$$$  $$ |$$  __|                                              
$$$  / \$$$ |$$ |      $$ |     $$ |  $$\ $$ |  $$ |$$ |\$  /$$ |$$ |                                                 
$$  /   \$$ |$$$$$$$$\ $$$$$$$$\\$$$$$$  | $$$$$$  |$$ | \_/ $$ |$$$$$$$$\                                            
\__/     \__|\________|\________|\______/  \______/ \__|     \__|\________|                                           
                                                                                                                      
                                                                                                                      
                                                                                                                      
$$$$$$$\   $$$$$$\  $$\   $$\ $$$$$$$$\ $$\   $$\ $$\      $$\   $$\  $$$$$$\  $$$$$$\ $$\    $$\ $$$$$$$$\       $$\ 
$$  __$$\ $$  __$$\ $$ | $$  |$$  _____|$$ |  $$ |$$ |     $$ |  $$ |$$  __$$\ \_$$  _|$$ |   $$ |$$  _____|      $$ |
$$ |  $$ |$$ /  $$ |$$ |$$  / $$ |      \$$\ $$  |$$ |     $$ |  $$ |$$ /  \__|  $$ |  $$ |   $$ |$$ |            $$ |
$$$$$$$  |$$ |  $$ |$$$$$  /  $$$$$\     \$$$$  / $$ |     $$ |  $$ |\$$$$$$\    $$ |  \$$\  $$  |$$$$$\          $$ |
$$  ____/ $$ |  $$ |$$  $$<   $$  __|    $$  $$<  $$ |     $$ |  $$ | \____$$\   $$ |   \$$\$$  / $$  __|         \__|
$$ |      $$ |  $$ |$$ |\$$\  $$ |      $$  /\$$\ $$ |     $$ |  $$ |$$\   $$ |  $$ |    \$$$  /  $$ |                
$$ |       $$$$$$  |$$ | \$$\ $$$$$$$$\ $$ /  $$ |$$$$$$$$\\$$$$$$  |\$$$$$$  |$$$$$$\    \$  /   $$$$$$$$\       $$\ 
\__|       \______/ \__|  \__|\________|\__|  \__|\________|\______/  \______/ \______|    \_/    \________|      \__|
                                                                                                                      
                                                                                                                      
                                                                                                                      
                                                                                                                
                                                                                                                
                                                                                                                
`
	fmt.Print(banner)
	fmt.Println()
}

func printInstructions() {
	overlayURL := fmt.Sprintf("http://localhost:%d/index.html", port)
	controlURL := fmt.Sprintf("http://localhost:%d/control.html", port)

	// ANSI escape codes for bold and red
	boldRed := "\033[1;31m"
	reset := "\033[0m"

	fmt.Println("========================================")
	fmt.Printf("%sOverlay Server is running ✅%s\n\n", boldRed, reset)
	fmt.Printf("%sOverlay URL (paste into OBS Browser Source):%s\n", boldRed, reset)
	fmt.Printf("%s%s%s\n\n", boldRed, overlayURL, reset)
	fmt.Printf("%sControl Panel (open in your browser):%s\n", boldRed, reset)
	fmt.Printf("%s%s%s\n\n", boldRed, controlURL, reset)
	fmt.Printf("%s1. Open OBS%s\n", boldRed, reset)
	fmt.Printf("%s2. Add a Browser Source%s\n", boldRed, reset)
	fmt.Printf("%s3. Paste the Overlay URL%s\n", boldRed, reset)
	fmt.Printf("%s4. Open the Control Panel link in your browser%s\n\n", boldRed, reset)
	fmt.Printf("%sDo NOT close this window while streaming.%s\n", boldRed, reset)
	fmt.Println("========================================")
}

func handleState(w http.ResponseWriter, r *http.Request) {
	// Set headers to prevent caching
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")

	switch r.Method {
	case http.MethodGet:
		handleGetState(w, r)
	case http.MethodPost:
		handlePostState(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetState(w http.ResponseWriter, r *http.Request) {
	stateMutex.RLock()
	defer stateMutex.RUnlock()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := json.NewEncoder(w).Encode(state); err != nil {
		http.Error(w, "Failed to encode state", http.StatusInternalServerError)
		return
	}
}

func handlePostState(w http.ResponseWriter, r *http.Request) {
	var newState OverlayState
	if err := json.NewDecoder(r.Body).Decode(&newState); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	stateMutex.Lock()
	state = newState
	stateMutex.Unlock()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return
	}
	cmd.Stderr = os.Stderr
	_ = cmd.Run() // Ignore errors
}
