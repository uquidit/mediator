package totp

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xlzd/gotp"
)

var (
	// These must be changed
	secretMS1 = "1234567890123456789012345678901234567890123456789012"
	secretMS2 = "ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func CheckKey(key string, c echo.Context) (bool, error) {
	totp1 := gotp.NewDefaultTOTP(secretMS1)
	totp2 := gotp.NewDefaultTOTP(secretMS2)
	now := time.Now()
	sec := now.Unix()
	return totp1.Verify(key[:6], sec) && totp2.Verify(key[6:], sec), nil
}

func GetKey() string {
	totp1 := gotp.NewDefaultTOTP(secretMS1)
	totp2 := gotp.NewDefaultTOTP(secretMS2)
	return totp1.Now() + totp2.Now()
}
