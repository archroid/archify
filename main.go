package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/charmbracelet/log"
	"github.com/common-nighthawk/go-figure"
	"github.com/joho/godotenv"
)

func main() {

	//ASCII art on startup
	myFigure := figure.NewColorFigure("Home Serve", "", "blue", false)
	myFigure.Print()

	fmt.Println("")

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}

	// HTTP routes and serves

	fs := http.FileServer(http.Dir("/home/ali"))

	http.Handle("/", fs)

	http.HandleFunc("/shutdown", handleShutdown)
	http.HandleFunc("/off", handleShutdown)

	http.HandleFunc("/reboot", handleReboot)
	http.HandleFunc("/restart", handleReboot)

	http.HandleFunc("/sleep", handleSleep)
	http.HandleFunc("/suspend", handleSleep)

	// run bots and server in goroutines

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

	// Wait here until CTRL-C or other term signal is received.

	log.Info("Press CTRL-C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

// HTTP hanlders

func handleSleep(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Suspention triggered"))

	err := sleep()
	if err != nil {
		http.Error(w, "Error suspending "+err.Error(), http.StatusInternalServerError)
	}
}

func handleShutdown(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Shutting down triggered"))

	err := shutdown()
	if err != nil {
		http.Error(w, "Error shutting down "+err.Error(), http.StatusInternalServerError)
	}
}

func handleReboot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Rebooting triggered"))

	err := reboot()
	if err != nil {
		http.Error(w, "Error rebooting "+err.Error(), http.StatusInternalServerError)
	}
}
