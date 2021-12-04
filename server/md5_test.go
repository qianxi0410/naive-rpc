package server_test

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMd5(t *testing.T) {
	b := md5.New().Sum([]byte("hello"))
	log.Printf("%s", b)
	b2 := md5.New().Sum([]byte("hello"))
	log.Printf("%s", b2)

	assert.Equal(t, b, b2)
}

func createRandomStr(name string) string {
	t := time.Now()
	h := md5.New()
	io.WriteString(h, name)
	io.WriteString(h, t.String())
	return fmt.Sprintf("%x", h.Sum(nil))
}

func TestRandomStr(t *testing.T) {
	log.Printf("%s --- %s", createRandomStr("qianxi"), createRandomStr("qianxi"))
}
