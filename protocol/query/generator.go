package query

import (
	"strconv"
)

// See https://facebook.github.io/watchman/docs/file-query#generators
type Generator string

const (
	// See https://facebook.github.io/watchman/docs/file-query#since-generator
	GSince Generator = "string"
	// See https://facebook.github.io/watchman/docs/file-query#suffix-generator
	GSuffix Generator = "suffix"
	// See https://facebook.github.io/watchman/docs/file-query#glob-generator
	GGlob Generator = "glob"
	// See https://facebook.github.io/watchman/docs/file-query#path-generator
	GPath Generator = "path"
)

type Generators map[Generator]any

// See https://facebook.github.io/watchman/docs/file-query#path-generator
type GPathPath struct {
	Path  string
	Depth int
}

func (p GPathPath) MarshalJSON() ([]byte, error) {
	if p.Depth == -1 {
		return []byte(`"` + p.Path + `"`), nil
	}

	return []byte(`["` + p.Path + `", ` + strconv.Itoa(p.Depth) + `]`), nil
}
