package clicommands

import "io"

type APIClient interface {
	RunGETwithToken(url string, content_type string, v any) (io.ReadCloser, error)
	RunPOSTwithToken(url string, body io.Reader, content_type string, v any) (io.ReadCloser, error)
	RunDELETEwithToken(url string, content_type string, v any) (io.ReadCloser, error)
	GetLastRequestStatusCode() int
}
