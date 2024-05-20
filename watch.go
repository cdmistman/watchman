package watchman

import (
	"path"
	"time"

	"github.com/cdmistman/watchman/protocol"
	"github.com/cdmistman/watchman/protocol/query"
)

// A Watch represents a directory, or watched root, that Watchman is watching for changes.
type Watch struct {
	client *Client
	root   string
	rel    string
}

// Clock returns the current clock value for a watched root.
//
// For details, see: https://facebook.github.io/watchman/docs/cmd/clock.html
func (w *Watch) Clock(syncTimeout time.Duration) (string, error) {
	timeout := syncTimeout.Nanoseconds() / int64(time.Millisecond)
	req := &protocol.ClockRequest{
		Path:        w.root,
		SyncTimeout: int(timeout),
	}
	pdu, err := w.client.send(req)
	var clock string
	if err == nil {
		res := protocol.NewClockResponse(pdu)
		clock = res.Clock()
	}
	return clock, err
}

// Subscribe requests notification when changes occur under a watched root.
func (w *Watch) Subscribe(name string, query *query.Query) (s *Subscription, err error) {
	req := &protocol.SubscribeRequest{
		Name:  name,
		Root:  w.root,
		Query: query,
	}

	if w.rel != "" {
		req.Query.RelativeRoot = w.rel
	}

	_, err = w.client.send(req)
	if err == nil {
		s = &Subscription{
			client: w.client,
			name:   name,
			root:   path.Join(w.root, w.rel),
		}
	}
	return
}

func (w *Watch) Root() string {
	return w.root
}

func (w *Watch) RelativePath() string {
	return w.rel
}
