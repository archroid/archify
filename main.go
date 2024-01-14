package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/charmbracelet/log"
	"github.com/common-nighthawk/go-figure"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var localip string
var homePath string

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

	homePath = os.Getenv("HOME_PATH")

	//find device local ip
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			localip = ipv4.String()
		}
	}

	// HTTP routes and serves

	// serve files in home directory
	r2 := mux.NewRouter()
	fs := http.FileServer(http.Dir(homePath))
	r2.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	// main http server
	r := mux.NewRouter()

	r.HandleFunc("/", handleSite)

	r.HandleFunc("/ping", handlePing)

	r.HandleFunc("/f/{dir:.*}", handleFolder)

	r.HandleFunc("/shutdown", handleShutdown)
	r.HandleFunc("/off", handleShutdown)

	r.HandleFunc("/reboot", handleReboot)
	r.HandleFunc("/restart", handleReboot)

	r.HandleFunc("/sleep", handleSleep)
	r.HandleFunc("/suspend", handleSleep)

	// run bots and server in goroutines
	go func() {
		log.Error(http.ListenAndServe(":8090", r2))
	}()

	go func() {
		log.Info("Server started on http://" + localip + ":8080")
		log.Error(http.ListenAndServe(":8080", r))

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

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong"))

}

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

func handleFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	directory := vars["dir"]

	directory = homePath + "/" + directory

	getDirectory(w, directory)
}

func handleSite(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/index.html")
}

func getDirectory(w http.ResponseWriter, directory string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the CSS styles

	fmt.Fprintln(w, `
	
	<head>
		<link rel='icon' href='https://avatars.githubusercontent.com/u/50708771?s=400&u=283e9b4589fc6d1455f2cea0356cc7f4156a5251&v=4'> <title>Home Serv</title>
		<style>
			table {
			width: 100%;
			}

			td {
			word-wrap: break-word;
			}

			@media only screen and (max-width: 600px) {
			td {
			font-size: small;
			}
			}
		</style>

	</head>
		
		`)

	fmt.Fprintln(w, `


		<style>
			@import url('https://fonts.googleapis.com/css2?family=Roboto&display=swap');
			body { font-family: 'Roboto', sans-serif; font-size: 16px; }
			a { color: black; text-decoration: none; }
			a:active { background-color: lightgray; }
			table { border-collapse: collapse; width: auto; border: 0; }
			td { border: 0; padding: 5px; font-size: 20px; }
			h4 { color: blue; 
		</style>
		<h1>Home Serve</h1>


		`)

	fmt.Fprintln(w, "<h4>"+directory+"</h4>")

	// Start the table
	fmt.Fprintln(w, "<table>")

	directory = strings.TrimPrefix(directory, homePath)

	if directory == "/" {
		directory = strings.TrimSuffix(directory, "/")

	}

	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			fileURL := "http://" + localip + ":8090/" + directory + file.Name() + "/"
			fileURLFolder := "http://" + localip + ":8080/f/" + directory + file.Name() + "/"
			// Check if the file is a directory
			if file.IsDir() {
				// If it is, prepend a document emoji to the file name
				fmt.Fprintf(w, "<tr><td>📁 <a href=\"%s\">%s</a></td></tr>", fileURLFolder, file.Name())
			} else {
				// Check if the file is a video
				if strings.HasSuffix(file.Name(), ".mp4") || strings.HasSuffix(file.Name(), ".avi") || strings.HasSuffix(file.Name(), ".mov") || strings.HasSuffix(file.Name(), ".mkv") || strings.HasSuffix(file.Name(), ".flv") || strings.HasSuffix(file.Name(), ".wmv") || strings.HasSuffix(file.Name(), ".webm") {
					// If it is, prepend a video camera emoji to the file name
					fmt.Fprintf(w, "<tr><td>🎥 <a href=\"%s\">%s</a></td></tr>", fileURL, file.Name())
				} else if strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".jpeg") || strings.HasSuffix(file.Name(), ".png") || strings.HasSuffix(file.Name(), ".gif") || strings.HasSuffix(file.Name(), ".bmp") {
					// If it is an image, prepend an image emoji to the file name
					fmt.Fprintf(w, "<tr><td>📸 <a href=\"%s\">%s</a></td></tr>", fileURL, file.Name())
				} else {
					fmt.Fprintf(w, "<tr><td>📄 <a href=\"%s\">%s</a></td></tr>", fileURL, file.Name())
				}
			}
		}
	}

	// End the table
	fmt.Fprintln(w, "</table>")
}
