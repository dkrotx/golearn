package main

import (
	"code.uber.internal/compute/control-plane.git/query-service-v2/libs/identifiers"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

type instanceInfo struct {
	ID           string               `json:"__id"`
	AuroraStatus identifiers.Status   `json:"status"`
	Revision     identifiers.Revision `json:"rev"`
	Host         identifiers.Host     `json:"host"`
	LatestChange time.Time            `json:"ts"`
}

func BenchmarkJsonUnmarshalling(t *testing.B) {
	data, err := ioutil.ReadFile("testdata/large.json")
	require.NoError(t, err)

	t.Run("json.Unmarshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var res []instanceInfo
			err = json.Unmarshal(data, &res)
			require.NoError(t, err)
		}
	})
}