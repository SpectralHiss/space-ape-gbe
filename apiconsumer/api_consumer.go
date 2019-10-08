package apiconsumer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
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

func (ae *APIErrors) Niller() error {
	if len(ae.retErr) == 0 {
		return nil
	} else {
		return ae
	}
}

func isEmptyArray(data json.RawMessage) bool {
	// ugly but cheap
	return bytes.Equal([]byte(data), []byte("[]"))
}

type APIResponse struct {
	StatusCode           int `json:"status_code"`
	Error                string
	NumberOfTotalResults int `json:"number_of_total_results"`
	NumberOfPageResults  int `json:"number_of_page_results"`
	Limit                int
	Offset               int
	Results              json.RawMessage
}

func BuildQuery(url string, path string, params url.Values) string {
	return fmt.Sprintf("%s/%s?%s", url, path, params.Encode())
}

func debugBody(r io.ReadCloser) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s", data)
}

func (ae *GbeConsumer) ApiTraverse(path, queryString string) ([]json.RawMessage, error) {

	params, err := url.ParseQuery(queryString)
	if err != nil {
		return []json.RawMessage{}, fmt.Errorf("bad query")
	}

	resp, err := http.Get(BuildQuery(ae.apiUrl, path, params))

	if err != nil {
		return []json.RawMessage{}, fmt.Errorf("The API endpoint is currently unreachable. %s", err.Error())
	}

	if resp.StatusCode == 401 {
		return []json.RawMessage{}, fmt.Errorf("Invalid API Key")
	}

	apiErrors := &APIErrors{}

	decoder := json.NewDecoder(resp.Body)

	var apiResponseP1 = APIResponse{}

	err = decoder.Decode(&apiResponseP1)

	if err != nil {
		return []json.RawMessage{}, fmt.Errorf("Error decoding json response: %s", err.Error())
	}

	// got a good first page response
	if apiResponseP1.StatusCode == 1 {

		// TODO 0 == 2
		var results = []json.RawMessage{}
		//fmt.Printf("%s", apiResponse1.Results)
		if !isEmptyArray(apiResponseP1.Results) {
			results = []json.RawMessage{apiResponseP1.Results}
		}

		numPages := int(math.Ceil(float64(apiResponseP1.NumberOfTotalResults) / float64(apiResponseP1.Limit)))

		for i := 2; i <= numPages; i++ {
			params.Set("page", fmt.Sprintf("%d", i))
			parsedResponse, err := ae.fetchPage(path, params)
			if err != nil {
				apiErrors.AddError(err.Error())
			}
			results = append(results, parsedResponse.Results)

		}

		return results, apiErrors.Niller()
	}
	apiErrors.AddError(apiResponseP1.Error)

	return []json.RawMessage{}, apiErrors.Niller()
}

func (ae *GbeConsumer) fetchPage(path string, params url.Values) (APIResponse, error) {
	//fmt.Println(params.Get("page"))

	pageQuery := BuildQuery(ae.apiUrl, path, params)
	//fmt.Println(pageQuery)

	pageResp, err := http.Get(pageQuery)
	if err != nil {
		return APIResponse{}, fmt.Errorf("Error fetching page %s", params.Get("page"))
	}

	var ParsedPageResponse = APIResponse{}
	decoder := json.NewDecoder(pageResp.Body)
	err = decoder.Decode(&ParsedPageResponse)

	if err != nil {
		return APIResponse{}, err
	}

	return ParsedPageResponse, nil
}
