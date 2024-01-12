package scworkflow

import (
	"fmt"
	"uqtu/mediator/apiclient"
)

type Requester interface {
	GetRequest(endpoint string) (*apiclient.Request, error)
}

// implement the Requester interface

type credentials_requester struct {
	username string
	password string
	url      string
}

func (r credentials_requester) GetRequest(endpoint string) (*apiclient.Request, error) {
	if client := apiclient.NewClient(r.url, r.username, r.password, true); client == nil {
		return nil, fmt.Errorf("cannot get API client")
	} else {
		return client.NewGETwithBasicAuth(endpoint, "xml")
	}
}
