package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os/exec"
	"sync"
	"text/template"
	"time"

	"github.com/skip2/go-qrcode"
)

func main() {
	var wg sync.WaitGroup

	port := ":32412"
	srv := &http.Server{Addr: port}

	wg.Add(2)
	go func() {
		defer wg.Done()
		open("http://localhost" + port)
	}()

	go func() {
		defer wg.Done()
		shutdown(srv)
	}()

	http.HandleFunc("/", showQr)
	srv.ListenAndServe()
	wg.Wait()
}

func showQr(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("layout.html"))
	out, err := readClipboard()
	if err != nil {
		fmt.Fprintln(w, err.Error())
		fmt.Fprintf(w, "You may need wl-clipboard\n")
	} else {
		png, _ := qrcode.Encode(string(out), qrcode.Low, 512)
		encoded := base64.StdEncoding.EncodeToString(png)
		tmpl.Execute(w, encoded)
	}
}

func readClipboard() (string, error) {
	out, err := exec.Command("wl-paste").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func open(url string) error {
	var args []string
	args = append(args, "--new-window")
	args = append(args, url)
	return exec.Command("firefox", args...).Start()
}

func shutdown(srv *http.Server) {
	time.Sleep(3 * time.Second)
	srv.Shutdown(context.TODO())
}
