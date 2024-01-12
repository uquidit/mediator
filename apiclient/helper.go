package apiclient

import (
	"io"
	"net/http"
)

type APIclientHelper struct {
	last_request_status int
	backend_url         string
	ssl_skip_verify     bool
}
type QueryParams map[string]string

func GetHelper(url string, skipVerify bool) *APIclientHelper {
	h := APIclientHelper{
		backend_url:     url,
		ssl_skip_verify: skipVerify,
	}
	return &h
}

func (h *APIclientHelper) GetLastRequestStatusCode() int {
	if h.last_request_status == 0 {
		return http.StatusInternalServerError
	}
	return h.last_request_status
}

func (h *APIclientHelper) RunGETwithToken(url string, content string, v any) (io.ReadCloser, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { h.saveStatusCode(r) }()

	if client, err = NewClientWithOTP(h.backend_url, h.ssl_skip_verify); err != nil {
		return nil, err
	}

	if r, err = client.NewGETwithToken(url, content); err != nil {
		return nil, err
	}

	if v == nil {
		// run without decode unless we get an error
		if buff, err := r.RunWithoutDecode(); err != nil {
			return nil, getError(buff, err)
		} else {
			return buff, err
		}
	} else {
		return nil, r.Run(v)
	}
}

func (h *APIclientHelper) RunPOSTwithToken(url string, body io.Reader, content string, v any) (io.ReadCloser, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { h.saveStatusCode(r) }()

	if client, err = NewClientWithOTP(h.backend_url, h.ssl_skip_verify); err != nil {
		return nil, err
	}

	if r, err = client.NewPOSTwithToken(url, body, content); err != nil {
		return nil, err
	} else if v == nil {
		// run without decode unless we get an error
		if buff, err := r.RunWithoutDecode(); err != nil {
			return nil, getError(buff, err)
		} else {
			return buff, err
		}
	} else {
		return nil, r.Run(v)
	}
}

func (h *APIclientHelper) RunDELETEwithToken(url string, content string, v any) (io.ReadCloser, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { h.saveStatusCode(r) }()

	if client, err = NewClientWithOTP(h.backend_url, h.ssl_skip_verify); err != nil {
		return nil, err
	}

	if r, err = client.NewDELETEwithToken(url, content); err != nil {
		return nil, err
	} else if v == nil {
		// run without decode unless we get an error
		if buff, err := r.RunWithoutDecode(); err != nil {
			return nil, getError(buff, err)
		} else {
			return buff, err
		}
	} else {
		return nil, r.Run(v)
	}
}

func (h *APIclientHelper) RunPOSTwithCredentials(url, u, p string, body io.Reader, content string, v any) (io.ReadCloser, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { h.saveStatusCode(r) }()

	client = NewClient(h.backend_url, u, p, h.ssl_skip_verify)

	if r, err = client.NewPOSTwithBasicAuth(url, body, content); err != nil {
		return nil, err
	} else if v == nil {
		// run without decode unless we get an error
		// if buff, err := r.RunWithoutDecode(); err != nil {
		// 	return nil, getError(buff, err)
		// } else {
		// 	return buff, err
		// }
		return r.RunWithoutDecode()
	} else {
		return nil, r.Run(v)
	}
}

func (h *APIclientHelper) RunGETwithCredentials(url, u, p string, content string, v any) (io.ReadCloser, error) {
	return h.RunGETwithCredentialsAndParams(url, u, p, nil, content, v)
}

func (h *APIclientHelper) RunGETwithCredentialsAndParams(url, u, p string, params QueryParams, content string, v any) (io.ReadCloser, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { h.saveStatusCode(r) }()

	client = NewClient(h.backend_url, u, p, h.ssl_skip_verify)

	if r, err = client.NewGETwithBasicAuth(url, content); err != nil {
		return nil, err
	}

	r.AddQueryParams(params)

	if v == nil {
		// run without decode unless we get an error
		// if buff, err := r.RunWithoutDecode(); err != nil {
		// 	return nil, getError(buff, err)
		// } else {
		// 	return buff, err
		// }
		return r.RunWithoutDecode()

	} else {
		return nil, r.Run(v)
	}
}

func (h *APIclientHelper) saveStatusCode(r *Request) {
	if r != nil {
		h.last_request_status = r.StatusCode
	}
}
