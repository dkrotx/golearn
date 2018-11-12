package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	tm := Timestamp{Time: time.Unix(1542017388, 0).UTC()}

	bytes, err := tm.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `"2018-11-12T10:09:48"`, string(bytes))

	tmRestored := Timestamp{}

	err = tmRestored.UnmarshalJSON(bytes)
	require.NoError(t, err)
	assert.Equal(t, tm, tmRestored, "decoded != original?!")
}

func TestTimestampWithinJSON(t *testing.T) {
	type record struct {
		Stamp Timestamp `json:"stamp"`
	}

	rec := record {
		Stamp: Timestamp{Time: time.Unix(1542017388, 0).UTC()},
	}

	bytes, err := json.Marshal(&rec)
	require.NoError(t, err)
	assert.Equal(t, `{"stamp":"2018-11-12T10:09:48"}`, string(bytes))

	var restored record

	err = json.Unmarshal(bytes, &restored)
	require.NoError(t, err)
	assert.Equal(t, rec, restored, "decoded != original?!")
}
