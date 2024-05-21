package protocol

import "github.com/cdmistman/watchman/protocol/query"

/*
["subscribe","/tmp","sub1",{"fields":["exists","name","type"]}]
{"clock":"c:1531594843:978:9:826","subscribe":"sub1","version":"4.9.0"}
{"unilateral":true,
 "subscription":"sub1",
 "root":"/tmp",
 "files":[{
  "exists", true,
  "name": "foo/main.go",
  "type": "f"
 }, {
  "exists", true,
  "name": "bar/main.go",
  "type": "f"
 }],
 "version":"4.9.0",
 "clock":"c:1531594843:978:9:826",
 "is_fresh_instance":true}
{"unilateral":true,
 "subscription":"sub1",
 "root":"/tmp",
 "files":[{
  "name": "foo/main.go",
  "type": "f"
 }],
 "version":"4.9.0",
 "since":"c:1531594843:978:9:826",
 "clock":"c:1531594843:978:9:827",
 "is_fresh_instance":false}
{"unilateral":true,"subscription":"sub1","root":"/tmp","canceled":true,"version":"4.9.0"}
*/

// A SubscribeRequest represents the Watchman subscribe command.
//
// See also: https://facebook.github.io/watchman/docs/cmd/subscribe.html
type SubscribeRequest struct {
	Root  string
	Name  string
	Query *query.Query
}

// Args returns values used to encode a request PDU.
func (req *SubscribeRequest) Args() []interface{} {
	res := []interface{}{
		"subscribe",
		req.Root,
		req.Name,
	}

	if req.Query != nil {
		res = append(res, req.Query)
	}

	return res
}

// A SubscribeResponse represents a response to the Watchman subscribe command.
type SubscribeResponse struct {
	response
	clock        string
	subscription string
}

// NewSubscribeResponse converts a ResponsePDU to SubscribeResponse
func NewSubscribeResponse(pdu ResponsePDU) (res *SubscribeResponse) {
	res = &SubscribeResponse{}
	res.response.init(pdu)

	if x, ok := pdu["clock"]; ok {
		if clock, ok := x.(string); ok {
			res.clock = clock
		}
	}
	if x, ok := pdu["subscribe"]; ok {
		if subscription, ok := x.(string); ok {
			res.subscription = subscription
		}
	}
	return
}

// Clock returns a value representing when the subscription started.
func (res *SubscribeResponse) Clock() string {
	return res.clock
}

// Subscription returns the name registered to the subscription.
func (res *SubscribeResponse) Subscription() string {
	return res.subscription
}

// A Subscription represents notifications generated as a result
// of the Watchman subscribe command.
type Subscription struct {
	response
	clock           string
	root            string
	subscription    string
	files           []interface{}
	isFreshInstance bool
}

// NewSubscription converts a ResponsePDU to Subscription
func NewSubscription(pdu ResponsePDU) (s *Subscription) {
	s = &Subscription{}
	s.response.init(pdu)

	if x, ok := pdu["clock"]; ok {
		if clock, ok := x.(string); ok {
			s.clock = clock
		}
	}
	if x, ok := pdu["files"]; ok {
		if files, ok := x.([]interface{}); ok {
			s.files = files
		}
	}
	if x, ok := pdu["is_fresh_instance"]; ok {
		if isFreshInstance, ok := x.(bool); ok {
			s.isFreshInstance = isFreshInstance
		}
	}
	if x, ok := pdu["root"]; ok {
		if root, ok := x.(string); ok {
			s.root = root
		}
	}
	if x, ok := pdu["subscription"]; ok {
		if subscription, ok := x.(string); ok {
			s.subscription = subscription
		}
	}
	return
}

// Clock returns a value representing when the notification was
// generated.
func (s *Subscription) Clock() string {
	return s.clock
}

// Files returns a list of of relative filepaths.
func (s *Subscription) Files() []interface{} {
	return s.files
}

// IsFreshInstance indicates if the notification was sent because
// of a newly established subscription, or observed changes.
func (s *Subscription) IsFreshInstance() bool {
	return s.isFreshInstance
}

// Root returns the directory registered to the subscription.
func (s *Subscription) Root() string {
	return s.root
}

// Subscription returns the name registered to the subscription.
func (s *Subscription) Subscription() string {
	return s.subscription
}
