package main

import (
	"go-zrbc/pkg/utils"
	"strings"
	"testing"
)

func TestCrypt_respect(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "respect")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("respect2024")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}

func TestCrypt_flynt(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "flynt2024")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("Flynt841010")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}

func TestCrypt_joshua(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "joshua")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("Joshua2024")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}

func TestCrypt_patrick(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "patrick")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("Patrick2024")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}

func TestCrypt_laobaby(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "caocao")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("caocao2024new")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}

func TestCrypt_zrnathan(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "ZRnathan")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("ZRnathan2024")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}

func TestCrypt_perra(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "perra")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("perra2024")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}

func TestCrypt_zrboss(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "zrboss")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("zrboss2024")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}

func TestCrypt_zrfeifei(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "feifei")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("feifei2024")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}

func TestCrypt_zroper01(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "oper01")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("oper012024")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}

func TestCrypt_zrsita(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "sita")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("sita2024")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}

func TestCrypt_clonedy(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "clonedy")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("clonedy2024")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}

func TestCrypt_zoo(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "zoo")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("zoo2024")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}

func TestCrypt_lucy(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "lucy")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("lucy2024")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}

func TestCrypt_zoo1234(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "zoo1234")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("zoozoo123")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}

func TestCrypt_rowan(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "rowan")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("rowan2025")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}

func TestCrypt_tom(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "tomtom")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("tom2025")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}

func TestCrypt_jida(t *testing.T) {
	salt := utils.GenerateSalt("588326785867908888", "jida")
	t.Logf("salt:%s", salt)

	passwordFromWeb := utils.Md5("jida2025")
	passwordFromWebEx := strings.ToUpper(passwordFromWeb)
	t.Logf("passwordFromWebEx:%s", passwordFromWebEx)
	pw := utils.HashPassword(salt, passwordFromWebEx)

	t.Logf("pw:%s", pw)
}
