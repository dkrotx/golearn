package main

import (
	"github.com/golang/mock/gomock"
	"golearn/mocking/mocks"
	"testing"
)

func TestWriteAndGetWithMock(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// get mocked object
	mockKV := mock_keyvalue.NewMockKeyValue(mockCtrl)

	// configure expectations and even re-define method
	mockKV.EXPECT().Set("test-key", gomock.Any()).AnyTimes()
	// mockKV.EXPECT().Get("test-key").Return("", true).Times(1)
	mockKV.EXPECT().Get(gomock.Any()).DoAndReturn(func(key string) (string, bool) {
		return "from-test", true
	})

	// provide mocked object by interface
	WriteAndGet(mockKV)
}
