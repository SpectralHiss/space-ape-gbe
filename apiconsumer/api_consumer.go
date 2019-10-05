package apiconsumer

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (ae *GbeConsumer) ApiTraverse(path, queryString string) ([]json.RawMessage, error) {

	params, err := url.ParseQuery(queryString)
	if err != nil {
		return []json.RawMessage{}, fmt.Errorf("bad query")
	}
	params.Set("limit", "100")

	resp, err := http.Get(BuildQuery(ae.apiUrl, path, params))

	if err != nil {
		return []json.RawMessage{}, fmt.Errorf("The API endpoint is currently unreachable. %s", err.Error())
	}

	if resp.StatusCode == 401 {
		return []json.RawMessage{}, fmt.Errorf("Invalid API Key")
	}

	apiErrors := &APIErrors{}
	decoder := json.NewDecoder(resp.Body)

	var apiResponse1 = &APIResponse{}

	for decoder.More() {
		err := decoder.Decode(apiResponse1)

		if err != nil {
			return []json.RawMessage{}, err
		}

		// got a good first page response
		if apiResponse1.StatusCode == 1 {

			if apiResponse1.Limit != 100 {
				apiErrors.AddError("api limit changes, this is unexpected")
			}
			var results []json.RawMessage
			if len(apiResponse1.Results) > 0 {
				results = []json.RawMessage{apiResponse1.Results}
			}
			numPages := int(math.Ceil(float64(apiResponse1.NumberOfTotalResults) / float64(apiResponse1.Limit)))

			for i := 2; i <= numPages; i++ {

				params.Set("page", fmt.Sprintf("%d", i))
				fmt.Println(params.Get("page"))

				pageQuery := BuildQuery(ae.apiUrl, path, params)
				fmt.Println(pageQuery)

				pageResp, err := http.Get(pageQuery)
				if err != nil {
					apiErrors.AddError(fmt.Sprintf("Error fetching page %d", i))
				}

				var ParsedPageResponse = &APIResponse{}
				decoder := json.NewDecoder(pageResp.Body)
				err = decoder.Decode(ParsedPageResponse)

				if err != nil {
					return []json.RawMessage{}, err
				}

				results = append(results, ParsedPageResponse.Results)

			}

			return results, apiErrors.Niller()

		}
		apiErrors.AddError(apiResponse1.Error)

	}

	return []json.RawMessage{}, apiErrors.Niller()
}
