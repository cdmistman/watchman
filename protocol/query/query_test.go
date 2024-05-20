package query

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type obj map[string]any

var queryTests = []struct {
	expect obj
	query  Query
}{
	{
		expect: obj{"sync_timeout": float64(1000)},
		query:  Query{SyncTimeout: 1000},
	},

	{
		expect: obj{"expression": []interface{}{"name", "foo"}},
		query:  Query{Expression: TName{Names: []string{"foo"}}},
	},

	{
		expect: obj{
			"suffix": "php",
			"expression": []any{
				"allof",
				[]any{"type", "f"},
				[]any{"not", "empty"},
				[]any{"ipcre", "test"},
			},
			"fields": []any{"name"},
		},
		query: Query{
			Generators: Generators{
				GSuffix: "php",
			},
			Expression: TAllof{
				TFileType("f"),
				TNot{EmptyT},
				TIPCRE{"test", MatchBaseName},
			},
			Fields: Fields{FName},
		},
	},

	{
		expect: obj{"fields": []any{"name", "exists", "new", "size", "mode"}},
		query:  Query{Fields: Fields{FName, FExists, FNew, FSize, FMode}},
	},

	{
		expect: obj{
			"expression":   "exists",
			"fields":       []any{"name"},
			"sync_timeout": float64(60000),
		},
		query: Query{
			Expression:  ExistsT,
			Fields:      Fields{FName},
			SyncTimeout: 60000,
		},
	},
}

func TestQueries(t *testing.T) {
	t.Parallel()

	for _, test := range queryTests {
		marshalled, err := json.Marshal(test.query)
		require.NoError(t, err)
		var actual obj
		err = json.Unmarshal(marshalled, &actual)
		require.NoError(t, err)
		require.Equal(t, test.expect, actual)
	}
}
