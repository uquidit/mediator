package totp

import (
	"encoding/base32"
	"errors"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/xlzd/gotp"
)

const (
	MS1 = "1234567890123456789012345678901234567890123456789012"
	MS2 = "ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	// These must be changed via compilation flags
	secretMS1 = MS1
	secretMS2 = MS2
)

// encode secrets (if need be) so they can be used by gotp
func init() {
	if secretMS1 == MS1 || secretMS2 == MS2 {
		logrus.Fatalln("tOTP secrets have not been changed. Stop.")
	}
	if err := checkSecrets(); err != nil {
		logrus.Println("Provided tOTP secrets are invalid. Trying to get them right...")
		encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
		secretMS1 = encoder.EncodeToString([]byte(secretMS1))
		secretMS2 = encoder.EncodeToString([]byte(secretMS2))
		if err := checkSecrets(); err != nil {
			logrus.Fatalf("I tried hard but tOTP secrets are still invalid: %v", err)
		}
	}
}

// check secrets are valid, ie properly encoded strings
// provided secrets are encoding in init so they should be valid
func checkSecrets() error {
	if !gotp.IsSecretValid(secretMS1) {
		return errors.New("secret 1 is invalid")
	}
	if !gotp.IsSecretValid(secretMS2) {
		return errors.New("secret 2 is invalid")
	}
	return nil
}

func CheckKey(key string, c echo.Context) (bool, error) {
	// sanity check: make sure secrets are properly encoded strings
	if err := checkSecrets(); err != nil {
		return false, err
	}

	totp1 := gotp.NewDefaultTOTP(secretMS1)
	totp2 := gotp.NewDefaultTOTP(secretMS2)
	now := time.Now()
	sec := now.Unix()
	return totp1.Verify(key[:6], sec) && totp2.Verify(key[6:], sec), nil
}

func GetKey() (string, error) {
	// sanity check: make sure secrets are properly encoded strings
	if err := checkSecrets(); err != nil {
		return "", err
	}

	totp1 := gotp.NewDefaultTOTP(secretMS1)
	totp2 := gotp.NewDefaultTOTP(secretMS2)
	return totp1.Now() + totp2.Now(), nil
}
