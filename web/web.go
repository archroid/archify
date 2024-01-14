package web

import (
	"net/http"

	"github.com/charmbracelet/log"
)

func Runweb() error {
	server2 := http.NewServeMux()

	server2.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	})
	server2.Handle("/", server2)
	http.ListenAndServe("80", server2)
	log.Info("web server running on port 8081")
	return nil
}
