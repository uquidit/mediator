package apiclient

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
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

var (
	http_client_verify   map[uint]*http.Client
	http_client_noverify map[uint]*http.Client
	mutex                sync.Mutex
)

func init() {
	http_client_verify = make(map[uint]*http.Client)
	http_client_noverify = make(map[uint]*http.Client)
}

func newClient(dial_timeout uint, InsecureSkipVerify bool) *http.Client {
	mutex.Lock()
	defer mutex.Unlock()

	if InsecureSkipVerify {
		if c, ok := http_client_noverify[dial_timeout]; ok {
			return c
		}
	} else {
		if c, ok := http_client_verify[dial_timeout]; ok {
			return c
		}

	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: InsecureSkipVerify},
		Dial: (&net.Dialer{
			Timeout: time.Duration(dial_timeout) * time.Second,
		}).Dial,
	}
	http_client := &http.Client{Transport: tr}

	if InsecureSkipVerify {
		http_client_noverify[dial_timeout] = http_client
	} else {
		http_client_verify[dial_timeout] = http_client
	}
	return http_client
}

func NewClientWithDialTimeout(urlprefix string, username string, password string, InsecureSkipVerify bool, dial_timeout uint) *Client {
	// get http client
	http_client := newClient(dial_timeout, InsecureSkipVerify)

	return newRichClientFromHTTPClient(urlprefix, username, password, http_client)
}

func NewClient(urlprefix string, username string, password string, InsecureSkipVerify bool) *Client {
	// get http client with fixed timeout
	http_client := newClient(10, InsecureSkipVerify)

	return newRichClientFromHTTPClient(urlprefix, username, password, http_client)
}

func newRichClientFromHTTPClient(urlprefix string, username string, password string, http_client *http.Client) *Client {
	// remove any trailing '/' from urlprefix
	urlprefix = strings.TrimSuffix(urlprefix, "/")

	c := Client{
		client:        http_client,
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

func (c *Client) NewGETwithCookie(url string, content string) (*Request, error) {
	return c.newRequest("GET", url, content, AuthMode_Cookie)
}

/****** POST ******/

func (c *Client) NewPOST(url string, body io.Reader, content string) (*Request, error) {
	return c.newWithBody("POST", url, body, content, AuthMode_None)
}

func (c *Client) NewPOSTwithToken(url string, body io.Reader, content string) (*Request, error) {
	return c.newWithBody("POST", url, body, content, AuthMode_Token)
}
func (c *Client) NewPOSTwithCookie(url string, body io.Reader, content string) (*Request, error) {
	return c.newWithBody("POST", url, body, content, AuthMode_Cookie)
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
			// do NOT set Content-Type header for XML request
			// Securetrack rejects it for some reason
			// cf https://gitlab.uquidit.corp/uqtu/mediator/-/issues/14 and https://gitlab.uquidit.corp/uqtu/back-end/-/issues/123
			// req.SetHeader("Content-Type", "application/xml")
		}
	}

	switch auth_mode {
	case AuthMode_Basic:
		if c.username == "" || c.password == "" {
			return nil, fmt.Errorf("no auth information")
		}
		req.httpreq.SetBasicAuth(c.username, c.password)

	case AuthMode_Token:
		if c.Token == "" {
			return nil, fmt.Errorf("token missing")
		}
		req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.Token))
		c.Token = ""
	case AuthMode_Cookie:
		if c.Token == "" {
			return nil, fmt.Errorf("token missing")
		}
		cookie := http.Cookie{
			Name:  "JSESSIONID",
			Value: c.Token,
		}

		req.AddCookie(&cookie)
	}

	return &req, nil
}
