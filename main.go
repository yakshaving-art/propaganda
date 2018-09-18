package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
)

func main() {
	setupLogger()

	args := parseArgs()

	http.HandleFunc("/github", handleGithub)

	logrus.Infof("listening on %s", args.Address)
	logrus.Fatal(http.ListenAndServe(args.Address, nil))
}

func handleGithub(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logrus.Debugf("failed to parse form on request %s: %#v", r, err)
		http.Error(w, fmt.Sprintf("bad request: %s", err), http.StatusBadRequest)
		return
	}

	payload := r.FormValue("payload")
	if payload == "" {
		logrus.Debugf("no payload in form %#v", r.Form)
		http.Error(w, "no payload in form", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	logrus.Println("received Webhook payload: %s", payload)
}

func setupLogger() {
	logrus.AddHook(filename.NewHook())
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
}

// Args represents the commandline arguments
type Args struct {
	Address string
}

func parseArgs() Args {
	var args Args

	flag.StringVar(&args.Address, "address", ":9092", "listening address")
	flag.Parse()

	return args
}
