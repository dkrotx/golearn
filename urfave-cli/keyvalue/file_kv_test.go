package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func mkTempFile(t *testing.T) string {
	t.Helper()

	fd, err := ioutil.TempFile("", "file_kv_test_*.txt")
	require.NoError(t, err)
	return fd.Name()
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
