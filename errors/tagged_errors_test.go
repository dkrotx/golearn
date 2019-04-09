package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestNewErrorWFields(t *testing.T) {
	orig := errors.New("error")

	withFields := NewErrorWFields(errors.Wrap(orig, "additional message"))
	assert.Equal(t, orig, errors.Cause(withFields))
}

func TestGetFields(t *testing.T) {
	orig := errors.New("error")
	err := NewErrorWFields(
		NewErrorWFields(orig, zap.String("service", "odin")),
		zap.String("updated_by", "jeniffer"))

	for _, fld := range GetFields(err) {
		fmt.Println(fld.Key, fld.String)
	}
}