package apiclient

import (
	"fmt"
	"io"
)

var (
	last_request_status int
	backend_url         string
	ssl_skip_verify     bool
)

func InitHelpers(url string, skipVerify bool) {
	backend_url = url
	ssl_skip_verify = skipVerify
}

func GetLastRequestStatusCode() int {
	return last_request_status
}

func RunGETwithToken(url string, content string, v any) (io.ReadCloser, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { saveStatusCode(r) }()

	if client, err = NewClientWithOTP(backend_url, ssl_skip_verify); err != nil {
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

func RunPOSTwithToken(url string, body io.Reader, content string, v any) (io.ReadCloser, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { saveStatusCode(r) }()

	if client, err = NewClientWithOTP(backend_url, ssl_skip_verify); err != nil {
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

func RunDELETEwithToken(url string, content string, v any) (io.ReadCloser, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { saveStatusCode(r) }()

	if client, err = NewClientWithOTP(backend_url, ssl_skip_verify); err != nil {
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

func RunPOSTwithCredentials(url, u, p string, body io.Reader, content string, v any) (io.ReadCloser, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { saveStatusCode(r) }()

	client = NewClient(backend_url, u, p, ssl_skip_verify)

	if r, err = client.NewPOSTwithBasicAuth(url, body, content); err != nil {
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

func RunGETwithCredentials(url, u, p string, content string, v any) (io.ReadCloser, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { saveStatusCode(r) }()

	client = NewClient(backend_url, u, p, ssl_skip_verify)

	if r, err = client.NewGETwithBasicAuth(url, content); err != nil {
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

func saveStatusCode(r *Request) {
	if r != nil {
		last_request_status = r.StatusCode
	}
}

func getError(r io.ReadCloser, err error) error {
	if r == nil {
		return err
	}
	defer r.Close()
	data, _ := io.ReadAll(r)
	return fmt.Errorf("%w: %s", err, string(data))
}
