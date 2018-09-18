package main

import (
	"fmt"
	"net/http"

	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
)

func main() {
	setupLogger()

	http.HandleFunc("/github", handleGithub)

	logrus.Fatal(http.ListenAndServe(":9999", nil))
}

func handleGithub(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logrus.Debugf("Failed to parse form on request %s: %#v", r, err)
		http.Error(w, fmt.Sprintf("bad request: %s", err), http.StatusBadRequest)
		return
	}

	payload := r.FormValue("payload")
	if payload == "" {
		logrus.Debugf("No payload in form %#v", r.Form)
		http.Error(w, "no payload in form", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	logrus.Println("Received Webhook payload: %s", payload)
}

func setupLogger() {
	logrus.AddHook(filename.NewHook())
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
}
