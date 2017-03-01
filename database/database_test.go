package database

import (
	"testing"

	validator "gopkg.in/go-playground/validator.v9"

	"github.com/stretchr/testify/assert"
)

func TestValid(t *testing.T) {
	a := NewAccount("Alex", "email@google.com", map[string]string{})
	err := a.Valid()
	assert.Nil(t, err)
}

func TestNotValid(t *testing.T) {
	a := NewAccount("", "", map[string]string{})
	err := a.Valid()
	assert.NotNil(t, err)
	assert.Equal(t, len(err.(validator.ValidationErrors)), 2)
}
