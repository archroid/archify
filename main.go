package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	log "github.com/charmbracelet/log"
	"github.com/common-nighthawk/go-figure"
	"github.com/joho/godotenv"
)

func main() {

	//ASCII art
	myFigure := figure.NewColorFigure("Home Serve", "", "blue", false)
	myFigure.Print()

	fmt.Println("\n")

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}

	fs := http.FileServer(http.Dir("/home/ali"))

	http.Handle("/", fs)

	http.HandleFunc("/shutdown", handleShutdown)
	http.HandleFunc("/off", handleShutdown)

	http.HandleFunc("/reboot", handleReboot)
	http.HandleFunc("/restart", handleReboot)

	http.HandleFunc("/sleep", handleSleep)
	http.HandleFunc("/suspend", handleSleep)

	go func() {
		log.Info("Server started on http://localhost:8080")
		log.Error(http.ListenAndServe(":8080", nil))

	}()

	go func() {
		err := discordBot()

		if err != nil {
			log.Error("Discord bot failed to start:  ", err)
			return
		}
	}()

	go func() {
		err := telegramBot()

		if err != nil {
			log.Error("Telegram bot failed to start:  ", err)
			return
		}

	}()

	log.Info("Press CTRL-C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

func handleSleep(w http.ResponseWriter, r *http.Request) {
	// Execute the sleep command
	cmd := exec.Command("systemctl", "suspend")
	err := cmd.Run()
	if err != nil {
		log.Error("Error sleeping: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Info("System suspend triggered")
}

func handleShutdown(w http.ResponseWriter, r *http.Request) {
	// Execute the shutdown command
	cmd := exec.Command("shutdown", "-h", "now")
	err := cmd.Run()
	if err != nil {
		log.Error("Error shutting down: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// If the command executed successfully, terminate the server
	log.Info("System shutdown triggered")
	os.Exit(0)
}

func handleReboot(w http.ResponseWriter, r *http.Request) {
	// Execute the reboot command
	cmd := exec.Command("reboot")
	err := cmd.Run()
	if err != nil {
		log.Error("Error rebooting: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// If the command executed successfully, terminate the server
	log.Info("System reboot triggered")
	os.Exit(0)
}
