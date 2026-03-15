package logstore_test

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lucasmaehn/journl/logstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sut *logstore.SQLiteLogstore

func TestMain(m *testing.M) {
	var err error
	sut, err = logstore.NewSQLite("./test-sqlite.db")
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func Test_CommitAndList(t *testing.T) {
	entry := "test-commit-" + uuid.NewString()

	err := sut.Commit(entry)
	require.NoError(t, err, "failed to commit a journal entry")

	entries, err := sut.List()
	require.NoError(t, err, "failed to lsit journal entries")
	require.NotEmpty(t, entries, "journal should not be empty after adding one entry")

	lastEntry := entries[len(entries)-1]
	assert.Equal(t, entry, lastEntry.Text)
	assert.WithinDuration(t, time.Now(), lastEntry.Timestamp, time.Second)
}
