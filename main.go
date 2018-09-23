package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/yakshaving.art/propaganda/configuration"
	"gitlab.com/yakshaving.art/propaganda/core"
	"gitlab.com/yakshaving.art/propaganda/github"
	"gitlab.com/yakshaving.art/propaganda/gitlab"
	"gitlab.com/yakshaving.art/propaganda/metrics"
	"gitlab.com/yakshaving.art/propaganda/server"
	"gitlab.com/yakshaving.art/propaganda/slack"
	"gitlab.com/yakshaving.art/propaganda/version"

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

	go func() {
		logrus.Fatal(s.ListenAndServe(args.Address))
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2)

	for sig := range signalCh {
		switch sig {
		case syscall.SIGHUP:
			logrus.Info("Reloading the configuration")
			loadConfiguration(args)

		case syscall.SIGUSR1:
			toggleLogLevel()

		case syscall.SIGINT:
			logrus.Info("Shutting down gracefully")
			s.Shutdown()
			os.Exit(0)
		}
	}
}

func setupLogger() {
	logrus.AddHook(filename.NewHook())
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func toggleLogLevel() {
	switch logrus.GetLevel() {
	case logrus.DebugLevel:
		logrus.Infof("setting info log level")
		logrus.SetLevel(logrus.InfoLevel)
	default:
		logrus.Infof("settings debug log level")
		logrus.SetLevel(logrus.DebugLevel)
	}
}

// Args represents the commandline arguments
type Args struct {
	Address     string
	MetricsPath string

	WebhookURL  string
	MatchString string

	ConfigFile  string
	Debug       bool
	ShowVersion bool
}

func parseArgs() Args {
	var args Args

	flag.StringVar(&args.Address, "address", ":9092", "listening address")
	flag.StringVar(&args.MetricsPath, "metrics", "/metrics", "metrics path")
	flag.StringVar(&args.WebhookURL, "webhook-url", os.Getenv("SLACK_WEBHOOK_URL"), "slack webhook url")
	flag.StringVar(&args.MatchString, "match-pattern", "\\[announce\\]", "match string")
	flag.StringVar(&args.ConfigFile, "config", "propaganda.yml", "configuration file to use")
	flag.BoolVar(&args.Debug, "debug", false, "enable debug logging")
	flag.BoolVar(&args.ShowVersion, "version", false, "show version and exit")
	flag.Parse()

	if args.ShowVersion {
		logrus.Printf("Version: %s Commit: %s Date: %s", version.Version, version.Commit, version.Date)
		os.Exit(0)
	}

	if args.Debug {
		toggleLogLevel()
	}

	if args.WebhookURL == "" {
		logrus.Fatalf("no slack webhook url, define it through -webhook-url argument or SLACK_WEBHOOK_URL env var")
	}

	loadConfiguration(args)

	return args
}

func loadConfiguration(args Args) {
	content, err := ioutil.ReadFile(args.ConfigFile)
	if err != nil {
		logrus.Errorf("failed to read configuration file %s: %s", args.ConfigFile, err)
	}

	if err = configuration.Load(content); err != nil {
		logrus.Errorf("failed to load configuration file %s: %s", args.ConfigFile, err)
	}
}
