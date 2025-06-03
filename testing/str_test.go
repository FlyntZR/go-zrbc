package main

import (
	"go-zrbc/pkg/utils"
	"testing"
)

func TestIsContainsSpecial(t *testing.T) {
	result := utils.IsContainsSpecial("'骗骗,喜欢你")
	// result := util.RemoveSpecail("Hello!@# World123")
	t.Logf("res:%v", result)
}

func TestRemoveSpecail(t *testing.T) {
	result := utils.RemoveSpecail("'骗骗?，：（,喜欢你")
	// result := util.RemoveSpecail("Hello!@# World123")
	t.Logf("res:%s", result)
}
