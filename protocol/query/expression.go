package query

import "encoding/json"

type Term json.Marshaler

// see https://facebook.github.io/watchman/docs/expr/allof
type TAllof []Term

func (t TAllof) MarshalJSON() ([]byte, error) {
	res := []any{"allof"}
	for _, term := range t {
		res = append(res, term)
	}

	return json.Marshal(res)
}

// see https://facebook.github.io/watchman/docs/expr/anyof
type TAnyof []Term

func (t TAnyof) MarshalJSON() ([]byte, error) {
	res := []any{"anyof"}
	for _, term := range t {
		res = append(res, term)
	}
	return json.Marshal(res)
}

// see https://facebook.github.io/watchman/docs/expr/size
type RelationalOp int

const (
	RelLt RelationalOp = iota
	RelLe
	RelEq
	RelNe
	RelGt
	RelGe
)

var relationalOpMap = map[RelationalOp]string{
	RelLt: "lt",
	RelLe: "le",
	RelEq: "eq",
	RelNe: "ne",
	RelGt: "gt",
	RelGe: "ge",
}

func (op RelationalOp) MarshalJSON() ([]byte, error) {
	return json.Marshal(relationalOpMap[op])
}

// see https://facebook.github.io/watchman/docs/expr/dirname
type TDirname struct {
	Name  string
	Op    RelationalOp
	Depth int
}

func (t TDirname) MarshalJSON() ([]byte, error) {
	res := []any{"dirname", t.Name}
	if t.Op != RelLe && t.Depth != 0 {
		res = append(res, t.Depth)
	}
	return json.Marshal(res)
}

// see https://facebook.github.io/watchman/docs/expr/dirname
type TIDirname struct {
	Name  string
	Op    RelationalOp
	Depth int
}

func (t TIDirname) MarshalJSON() ([]byte, error) {
	res := []any{"idirname", t.Name}
	if t.Op != RelLe && t.Depth != 0 {
		res = append(res, t.Depth)
	}
	return json.Marshal(res)
}

// see https://facebook.github.io/watchman/docs/expr/empty
type TEmpty struct{}

var EmptyT = TEmpty{}

func (t TEmpty) MarshalJSON() ([]byte, error) {
	return json.Marshal("empty")
}

// see https://facebook.github.io/watchman/docs/expr/exists
type TExists struct{}

var ExistsT = TExists{}

func (t TExists) MarshalJSON() ([]byte, error) {
	return json.Marshal("exists")
}

// see https://facebook.github.io/watchman/docs/expr/false
type TFalse struct{}

var FalseT = TFalse{}

func (t TFalse) MarshalJSON() ([]byte, error) {
	return json.Marshal("false")
}

// see https://facebook.github.io/watchman/docs/expr/match
type TMatchType int

const (
	MatchBaseName TMatchType = iota
	MatchWholeName
)

func (t TMatchType) MarshalJSON() ([]byte, error) {
	var str string
	switch t {
	case MatchBaseName:
		str = "basename"
	case MatchWholeName:
		str = "wholename"
	}
	return json.Marshal(str)
}

type TMatchFlags int

const (
	MatchIncludeDotFiles TMatchFlags = 1 << iota
	MatchNoEscape
)

// see https://facebook.github.io/watchman/docs/expr/match
type TMatch struct {
	Glob      string
	MatchType TMatchType
	Flags     TMatchFlags
}

func (t TMatch) MarshalJSON() ([]byte, error) {
	res := []any{"match", t.Glob}
	if t.MatchType != MatchBaseName || t.Flags != 0 {
		res = append(res, t.MatchType)
	}

	if t.Flags != 0 {
		flags := map[string]bool{}
		if t.Flags&MatchIncludeDotFiles != 0 {
			flags["includedotfiles"] = true
		}

		if t.Flags&MatchNoEscape != 0 {
			flags["noescape"] = true
		}

		res = append(res, flags)
	}

	return json.Marshal(res)
}

// see https://facebook.github.io/watchman/docs/expr/match
type TIMatch struct {
	Glob      string
	MatchType TMatchType
	Flags     TMatchFlags
}

