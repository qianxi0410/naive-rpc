package errors_test

import (
	"fmt"
	"testing"

	"github.com/qianxi0410/naive-rpc/errors"
	"github.com/stretchr/testify/assert"
)

func TestIsFrameworkError(t *testing.T) {
	fe := errors.CodecDecodeErr
	assert.Equal(t, true, errors.IsFrameworkError(fe))
	assert.Equal(t, false, errors.IsFrameworkError(nil))
	assert.Equal(t, false, errors.IsFrameworkError(fmt.Errorf("")))
	assert.Equal(t, false, errors.IsFrameworkError(errors.New(1, "")))
}

func TestIsBussinessError(t *testing.T) {
	be := errors.New(114514, "bussiness error")
	assert.Equal(t, true, errors.IsBusinessError(be))
	assert.Equal(t, false, errors.IsBusinessError(nil))
	assert.Equal(t, false, errors.IsBusinessError(fmt.Errorf("")))
	assert.Equal(t, false, errors.IsBusinessError(errors.RouteNotFoundErr))
}

func TestErrorCode(t *testing.T) {
	be := errors.New(666, "")
	assert.Equal(t, 666, errors.ErrorCode(be))
	assert.Equal(t, 0, errors.ErrorCode(nil))
	assert.Equal(t, 0, errors.ErrorCode(fmt.Errorf("")))
	assert.Equal(t, 2000, errors.ErrorCode(errors.CodecDecodeErr))
}

func TestErrorMsg(t *testing.T) {
	assert.Equal(t, "123", errors.ErrorMsg(fmt.Errorf("123")))
	assert.Equal(t, "", errors.ErrorMsg(nil))
	assert.Equal(t, fmt.Sprintf("err code: %d, error msg: %s, error type: %s", 666, "oops", "bussiness error"), errors.ErrorMsg(errors.New(666, "oops")))
}
