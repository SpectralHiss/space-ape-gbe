package apiconsumer

import (
	"bytes"
	"net/http"
)

type GbeConsumer struct {
	apiUrl string
}

func NewGBE(url string) *GbeConsumer {
	return &GbeConsumer{apiUrl: url}
}

type APIErrors struct {
	retErr []string
}

func (ae *APIErrors) Error() string {
	var out bytes.Buffer

	for _, err := range ae.retErr {
		out.WriteString(err)
	}
	return out.String()
}

func (ae *APIErrors) AddError(errString string) {
	ae.retErr = append(ae.retErr, errString)
}

func (*GbeConsumer) ApiTraverse(query *http.Request) ([]string, error) {

	//resp, err := http.Get(query.)
	return []string{}, nil
}
