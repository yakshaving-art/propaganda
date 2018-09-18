package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
	// This requires registering the webhooks using json format and only receive
	// pull request events

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("failed to read body: %s", err)
		http.Error(w, fmt.Sprintf("bad request: %s", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusAccepted)

	logrus.Infof("received Webhook payload: %s", string(body))
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
