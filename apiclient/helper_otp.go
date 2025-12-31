//go:build !no_otp

package apiclient

import (
	"io"
)

func (h *APIclientHelper) RunGETwithToken(url string, content string, v any) (io.Reader, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { h.saveStatusCodeAndSession(r) }()

	if h.time_out != 0 {
		if client, err = NewClientWithOTPAndTimeout(h.backend_url, h.ssl_skip_verify, h.time_out); err != nil {
			return nil, err
		}
	} else {
		if client, err = NewClientWithOTP(h.backend_url, h.ssl_skip_verify); err != nil {
			return nil, err
		}
	}

	if r, err = client.NewGETwithToken(url, content); err != nil {
		return nil, err
	}

	if v == nil {
		// run without decode
		return r.RunWithoutDecode()
	} else {
		return nil, r.Run(v)
	}
}

func (h *APIclientHelper) RunGETwithBodyAndToken(url string, body io.Reader, content string, v any) (io.Reader, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { h.saveStatusCodeAndSession(r) }()

	if h.time_out != 0 {
		if client, err = NewClientWithOTPAndTimeout(h.backend_url, h.ssl_skip_verify, h.time_out); err != nil {
			return nil, err
		}
	} else {
		if client, err = NewClientWithOTP(h.backend_url, h.ssl_skip_verify); err != nil {
			return nil, err
		}
	}

	if r, err = client.NewGETwithBodyAndToken(url, body, content); err != nil {
		return nil, err
	} else if v == nil {
		// run without decode
		return r.RunWithoutDecode()
	} else {
		return nil, r.Run(v)
	}
}

func (h *APIclientHelper) RunPOSTwithToken(url string, body io.Reader, content string, v any) (io.Reader, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { h.saveStatusCodeAndSession(r) }()

	if h.time_out != 0 {
		if client, err = NewClientWithOTPAndTimeout(h.backend_url, h.ssl_skip_verify, h.time_out); err != nil {
			return nil, err
		}
	} else {
		if client, err = NewClientWithOTP(h.backend_url, h.ssl_skip_verify); err != nil {
			return nil, err
		}
	}

	if r, err = client.NewPOSTwithToken(url, body, content); err != nil {
		return nil, err
	} else if v == nil {
		// run without decode
		return r.RunWithoutDecode()
	} else {
		return nil, r.Run(v)
	}
}

func (h *APIclientHelper) RunPUTwithToken(url string, body io.Reader, content string, v any) (io.Reader, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { h.saveStatusCodeAndSession(r) }()

	if h.time_out != 0 {
		if client, err = NewClientWithOTPAndTimeout(h.backend_url, h.ssl_skip_verify, h.time_out); err != nil {
			return nil, err
		}
	} else {
		if client, err = NewClientWithOTP(h.backend_url, h.ssl_skip_verify); err != nil {
			return nil, err
		}
	}

	if r, err = client.NewPUTwithToken(url, body, content); err != nil {
		return nil, err
	} else if v == nil {
		// run without decode
		return r.RunWithoutDecode()
	} else {
		return nil, r.Run(v)
	}
}

func (h *APIclientHelper) RunDELETEwithToken(url string, content string, v any) (io.Reader, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { h.saveStatusCodeAndSession(r) }()

	if h.time_out != 0 {
		if client, err = NewClientWithOTPAndTimeout(h.backend_url, h.ssl_skip_verify, h.time_out); err != nil {
			return nil, err
		}
	} else {
		if client, err = NewClientWithOTP(h.backend_url, h.ssl_skip_verify); err != nil {
			return nil, err
		}
	}

	if r, err = client.NewDELETEwithToken(url, content); err != nil {
		return nil, err
	} else if v == nil {
		// run without decode
		return r.RunWithoutDecode()
	} else {
		return nil, r.Run(v)
	}
}

func (h *APIclientHelper) RunDELETEwithBodyAndToken(url string, body io.Reader, content string, v any) (io.Reader, error) {
	var (
		err    error
		r      *Request
		client *Client
	)
	defer func() { h.saveStatusCodeAndSession(r) }()

	if h.time_out != 0 {
		if client, err = NewClientWithOTPAndTimeout(h.backend_url, h.ssl_skip_verify, h.time_out); err != nil {
			return nil, err
		}
	} else {
		if client, err = NewClientWithOTP(h.backend_url, h.ssl_skip_verify); err != nil {
			return nil, err
		}
	}

	if r, err = client.NewDELETEwithBodyAndToken(url, body, content); err != nil {
		return nil, err
	} else if v == nil {
		// run without decode
		return r.RunWithoutDecode()
	} else {
		return nil, r.Run(v)
	}
}
