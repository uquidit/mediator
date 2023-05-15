package apiclient

type AuthenticationMode int

const (
	AuthMode_None AuthenticationMode = iota //this will be default value as first call to iota returns 0
	AuthMode_Basic
	AuthMode_FormData
	AuthMode_Token
)
