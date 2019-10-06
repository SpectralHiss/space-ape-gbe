package search

import (
	"encoding/json"
	"net/url"

	"github.com/SpectralHiss/space-ape-gbe/apiconsumer"
)

type SearchResponse struct {
	Name        string `json:"name"`
	Deck        string `json:"deck"`
	ReleaseDate int    `json:"expected_release_year"`
	Platforms   []struct {
		Name string
	} `json:"platforms"`

	ID int `json:"id"`
}

type Searcher struct {
	Url    string
	ApiKey string
}

func NewSearcher(url string, apiK string) *Searcher {
	return &Searcher{Url: url, ApiKey: apiK}
}

func (s *Searcher) SearchTitles(titleQuery string) ([]SearchResponse, error) {
	consumer := apiconsumer.NewGBE(s.Url)

	params := url.Values{}
	params.Set("limit", "100")
	params.Set("resources", "game")
	params.Set("query", titleQuery)
	params.Set("api_key", s.ApiKey)
	params.Set("format", "json")
	rawJSONs, err := consumer.ApiTraverse("api/search", params.Encode())

	//fmt.Println(string(rawJSONs[0]))

	if err != nil {
		return []SearchResponse{}, err
	}
	allResp := []SearchResponse{}
	for _, rawJSON := range rawJSONs {
		// decode json
		searchResps := []SearchResponse{}
		json.Unmarshal([]byte(rawJSON), &searchResps)

		allResp = append(allResp, searchResps...)
	}

	return allResp, nil
}
