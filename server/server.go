package server

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.com/yakshaving.art/propaganda/core"
	"gitlab.com/yakshaving.art/propaganda/metrics"

	"github.com/sirupsen/logrus"
)

// New returns a new Server with the provided parsers
func New(announcer core.Announcer, parsers []core.Parser) Server {
	return Server{
		parsers:   parsers,
		announcer: announcer,
	}
}

// Server is the http serving object
type Server struct {
	parsers   []core.Parser
	announcer core.Announcer
}

// ListenAndServe starts listening and serving traffic
func (s Server) ListenAndServe(addr string) error {
	http.HandleFunc("/", s.handle)

	metrics.Up.Set(1)
	logrus.Infof("listening on %s", addr)
	return http.ListenAndServe(addr, nil)
}

func (s Server) handle(w http.ResponseWriter, r *http.Request) {
	// This requires registering the webhooks using json format and only receive
	// pull request events

	metrics.WebhooksReceived.Inc()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		metrics.WebhooksErrors.Inc()

		logrus.Errorf("failed to read body: %s", err)
		http.Error(w, fmt.Sprintf("bad request: %s", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	metrics.WebhooksBytesRead.Add(float64(len(body)))

	for _, p := range s.parsers {
		if p.Match(r.Header) {
			a, err := p.Parse(body)
			if err != nil {
				http.Error(w, fmt.Sprintf("Parser failed to parse: %s", err), http.StatusUnprocessableEntity)
				metrics.WebhooksInvalid.WithLabelValues("failed_parsing").Inc()
				return
			}

			if !a.ShouldAnnounce() {
				w.WriteHeader(http.StatusOK)
				metrics.WebhooksInvalid.WithLabelValues("non_announceable").Inc()
				return
			}

			metrics.WebhooksValid.WithLabelValues(a.ProjectName()).Inc()

			s.announcer.Announce(a)

			w.WriteHeader(http.StatusAccepted)
			logrus.Debugf("received Webhook\nHeaders: %#v\nPayload: %s", r.Header, string(body))
			return
		}
	}

	metrics.WebhooksInvalid.WithLabelValues("no_parser").Inc()
	http.Error(w, fmt.Sprintf("No parser defined for this hook"), http.StatusUnprocessableEntity)
}

func (s Server) process(a core.Announcement) {

}
