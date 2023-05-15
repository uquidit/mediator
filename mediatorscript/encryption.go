package mediatorscript

import (
	"crypto/hmac"
	"crypto/sha512"
	"io"
	"os"
)

var (
	// These 3 strings must be changed
	salt      = "Some random string you must change"
	pepper    = "And this is another random string you must change too"
	secretKey = "a secret key must be provided here"
)

func (s *Script) computeHash() ([]byte, error) {
	if file, err := os.Open(s.Fullpath); err != nil {
		return nil, err
	} else if content, err := io.ReadAll(file); err != nil {
		return nil, err
	} else {
		h := hmac.New(sha512.New, []byte(secretKey))
		h.Write([]byte(salt))
		h.Write(content)
		h.Write([]byte(pepper))
		return h.Sum(nil), nil
	}
}
