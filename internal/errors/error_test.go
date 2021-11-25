package errors_test

import (
	"fmt"
	"testing"

	"github.com/qianxi0410/naive-rpc/internal/errors"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := errors.New(114514, "error for test", errors.Business)
	assert.Equal(t, 114514, err.Code)
	assert.Equal(t, "error for test", err.Msg)
	assert.Equal(t, errors.Business, err.Type)

	err2 := errors.New(114514, "error for bro", errors.Framework)
	assert.Equal(t, 114514, err2.Code)
	assert.Equal(t, "error for bro", err2.Msg)
	assert.Equal(t, errors.Framework, err2.Type)
}

func TestErrorTypeString(t *testing.T) {
	assert.Equal(t, "framework error", fmt.Sprintf("%s", errors.Framework))
	assert.Equal(t, "bussiness error", fmt.Sprintf("%s", errors.Business))
}

func TestErrorString(t *testing.T) {
	expect := fmt.Sprintf("err code: %d, error msg: %s, error type: %s", 114514, "oops", errors.Framework)
	err := errors.New(114514, "oops", errors.Framework)
	assert.Equal(t, err.Error(), expect)
}
