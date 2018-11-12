package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func mkTempFile(t *testing.T) string {
	t.Helper()

	fd, err := ioutil.TempFile("", "file_kv_test_*.txt")
	require.NoError(t, err)

	path := fd.Name()

	require.NoError(t, fd.Close())
	require.NoError(t, os.Remove(path)) // we should start w/o file
	return path
}

func TestSetAndGetWorks(t *testing.T) {
	tempFile := mkTempFile(t)
	defer os.Remove(tempFile)

	kv, err := NewFileKeyValue(tempFile)
	require.NoError(t, err)

	require.NoError(t, kv.Set("name", "Jan"))

	value, err := kv.Get("name")
	assert.NoError(t, err)
	assert.Equal(t, "Jan", value)
}

func TestTTLWorks(t *testing.T) {
	tempFile := mkTempFile(t)
	defer os.Remove(tempFile)

	kv, err := NewFileKeyValue(tempFile)
	require.NoError(t, err)

	require.NoError(t, kv.Set("name", "Jan", TTL(time.Second / 2)))

	value, err := kv.Get("name")
	require.NoError(t, err, "key should not expire yet")
	assert.Equal(t, "Jan", value)

	time.Sleep(600 * time.Millisecond)
	_, err = kv.Get("name")
	require.Error(t, err, "key must be expired!")
	assert.IsType(t, KeyNotFoundErr, err)
}

func TestGetReturnsErrorWhenKeyNotFound(t *testing.T) {
	tempFile := mkTempFile(t)
	defer os.Remove(tempFile)

	kv, err := NewFileKeyValue(tempFile)
	require.NoError(t, err)

	require.NoError(t, kv.Set("name", "Jan"))

	value, err := kv.Get("family")
	require.Error(t, err)
	assert.IsType(t, KeyNotFoundErr, err)
	assert.Empty(t, value, "return value in case of error should be empty")
}

func TestDelete(t *testing.T) {
	tempFile := mkTempFile(t)
	defer os.Remove(tempFile)

	kv, err := NewFileKeyValue(tempFile)
	require.NoError(t, err)

	require.NoError(t, kv.Set("name", "Jan"))
	require.NoError(t, kv.Set("family", "Kisel"))

	require.NoError(t, kv.Delete("name"))

	_, err = kv.Get("name")
	require.Error(t, err)
	assert.IsType(t, KeyNotFoundErr, err)

	value, err := kv.Get("family")
	require.NoError(t, err)
	assert.Equal(t, "Kisel", value)

	// remove last key-value
	require.NoError(t, kv.Delete("family"))
	_, err = kv.Get("family")
	require.Error(t, err)
	assert.IsType(t, KeyNotFoundErr, err)
}

