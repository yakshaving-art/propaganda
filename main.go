package main

import (
	"flag"
	"os"

	"gitlab.com/yakshaving.art/propaganda/core"
	"gitlab.com/yakshaving.art/propaganda/github"
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
		slack.Announcer{
			WebhookURL: args.WebhookURL,
		},
		[]core.Parser{
			github.NewParser(args.MatchString),
			gitlab.NewParser(args.MatchString),
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

	WebhookURL  string
	MatchString string

	ConfigFile string
}

func parseArgs() Args {
	var args Args

	flag.StringVar(&args.Address, "address", ":9092", "listening address")
	flag.StringVar(&args.MetricsPath, "metrics", "/metrics", "metrics path")
	flag.StringVar(&args.WebhookURL, "webhook-url", os.Getenv("SLACK_WEBHOOK_URL"), "slack webhook url")
	flag.StringVar(&args.MatchString, "match-pattern", "\\[announce\\]", "match string")
	flag.StringVar(&args.ConfigFile, "config", "propaganda.yml", "configuration file to use")
	flag.Parse()

	if args.WebhookURL == "" {
		logrus.Fatalf("no slack webhook url, define it through -webhook-url argument or SLACK_WEBHOOK_URL env var")
	}

	if _, err := os.Stat(args.ConfigFile); err != nil {
		logrus.Fatalf("failed to stat configuration file %s: %s", args.ConfigFile, err)
	}

	return args
}
