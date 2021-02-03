package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckHost(t *testing.T) {
	isTest = true
	ok, err := checkHost("https://www.pixiv.net/")
	if err != nil {
		t.Error(err)
	}
	assert.True(t, ok, "pixiv.net 应该有效")
	ok, err = checkHost("https://www.baidu.com/")
	if err != nil {
		t.Error(err)
	}
	assert.False(t, ok, "baidu.com 应该无效")
}
