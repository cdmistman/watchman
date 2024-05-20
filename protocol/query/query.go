package query

import "encoding/json"

type Query struct {
	Generators   Generators
	Expression   Term
	Fields       Fields
	DedupResults bool
	RelativeRoot string
	SyncTimeout  int
	LockTimeout  int
	Case         Case
}

type Case int

const (
	CaseSensitive Case = iota
	CaseInsensitive
)

func (q Query) MarshalJSON() ([]byte, error) {
	res := map[string]any{}

	for generator, arg := range q.Generators {
		res[string(generator)] = arg
	}

	if q.Expression != nil {
		res["expression"] = q.Expression
	}

	if q.Fields != nil {
		res["fields"] = q.Fields
	}

	if q.DedupResults {
		res["dedup_results"] = true
	}

	if q.RelativeRoot != "" {
		res["relative_root"] = q.RelativeRoot
	}

	if q.SyncTimeout != 0 {
		res["sync_timeout"] = q.SyncTimeout
	}

	if q.LockTimeout != 0 {
		res["lock_timeout"] = q.LockTimeout
	}

	if q.Case != CaseSensitive {
		res["case_sensitive"] = false
	}

	return json.Marshal(res)
}
