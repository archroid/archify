// in the name of ALLAH
package main

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/template"

	discordbot "archroid/archify/discordbot"
	"archroid/archify/telegrambot"
	utils "archroid/archify/utils"

	log "github.com/charmbracelet/log"
	"github.com/common-nighthawk/go-figure"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var localip string
var homePath string
var showHiddenFiles = false

func main() {
	//ASCII art on startup
	myFigure := figure.NewColorFigure("ARCHIFY", "", "blue", true)
	myFigure.Print()

	// save logs into a file
	os.Remove("archify.log")
	f, err := os.OpenFile("archify.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	// Load the .env file
	err = godotenv.Load()
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

	// HTTP routes and serve

	// main http server
	r := mux.NewRouter()

	r.HandleFunc("/", handleSite)
	r.PathPrefix("/web").Handler(http.StripPrefix("/web", http.FileServer(http.Dir("./web"))))

	r.PathPrefix("/d/").Handler(http.StripPrefix("/d/", http.FileServer(http.Dir(homePath))))

	r.HandleFunc("/log", handleLog)

	r.HandleFunc("/ping", handlePing)

	r.HandleFunc("/dir/{dir:.*}", handleDirectory)
	r.HandleFunc("/hiddinfiles/{bool}", handleHiddenFiles)
	r.HandleFunc("/upload", handleUpload).Methods("POST")

	r.HandleFunc("/shutdown", handleShutdown)
	r.HandleFunc("/off", handleShutdown)

	r.HandleFunc("/reboot", handleReboot)
	r.HandleFunc("/restart", handleReboot)

	r.HandleFunc("/sleep", handleSleep)
	r.HandleFunc("/suspend", handleSleep)

	// run bots and server in goroutines
	go func() {
		log.Info("Server started: http://" + localip + ":8080")
		log.Error(http.ListenAndServe(":8080", r))

	}()

	go func() {
		err := telegrambot.Run()

		if err != nil {
			log.Error("Telegram bot failed to start:  ", err)
			return
		}
	}()

	go func() {
		err = discordbot.RunSession()
		if err != nil {
			log.Error("Discord bot failed to start:  ", err)
		}
	}()

	log.Info("Press CTRL-C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

// HTTP hanlders

func handleLog(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./archify.log")
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"resp": "pong"})
}

func handleSite(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/index.html")
}

func handleSleep(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"resp": "Sleeping"})

	err := utils.Sleep()
	if err != nil {
		http.Error(w, "Error suspending "+err.Error(), http.StatusInternalServerError)
	}
}

func handleShutdown(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"resp": "Shutting down"})

	err := utils.Shutdown()
	if err != nil {
		http.Error(w, "Error shutting down "+err.Error(), http.StatusInternalServerError)
	}
}

func handleReboot(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"resp": "Rebooting"})

	err := utils.Reboot()
	if err != nil {
		http.Error(w, "Error rebooting "+err.Error(), http.StatusInternalServerError)
	}
}

// structs requiered for dir template
type File struct {
	IsDir         bool
	Name          string
	FileURL       string
	FileURLFolder string
	CurrentDir    string
}

type Directory struct {
	CurrentDir string
	Files      []File
}

func handleDirectory(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	directory := vars["dir"]

	directory = homePath + "/" + directory

	files, err := os.ReadDir(directory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	var filesData []File

	directory = strings.TrimPrefix(directory, homePath)

	if directory == "/" {
		directory = strings.TrimSuffix(directory, "/")

	}
	if !showHiddenFiles {
		for _, file := range files {
			if !strings.HasPrefix(file.Name(), ".") {
				filesData = append(filesData, File{
					IsDir:         file.IsDir(),
					Name:          file.Name(),
					FileURL:       "http://" + localip + ":8080/d/" + directory + file.Name() + "/",
					FileURLFolder: "http://" + localip + ":8080/dir/" + directory + file.Name() + "/",
				})
			}
		}
	} else {
		for _, file := range files {
			filesData = append(filesData, File{
				IsDir:         file.IsDir(),
				Name:          file.Name(),
				FileURL:       "http://" + localip + ":8080/d/" + directory + file.Name() + "/",
				FileURLFolder: "http://" + localip + ":8080/dir/" + directory + file.Name() + "/",
			})

		}
	}

	if directory == "" {
		directory = "/"

	}
	dir := Directory{
		CurrentDir: directory,
		Files:      filesData,
	}

	// Parse and execute the template
	t, err := template.New("directory.html").Funcs(template.FuncMap{
		"hasSuffix": strings.HasSuffix,
	}).ParseFiles("web/directory.html")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, dir)
	if err != nil {
		log.Fatal(err)
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]

	log.Warnf("Dir Access: %s", ip)

}

func handleHiddenFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bool := vars["bool"]

	if bool == "true" {
		showHiddenFiles = true
	} else {
		showHiddenFiles = false
	}

	http.Redirect(w, r, "/dir/", http.StatusSeeOther)

}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	log.Warn("Uploading something")
	// Parse the multipart form in the request
	err := r.ParseMultipartForm(10 << 20) // limit your maxMultipartMemory to 10MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	dirPath := homePath + "/uploads/"
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	// Retrieve the file from form data
	file, handler, err := r.FormFile("myFile") // "myFile" is the key of the input form
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	defer file.Close()

	// Create a new file in the server's file system with the same name
	// Assuming you have a folder named 'uploads'
	dst, err := os.Create(homePath + "/uploads/" + handler.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "done"})

}
