package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Devices struct {
	List []string
}

type ServerResponse struct {
	Devices Devices
}

func main() {
	http.HandleFunc("/", home)
	// static routes
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// api
	http.HandleFunc("/deviceCount", deviceCount)

	fmt.Println("Server Running...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

// api
func home(w http.ResponseWriter, r *http.Request) {
	devices := getDevices()
	response := ServerResponse{devices}

	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func deviceCount(w http.ResponseWriter, r *http.Request) {
	resp := getDevices()
	sendJson(w, r, resp)
}

// helpers
func execCmd(cmdName string, cmdArgs []string) string {
	var (
		cmdOut []byte
		err    error
	)
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running exec command: ", err)
		os.Exit(1)
	}
	return string(cmdOut)
}

func sendJson(w http.ResponseWriter, r *http.Request, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getDevices() Devices {
	devices := execCmd("arp", []string{"-a"})
	deviceList := strings.Split(devices, "\n")
	return Devices{deviceList}
}
