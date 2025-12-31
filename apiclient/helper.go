package apiclient

import (
	"errors"
	"io"
	"net/http"
)

type APIclientHelper struct {
	last_request_status  int
	last_request_session string
	backend_url          string
	ssl_skip_verify      bool
	time_out             uint
	response             *http.Response
	cookie               string
}
type QueryParams map[string]string

func GetHelper(url string, skipVerify bool) *APIclientHelper {
	h := APIclientHelper{
		backend_url:     url,
		ssl_skip_verify: skipVerify,
	}
	return &h
}

func GetHelperWithCookie(url string, skipVerify bool, cookie string) *APIclientHelper {
	h := APIclientHelper{
		backend_url:     url,
		ssl_skip_verify: skipVerify,
		cookie:          cookie,
	}
	return &h
}

func GetHelperWithTimeout(url string, skipVerify bool, timeout uint) *APIclientHelper {
	h := APIclientHelper{
		backend_url:     url,
		ssl_skip_verify: skipVerify,
		time_out:        timeout,
	}
	return &h
}

func (h *APIclientHelper) GetLastRequestStatusCode() int {
	if h.last_request_status == 0 {
		return http.StatusInternalServerError
	}
	return h.last_request_status
}

func (h *APIclientHelper) GetLastRequestSessionID() string {
	return h.last_request_session
}

func (h *APIclientHelper) GetIDFromLocationHeader() (int, error) {
	// try and get app new ID
	if h.response == nil {
		return 0, errors.New("response is empty. Run request first")
	}
	return getIDFromLocationHeader(h.response.Header)
}

func (h *APIclientHelper) RunPOSTwithCredentials(url, u, p string, body io.Reader, content string, v any) (io.Reader, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { h.saveStatusCodeAndSession(r) }()

	if h.time_out != 0 {
		client = NewClientWithDialTimeout(h.backend_url, u, p, h.ssl_skip_verify, h.time_out)
	} else {
		client = NewClient(h.backend_url, u, p, h.ssl_skip_verify)
	}

	if r, err = client.NewPOSTwithBasicAuth(url, body, content); err != nil {
		return nil, err
	} else if v == nil {
		// run without decode
		return r.RunWithoutDecode()
	} else {
		return nil, r.Run(v)
	}
}

func (h *APIclientHelper) RunDELETEwithCredentials(url, u, p string, content string, v any) (io.Reader, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { h.saveStatusCodeAndSession(r) }()

	if h.time_out != 0 {
		client = NewClientWithDialTimeout(h.backend_url, u, p, h.ssl_skip_verify, h.time_out)
	} else {
		client = NewClient(h.backend_url, u, p, h.ssl_skip_verify)
	}

	if r, err = client.NewDELETEwithBasicAuth(url, content); err != nil {
		return nil, err
	} else if v == nil {
		// run without decode
		return r.RunWithoutDecode()
	} else {
		return nil, r.Run(v)
	}
}

func (h *APIclientHelper) RunGETwithCredentials(url, u, p string, content string, v any) (io.Reader, error) {
	return h.RunGETwithCredentialsAndParams(url, u, p, nil, content, v)
}

func (h *APIclientHelper) RunGETwithCredentialsAndParams(url, u, p string, params QueryParams, content string, v any) (io.Reader, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { h.saveStatusCodeAndSession(r) }()

	if h.time_out != 0 {
		client = NewClientWithDialTimeout(h.backend_url, u, p, h.ssl_skip_verify, h.time_out)
	} else {
		client = NewClient(h.backend_url, u, p, h.ssl_skip_verify)
	}

	if r, err = client.NewGETwithBasicAuth(url, content); err != nil {
		return nil, err
	}

	r.AddQueryParams(params)

	if v == nil {
		// run without decode
		return r.RunWithoutDecode()
	} else {
		return nil, r.Run(v)
	}
}

func (h *APIclientHelper) saveStatusCodeAndSession(r *Request) {
	if r != nil {
		h.last_request_status = r.StatusCode
		h.last_request_session = r.JsessionID
		h.response = r.response
	}
}

func (h *APIclientHelper) RunGETwithCookie(url string, content string, v any) (io.Reader, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { h.saveStatusCodeAndSession(r) }()

	if h.time_out != 0 {
		client = NewClientWithDialTimeout(h.backend_url, "", "", h.ssl_skip_verify, h.time_out)
	} else {
		client = NewClient(h.backend_url, "", "", h.ssl_skip_verify)
	}

	client.Token = h.cookie

	if r, err = client.NewGETwithCookie(url, content); err != nil {
		return nil, err
	}

	if v == nil {
		// run without decode
		return r.RunWithoutDecode()
	} else {
		return nil, r.Run(v)
	}
}
