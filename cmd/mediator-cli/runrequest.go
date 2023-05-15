package main

import (
	"fmt"
	"io"
	"uqtu/mediator/apiclient"
	"uqtu/mediator/totp"
)

// this is a wrapper around the apiclient.Client struct
// we need it here to add extra method to apiclient.Client
// to it fits clicommands interface.
type ClientRequester struct {
	client *apiclient.Client
}

func GetClientRequester(url string, skipVerify bool) *ClientRequester {
	if len(url) == 0 {
		return nil // should return an error here
	}
	cr := ClientRequester{}
	cr.client = apiclient.NewClient(url, "", "", skipVerify)
	return &cr
}

func (cr ClientRequester) RunGETwithToken(url string, content string, v any) (io.ReadCloser, error) {
	if cr.client == nil {
		return nil, fmt.Errorf("no API client")
	}
	cr.client.Token = totp.GetKey()

	if r, err := cr.client.NewGETwithToken(url, content); err != nil {
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

func (cr *ClientRequester) RunPOSTwithToken(url string, body io.Reader, content string, v any) (io.ReadCloser, error) {
	if cr.client == nil {
		return nil, fmt.Errorf("no API client")
	}
	cr.client.Token = totp.GetKey()

	if r, err := cr.client.NewPOSTwithToken(url, body, content); err != nil {
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

func (cr *ClientRequester) RunDELETEwithToken(url string, content string, v any) (io.ReadCloser, error) {
	if cr.client == nil {
		return nil, fmt.Errorf("no API client")
	}
	cr.client.Token = totp.GetKey()

	if r, err := cr.client.NewDELETEwithToken(url, content); err != nil {
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

func getError(r io.ReadCloser, err error) error {
	defer r.Close()
	data, _ := io.ReadAll(r)
	return fmt.Errorf("%w: %s", err, string(data))
}
