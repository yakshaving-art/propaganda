package main

import (
	"flag"
	"os"

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
		slack.Announcer{
			WebhookURL: args.WebhookURL,
		},
		[]core.Parser{
			gitlab.Parser{
				MatchString: "[announce]",
			},
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

	WebhookURL string
}

func parseArgs() Args {
	var args Args

	flag.StringVar(&args.Address, "address", ":9092", "listening address")
	flag.StringVar(&args.MetricsPath, "metrics", "/metrics", "metrics path")
	flag.StringVar(&args.WebhookURL, "webhook-url", os.Getenv("SLACK_WEBHOOK_URL"), "slack webhook url")
	flag.Parse()

	if args.WebhookURL == "" {
		logrus.Fatalf("No slack webhook url, define it through -webhook-url argument or SLACK_WEBHOOK_URL env var")
	}

	return args
}
