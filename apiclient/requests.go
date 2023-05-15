package apiclient

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
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

	loc := req.response.Header.Get("Location")
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
	if _, err := req.RunWithoutDecode(); err != nil {
		// rawBody, _ := ioutil.ReadAll(b)
		// fmt.Println(string(rawBody))
		return err
	}

	if v != nil {
		var err error

		// let's see if body is meaningful
		if req.content == "json" {
			// rawBody, _ := ioutil.ReadAll(req.response.Body)
			// fmt.Println(string(rawBody))
			// req.response.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
			err = json.NewDecoder(req.response.Body).Decode(v)
		} else {
			// rawBody, _ := ioutil.ReadAll(req.response.Body)
			// fmt.Println(string(rawBody))
			// req.response.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
			err = xml.NewDecoder(req.response.Body).Decode(v)
		}
		if err != nil {
			return fmt.Errorf("cannot decode %s data: %v", req.content, err)
		}
	}

	return nil
}

func (req *Request) RunWithoutDecode() (io.ReadCloser, error) {
	var err error
	req.response, err = req.client.Do(req.httpreq)
	if err != nil {
		return nil, fmt.Errorf("error while running request %s: %v", req.httpreq.URL, err)
	}

	req.StatusCode = req.response.StatusCode
	switch req.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent:
		//get JSESSIONID cookie
		for _, cookie := range req.response.Cookies() {
			if cookie.Name == "JSESSIONID" {
				req.JsessionID = cookie.Value
			}
		}

		return req.response.Body, nil

	default:
		return req.response.Body, req.decodeError()

	}
}

func (req *Request) decodeError() error {
	//duplicate request body in case caller needs it
	// cf https://stackoverflow.com/questions/43021058/golang-read-request-body-multiple-times
	rawBody, err := io.ReadAll(req.response.Body)
	if err != nil {
		// can't read body
		// too bad: return status code error
		// further attempts to read req body will return EOF.
		// guess it's fair enough
		return fmt.Errorf("%s", http.StatusText(req.StatusCode))
	}
	req.response.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	// let's see if body is meaningful
	// can we unmarshall body in an error struct.
	// this is very specific to uQuidIT API
	uqt_err := struct {
		Error string `json:"error"`
	}{}
	if errr := json.Unmarshal(rawBody, &uqt_err); errr == nil {
		return fmt.Errorf("%s: %s", http.StatusText(req.StatusCode), uqt_err.Error)
	} else {
		// can't unmarshall: return raw body
		return fmt.Errorf("%s: %s", http.StatusText(req.StatusCode), string(rawBody))
	}
}

func (req *Request) AddQueryParams(key, value string) {
	q := req.httpreq.URL.Query()          // Get a copy of the query values.
	q.Add(key, value)                     // Add a new value to the set.
	req.httpreq.URL.RawQuery = q.Encode() // Encode and assign back to the original query.
}

func (req *Request) String() string {
	data, _ := io.ReadAll(req.response.Body)
	return string(data)
}
