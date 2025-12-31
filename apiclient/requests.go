package apiclient

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Request struct {
	httpreq    *http.Request
	response   *http.Response
	StatusCode int
	JsessionID string
	client     *http.Client
	content    string
}

func (req *Request) AddCookie(cookie *http.Cookie) {
	req.httpreq.AddCookie(cookie)
}

func (req *Request) SetHeader(key string, value string) {
	req.httpreq.Header.Set(key, value)
}

func (req *Request) GetAllHeaders() http.Header {
	return req.httpreq.Header
}

func (req *Request) GetIDFromLocationHeader() (int, error) {
	// try and get app new ID
	if req.response == nil {
		return 0, fmt.Errorf("response is empty. Run request first")
	}
	return getIDFromLocationHeader(req.response.Header)
}

func getIDFromLocationHeader(header http.Header) (int, error) {
	loc := header.Get("Location")
	if loc == "" {
		return 0, fmt.Errorf("location header is empty")
	}
	parts := strings.Split(loc, "/")
	if len(parts) == 0 {
		return 0, fmt.Errorf("location header is ill-formed")
	}
	id := parts[len(parts)-1]
	if id, err := strconv.Atoi(id); err != nil {
		return 0, err
	} else {
		return id, nil
	}
}

func (req *Request) Run(v any) error {
	if resp_body, err := req.RunWithoutDecode(); err != nil {
		return err

	} else if v != nil {
		var err error

		// let's see if body is meaningful
		if req.content == "json" {
			err = json.NewDecoder(resp_body).Decode(v)
		} else {
			err = xml.NewDecoder(resp_body).Decode(v)
		}
		if err != nil {
			return fmt.Errorf("cannot decode %s data: %v", req.content, err)
		}
	}

	return nil
}

func (req *Request) RunWithoutDecode() (io.Reader, error) {
	var err error
	req.response, err = req.client.Do(req.httpreq)
	defer func() {
		if req.response != nil && req.response.Body != nil {
			req.response.Body.Close()
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("error while running request %s: %w", req.httpreq.URL, err)
	}

	response, err := io.ReadAll(req.response.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading request %s response: %w", req.httpreq.URL, err)
	}
	req.StatusCode = req.response.StatusCode
	switch {
	case req.StatusCode >= 200 && req.StatusCode < 300:
		//get JSESSIONID cookie
		for _, cookie := range req.response.Cookies() {
			if cookie.Name == "JSESSIONID" {
				req.JsessionID = cookie.Value
			}
		}

		return bytes.NewReader(response), nil

	default:
		return bytes.NewReader(response), req.decodeError(response)
	}
}

func (req *Request) decodeError(resp_body []byte) error {
	// let's see if body is meaningful
	// can we unmarshall body in an error struct.
	// this is very specific to uQuidIT API
	uqt_err := struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}{}
	if errr := json.Unmarshal(resp_body, &uqt_err); errr != nil {
		// can't unmarshall: return status code error
		return req.getErrorFromStatus()

	} else if uqt_err.Message != "" {
		return errors.New(uqt_err.Message)
	} else if uqt_err.Error != "" {
		return fmt.Errorf("%w: %s", req.getErrorFromStatus(), uqt_err.Error)
	} else {
		// no message
		return req.getErrorFromStatus()
	}
}

// Returns an error using request status code text
func (req *Request) getErrorFromStatus() error {
	return errors.New(http.StatusText(req.StatusCode))
}

func (req *Request) AddQueryParam(key, value string) {
	q := req.httpreq.URL.Query()          // Get a copy of the query values.
	q.Add(key, value)                     // Add a new value to the set.
	req.httpreq.URL.RawQuery = q.Encode() // Encode and assign back to the original query.
}

func (req *Request) AddQueryParams(params QueryParams) {
	if len(params) > 0 {
		q := req.httpreq.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.httpreq.URL.RawQuery = q.Encode()
	}
}

func (req *Request) GetQueryLength() int {
	if req == nil || req.httpreq == nil {
		return 0
	}
	return len(req.httpreq.URL.String())
}

func (req *Request) String() string {
	data, _ := io.ReadAll(req.response.Body)
	return string(data)
}
