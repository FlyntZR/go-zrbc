package main

import (
	"go-zrbc/pkg/utils"
	"testing"
)

func TestAuth_bib(t *testing.T) {
	sig := utils.HmacSha256("decc228dd39d1324f159aefb5121a7512.0MD51703741966000", "97e280b444b3fdba554e9167f1407bd1")
	t.Logf("sig : %s\n", sig)
}
