package main

import (
	"flag"

	"gitlab.com/yakshaving.art/propaganda/core"
	"gitlab.com/yakshaving.art/propaganda/gitlab"
	"gitlab.com/yakshaving.art/propaganda/metrics"
	"gitlab.com/yakshaving.art/propaganda/server"
	"gitlab.com/yakshaving.art/propaganda/slack"

	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
)

func main() {
	setupLogger()

	args := parseArgs()

	metrics.Register(args.MetricsPath)

	s := server.New(
		slack.Announcer{},
		[]core.Parser{
			gitlab.Parser{},
		})

	logrus.Fatal(s.ListenAndServe(args.Address))
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
	Address     string
	MetricsPath string
}

func parseArgs() Args {
	var args Args

	flag.StringVar(&args.Address, "address", ":9092", "listening address")
	flag.StringVar(&args.MetricsPath, "metrics", "/metrics", "metrics path")
	flag.Parse()

	return args
}
