package slack

import (
	"gitlab.com/yakshaving.art/propaganda/core"
)

// Announcer announces things to slack
type Announcer struct {
}

// Announce implements core.Announcer interface
func (Announcer) Announce(a core.Announcement) {

}
