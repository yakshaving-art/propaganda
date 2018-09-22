package core

// Announcement represents a message to be shouted out
type Announcement interface {
	Title() string
	Text() string
	URL() string
	ShouldAnnounce() bool
	ProjectName() string
}

// Parser provides an interface that allows to identify if a request can be
// parsed, and then will extract the announcement if there is any.
type Parser interface {
	Match(map[string][]string) bool
	Parse([]byte) (Announcement, error)
}

// Announcer provides a simple interface to announce things
type Announcer interface {
	Announce(Announcement)
}