func (t TIMatch) MarshalJSON() ([]byte, error) {
	res := []any{"imatch", t.Glob}
	if t.MatchType != MatchBaseName || t.Flags != 0 {
		res = append(res, t.MatchType)
	}

	if t.Flags != 0 {
		flags := map[string]bool{}
		if t.Flags&MatchIncludeDotFiles != 0 {
			flags["includedotfiles"] = true
		}

		if t.Flags&MatchNoEscape != 0 {
			flags["noescape"] = true
		}

		res = append(res, flags)
	}

	return json.Marshal(res)
}

// see https://facebook.github.io/watchman/docs/expr/name
type TName struct {
	Names     []string
	MatchType TMatchType
}

func (t TName) MarshalJSON() ([]byte, error) {
	res := []any{"name"}
	if len(t.Names) == 1 {
		res = append(res, t.Names[0])
	} else {
		res = append(res, t.Names)
	}

	if t.MatchType != MatchBaseName {
		res = append(res, t.MatchType)
	}

	return json.Marshal(res)
}

// see https://facebook.github.io/watchman/docs/expr/name
type TIName struct {
	Names     []string
	MatchType TMatchType
}

func (t TIName) MarshalJSON() ([]byte, error) {
	res := []any{"iname"}
	if len(t.Names) == 1 {
		res = append(res, t.Names[0])
	} else {
		res = append(res, t.Names)
	}

	if t.MatchType != MatchBaseName {
		res = append(res, t.MatchType)
	}

	return json.Marshal(res)
}

type TNot struct {
	Not Term
}

func (t TNot) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{"not", t.Not})
}

// see https://facebook.github.io/watchman/docs/expr/pcre
type TPCRE struct {
	Regexp    string
	MatchType TMatchType
}

func (t TPCRE) MarshalJSON() ([]byte, error) {
	res := []any{"pcre", t.Regexp}
	if t.MatchType != MatchBaseName {
		res = append(res, t.MatchType)
	}
	return json.Marshal(res)
}

// see https://facebook.github.io/watchman/docs/expr/pcre
type TIPCRE struct {
	Regexp    string
	MatchType TMatchType
}

func (t TIPCRE) MarshalJSON() ([]byte, error) {
	res := []any{"ipcre", t.Regexp}
	if t.MatchType != MatchBaseName {
		res = append(res, t.MatchType)
	}
	return json.Marshal(res)
}

// see https://facebook.github.io/watchman/docs/expr/since
type TClockSource int

const (
	OClock TClockSource = iota
	CClock
	CTime
	MTime
)

// see https://facebook.github.io/watchman/docs/expr/since
type TSince struct {
	Timestamp any
	Source    TClockSource
}

func (t TSince) MarshalJSON() ([]byte, error) {
	res := []any{"since", t.Timestamp}
	if t.Source != OClock {
		res = append(res, t.Source)
	}

	return json.Marshal(res)
}

// see https://facebook.github.io/watchman/docs/expr/size
type TSize struct {
	Op   RelationalOp
	Size int
}

func (t TSize) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{"size", t.Op, t.Size})
}

// see https://facebook.github.io/watchman/docs/expr/suffix
type TSuffix []string

func (t TSuffix) MarshalJSON() ([]byte, error) {
	res := []any{"suffix"}
	if len(t) == 1 {
		res = append(res, t[0])
	} else {
		res = append(res, t)
	}
	return json.Marshal(res)
}

// see https://facebook.github.io/watchman/docs/expr/true
type TTrue struct{}

var TrueT = TTrue{}

func (t TTrue) MarshalJSON() ([]byte, error) {
	return json.Marshal("true")
}

// see https://facebook.github.io/watchman/docs/expr/type
type TFileType string

const (
	TFileBlock       TFileType = "b"
	TFileChar        TFileType = "c"
	TFileDir         TFileType = "d"
	TFileRegular     TFileType = "f"
	TFileFIFO        TFileType = "p"
	TFileLink        TFileType = "l"
	TFileSocket      TFileType = "s"
	TFileSolarisDoor TFileType = "D"
	TFileUnknown     TFileType = "?"
)

func (t TFileType) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{"type", string(t)})
}
