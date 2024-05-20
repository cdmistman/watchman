package watchman_test

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/stretchr/testify/require"

	"github.com/cdmistman/watchman"
	"github.com/cdmistman/watchman/protocol/query"
)

const pause = 250 * time.Millisecond

func collect(updates <-chan interface{}) []interface{} {
	messages := make([]interface{}, 0, 3)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		timeout := time.After(pause)
		for {
			select {
			case msg := <-updates:
				messages = append(messages, msg)
			case <-timeout:
				return
			}
		}
	}()

	wg.Wait()
	return messages
}

func tmpdir(t *testing.T) (string, error) {
	dir := t.TempDir()

	dir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, ".watchmanconfig")
	err = os.WriteFile(path, []byte(`{"idle_reap_age_seconds": 300}`+"\n"), os.ModePerm)
	return dir, err
}

func mkdir(dir string, names ...string) error {
	for _, name := range names {
		name = filepath.Join(dir, name)
		err := os.Mkdir(name, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func remove(dir string, names ...string) error {
	for _, name := range names {
		path := filepath.Join(dir, name)
		err := os.Remove(path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func symlink(dir, name, target string) error {
	name = filepath.Join(dir, name)
	return os.Symlink(target, name)
}

func touch(dir string, names ...string) error {
	data := []byte("Kilroy was here.")
	for _, name := range names {
		path := filepath.Join(dir, name)

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}

		n, err := f.Write(data)
		if err == nil && n < len(data) {
			return io.ErrShortWrite
		}

		if err = f.Sync(); err != nil {
			return err
		}

		if err = f.Close(); err != nil {
			return err
		}
	}
	return nil
}

func TestClient(t *testing.T) {
	require := require.New(t)
	defer leaktest.Check(t)()

	dir, err := tmpdir(t)
	require.NoError(err)
	defer os.RemoveAll(dir)

	// connect
	c, err := watchman.Connect()
	require.NoError(err)

	// connection metadata
	require.NotEmpty(c.SockName())
	require.NotEmpty(c.Version())

	// capabilities
	require.Equal(true, c.HasCapability("cmd-subscribe"))
	require.Equal(false, c.HasCapability("grant-three-wishes"))

	// watch-project
	watch, err := c.AddWatch(dir)
	require.NoError(err)

	updates := c.Notifications()
	n := len(collect(updates))
	require.Equal(0, n)

	// watch-list
	roots, err := c.ListWatches()
	require.NoError(err)
	require.NotEmpty(roots)

	// subscribe
	s, err := watch.Subscribe("Spoon!", &query.Query{
		Fields: query.Fields{query.FName, query.FType, query.FNew},
	})

	require.NoError(err)

	n = len(collect(updates))
	require.NotEqual(0, n)

	// clock
	clock1, err := watch.Clock(0)
	require.NoError(err)
	require.NotEmpty(clock1)

	err = remove(dir, "qux", "quux")
	require.NoError(err)

	err = touch(dir, "foo", "bar", "baz")
	require.NoError(err)

	n = len(collect(updates))
	require.NotEqual(0, n)

	clock2, err := watch.Clock(pause)
	require.NoError(err)
	require.NotEmpty(clock2)
	require.NotEqual(clock1, clock2)

	// state changes
	err = touch(dir, "baz", "qux", "quux")
	require.NoError(err)

	err = remove(dir, "foo", "bar", "quux")
	require.NoError(err)

	err = mkdir(dir, "corge", "grault")
	require.NoError(err)

	switch runtime.GOOS {
	case "darwin", "freebsd", "linux":
		err = symlink(dir, "garply", "grault")
		require.NoError(err)
	}

	messages := collect(updates)
	for _, msg := range messages {
		cn, ok := msg.(*watchman.ChangeNotification)
		if !ok || cn.IsFreshInstance {
			continue
		}

		for _, f := range cn.Files {
			file, ok := f.(map[string]interface{})
			require.True(ok)
			n, ok := file["name"]
			require.True(ok)
			name, ok := n.(string)
			require.True(ok)
			switch name {
			case "foo", "bar":
				tyVal, ok := file["type"]
				require.True(ok)
				ty, ok := tyVal.(string)
				require.True(ok)
				require.Equal("f", ty)

				// changeVal, ok := file["change"]
				// require.True(ok)
				// change, ok := changeVal.(watchman.StateChange)
				// require.True(ok)
				// require.Equal(watchman.Removed, change)

			case "baz":
				tyVal, ok := file["type"]
				require.True(ok)
				ty, ok := tyVal.(string)
				require.True(ok)
				require.Equal("f", ty)

				// changeVal, ok := file["change"]
				// require.True(ok)
				// change, ok := changeVal.(watchman.StateChange)
				// require.True(ok)
				// require.Equal(watchman.Updated, change)

			case "qux":
				tyVal, ok := file["type"]
				require.True(ok)
				ty, ok := tyVal.(string)
				require.True(ok)
				require.Equal("f", ty)

				// changeVal, ok := file["change"]
				// require.True(ok)
				// change, ok := changeVal.([]watchman.StateChange)
				// require.True(ok)
				// // require.Equal(watchman.Created, file.Change)
				// require.Contains(
				// 	[]watchman.StateChange{watchman.Created, watchman.Updated},
				// 	change,
				// )

			case "quux":
				// require.Equal("?", file.Type)
				// require.Equal(watchman.Ephemeral, file.Change)
				tyVal, ok := file["type"]
				require.True(ok)
				ty, ok := tyVal.(string)
				require.True(ok)
				require.Contains("f?", ty)

				// changeVal, ok := file["change"]
				// require.True(ok)
				// change, ok := changeVal.([]watchman.StateChange)
				// require.True(ok)
				// require.Contains(
				// 	[]watchman.StateChange{watchman.Ephemeral, watchman.Removed},
				// 	change,
				// )

			case "corge", "grault":
				tyVal, ok := file["type"]
				require.True(ok)
				ty, ok := tyVal.(string)
				require.True(ok)
				require.Equal("d", ty)

				// changeVal, ok := file["change"]
				// require.True(ok)
				// change, ok := changeVal.([]watchman.StateChange)
				// require.True(ok)
				// // require.Equal(watchman.Created, file.Change)
				// require.Contains(
				// 	[]watchman.StateChange{watchman.Created, watchman.Updated},
				// 	change,
				// )

			case "garply":
				tyVal, ok := file["type"]
				require.True(ok)
				ty, ok := tyVal.(string)
				require.True(ok)
				require.Equal("l", ty)

				// changeVal, ok := file["change"]
				// require.True(ok)
				// change, ok := changeVal.([]watchman.StateChange)
				// require.True(ok)
				// require.Equal(watchman.Created, change)
			}
		}
	}

	// unsubscribe
	err = s.Unsubscribe()
	require.NoError(err)

	// close
	err = c.Close()
	require.NoError(err)
}
