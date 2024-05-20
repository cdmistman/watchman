package watchman

import (
	"github.com/cdmistman/watchman/protocol"
)

// A ChangeNotification represents changes two one or more filesystem entries.
type ChangeNotification struct {
	IsFreshInstance bool
	Clock           string
	Subscription    string
	Files           []interface{}
}

func newChangeNotification(sub *protocol.Subscription) *ChangeNotification {
	clock := sub.Clock()
	files := sub.Files()
	cn := &ChangeNotification{
		IsFreshInstance: sub.IsFreshInstance(),
		Clock:           clock,
		Subscription:    sub.Subscription(),
		Files:           files,
	}
	return cn
}
