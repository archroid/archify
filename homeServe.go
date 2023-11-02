package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	fs := http.FileServer(http.Dir("/home/ali"))

	http.Handle("/", fs)

	// http.HandleFunc("/latest", redirectHandler)
	// http.HandleFunc("/fringe", fringeHandler)

	http.HandleFunc("/shutdown", handleShutdown)
	http.HandleFunc("/reboot", handleReboot)

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleShutdown(w http.ResponseWriter, r *http.Request) {

	// Execute the shutdown command
	cmd := exec.Command("shutdown", "-h", "now")
	err := cmd.Run()
	if err != nil {
		log.Println("Error shutting down:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// If the command executed successfully, terminate the server
	log.Println("System shutdown triggered")
	os.Exit(0)
}

func handleReboot(w http.ResponseWriter, r *http.Request) {

	// Execute the shutdown command
	cmd := exec.Command("reboot")
	err := cmd.Run()
	if err != nil {
		log.Println("Error rebooting:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// If the command executed successfully, terminate the server
	log.Println("System reboot triggered")
	os.Exit(0)
}

// func redirectHandler(w http.ResponseWriter, r *http.Request) {
// 	dir := "/home/ali/Videos/Fringe"
// 	extension := ".mkv" // Specify the file extension

// 	// Read the directory entries
// 	entries, err := ioutil.ReadDir(dir)
// 	if err != nil {
// 		fmt.Println("Error reading directory:", err)
// 		return
// 	}

// 	// Sort the entries by modification time
// 	sort.Slice(entries, func(i, j int) bool {
// 		return entries[i].ModTime().After(entries[j].ModTime())
// 	})

// 	// Find the last file with the specified extension
// 	var lastFileName string
// 	for _, entry := range entries {
// 		if !entry.IsDir() && strings.HasSuffix(entry.Name(), extension) {
// 			lastFileName = entry.Name()
// 			break
// 		}
// 	}
// 	log.Println("Redirected to " + lastFileName)
// 	http.Redirect(w, r, "/Videos/Fringe/"+lastFileName, http.StatusFound)

// }
// func fringeHandler(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, "/Videos/Fringe/", http.StatusFound)
// }



