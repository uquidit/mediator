//go:build !no_otp

package apiclient

import "mediator/totp"

func NewClientWithOTP(url_prefix string, insecure_skip_verify bool) (*Client, error) {

	client := NewClient(url_prefix, "", "", insecure_skip_verify)

	if err := client.SetToken(); err != nil {
		return nil, err
	}
	return client, nil
}

func NewClientWithOTPAndTimeout(url_prefix string, insecure_skip_verify bool, dial_timeout uint) (*Client, error) {

	client := NewClientWithDialTimeout(url_prefix, "", "", insecure_skip_verify, dial_timeout)

	if err := client.SetToken(); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) SetToken() error {
	var err error
	if c.Token, err = totp.GetKey(); err != nil {
		return err
	}
	return nil
}
