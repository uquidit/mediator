package mediatorscript

import (
	"crypto/hmac"
	"crypto/sha512"
	"io"
	"log"
	"os"
)

const (
	SALT      = "Some random string you must change"
	PEPPER    = "And this is another random string you must change too"
	SECRETKEY = "a secret key must be provided here"
)

var (
	// These 3 strings must be changed at build via compilation flag
	salt      = SALT
	pepper    = PEPPER
	secretKey = SECRETKEY
)

func init() {
	if salt == SALT {
		log.Fatalln("mediatorscript package salt has not been changed. Stop.")

	}
	if pepper == PEPPER {
		log.Fatalln("mediatorscript package pepper has not been changed. Stop.")

	}
	if secretKey == SECRETKEY {
		log.Fatalln("mediatorscript package secretKey has not been changed. Stop.")

	}
}

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
