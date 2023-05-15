package apiclient

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type Client struct {
	client        *http.Client
	urlprefix     string
	username      string
	password      string
	Token         string
	UsernameField string
	PasswordField string
}

func NewClient(urlprefix string, username string, password string, InsecureSkipVerify bool) *Client {
	// get http client with certificate validation disabled
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: InsecureSkipVerify},
	}

	// remove any trailing '/' from urlprefix
	urlprefix = strings.TrimSuffix(urlprefix, "/")

	c := Client{
		client:        &http.Client{Transport: tr},
		urlprefix:     urlprefix,
		username:      username,
		password:      password,
		UsernameField: "username",
		PasswordField: "password",
	}
	return &c
}

/****** GET ******/

func (c *Client) NewGET(url string, content string) (*Request, error) {
	return c.newRequest("GET", url, content, AuthMode_None)
}

func (c *Client) NewGETwithToken(url string, content string) (*Request, error) {
	return c.newRequest("GET", url, content, AuthMode_Token)
}

func (c *Client) NewGETwithBodyAndToken(url string, body io.Reader, content string) (*Request, error) {
	return c.newWithBody("GET", url, body, content, AuthMode_Token)
}

func (c *Client) NewGETwithBasicAuth(url string, content string) (*Request, error) {
	return c.newRequest("GET", url, content, AuthMode_Basic)
}

/****** POST ******/

func (c *Client) NewPOST(url string, body io.Reader, content string) (*Request, error) {
	return c.newWithBody("POST", url, body, content, AuthMode_None)
}

func (c *Client) NewPOSTwithToken(url string, body io.Reader, content string) (*Request, error) {
	return c.newWithBody("POST", url, body, content, AuthMode_Token)
}

func (c *Client) NewPOSTwithBasicAuth(url string, body io.Reader, content string) (*Request, error) {
	return c.newWithBody("POST", url, body, content, AuthMode_Basic)
}

func (c *Client) NewPOSTwithFormData(url string, body io.Reader, content string) (*Request, error) {
	return c.newWithBody("POST", url, body, content, AuthMode_FormData)
}

/****** PUT ******/

func (c *Client) NewPUT(url string, body io.Reader, content string) (*Request, error) {
	return c.newWithBody("PUT", url, body, content, AuthMode_None)
}

func (c *Client) NewPUTwithToken(url string, body io.Reader, content string) (*Request, error) {
	return c.newWithBody("PUT", url, body, content, AuthMode_Token)
}

func (c *Client) NewPUTwithBasicAuth(url string, body io.Reader, content string) (*Request, error) {
	return c.newWithBody("PUT", url, body, content, AuthMode_Basic)
}

/****** DELETE ******/

func (c *Client) NewDELETE(url string, content string) (*Request, error) {
	return c.newRequest("DELETE", url, content, AuthMode_None)
}

func (c *Client) NewDELETEwithBodyAndToken(url string, body io.Reader, content string) (*Request, error) {
	return c.newWithBody("DELETE", url, body, content, AuthMode_Token)
}

func (c *Client) NewDELETEwithToken(url string, content string) (*Request, error) {
	return c.newRequest("DELETE", url, content, AuthMode_Token)
}

func (c *Client) NewDELETEwithBasicAuth(url string, content string) (*Request, error) {
	return c.newRequest("DELETE", url, content, AuthMode_Basic)
}

/****** private functions ******/

func (c *Client) newRequest(method string, url string, content string, auth_mode AuthenticationMode) (*Request, error) {
	return c.newWithBody(method, url, nil, content, auth_mode)
}

func (c *Client) newWithBody(method string, url_suffix string, body io.Reader, content string, auth_mode AuthenticationMode) (*Request, error) {
	var err error

	req := Request{
		client:  c.client,
		content: strings.ToLower(content),
	}

	if req.content != "xml" && req.content != "json" {
		return nil, fmt.Errorf("unsupported content type: %s", content)
	}

	url := c.urlprefix
	if url_suffix != "" {
		// remove any '/' from begining of url_suffix
		url_suffix = strings.TrimPrefix(url_suffix, "/")

		url = fmt.Sprintf("%s/%s", c.urlprefix, url_suffix)
	}

	if auth_mode == AuthMode_FormData {

		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		if err := writer.WriteField(c.UsernameField, c.username); err != nil {
			return nil, err
		}
		if err := writer.WriteField(c.PasswordField, c.password); err != nil {
			return nil, err
		}
		if err := writer.Close(); err != nil {
			return nil, err
		}

		if body != nil {
			r := io.MultiReader(payload, body)
			req.httpreq, err = http.NewRequest(method, url, r)
			if err != nil {
				return nil, fmt.Errorf("cannot create request for url %s: %v", url, err)
			}
		} else {
			req.httpreq, err = http.NewRequest(method, url, payload)
			if err != nil {
				return nil, fmt.Errorf("cannot create request for url %s: %v", url, err)
			}
		}
		req.SetHeader("Authorization", "Basic Og==")
		req.SetHeader("Content-Type", writer.FormDataContentType())
		if req.content == "json" {
			req.SetHeader("Accept", "application/json")
		} else {
			req.SetHeader("Accept", "application/xml")
		}

	} else {
		req.httpreq, err = http.NewRequest(method, url, body)
		if err != nil {
			return nil, fmt.Errorf("cannot create request for url %s: %v", url, err)
		}

		if req.content == "json" {
			req.SetHeader("Accept", "application/json")
			req.SetHeader("Content-Type", "application/json")
		} else {
			req.SetHeader("Accept", "application/xml")
			req.SetHeader("Content-Type", "application/xml")
		}
	}

	if auth_mode == AuthMode_Basic {
		if c.username == "" || c.password == "" {
			return nil, fmt.Errorf("no auth information")
		}
		req.httpreq.SetBasicAuth(c.username, c.password)
	} else if auth_mode == AuthMode_Token {
		if c.Token == "" {
			return nil, fmt.Errorf("token missing")
		}
		req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}

	return &req, nil
}
