package fetch

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"time"

	"sync"

	"github.com/SpectralHiss/space-ape-gbe/apiconsumer"
)

type NameEntry struct {
	Name string `json:"name"`
}

type GameResponse struct {
	Name        string      `json:"name"`
	Deck        string      `json:"deck"`
	ReleaseDate int         `json:"expected_release_year"`
	ID          int         `json:"id"`
	Description string      `json:"description"`
	Genres      []NameEntry `json:"genres,omitempty"`
	Platforms   []NameEntry `json:"platforms,omitempty"`
	GameRating  []NameEntry `json:"original_game_rating,omitempty"`
	Characters  []NameEntry `json:"characters,omitempty"`
	Publishers  []NameEntry `json:"publishers,omitempty"`
	Developers  []NameEntry `json:"developers,omitempty"`
	Franchises  []NameEntry `json:"franchises,omitempty"`
}

type DLCDetails struct {
	ID          int
	Name        string
	ReleaseDate string `json:"release_date"`
}

type ByReleaseD []DLCDetails

func (a ByReleaseD) Len() int {
	return len(a)
}

func (a ByReleaseD) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]

}

func (a ByReleaseD) Less(i, j int) bool {
	if a[i].ReleaseDate == "" {
		return true
	}

	if a[j].ReleaseDate == "" {
		return false
	}

	d1, err := time.Parse("2006-01-02 15:04:05", a[i].ReleaseDate)
	if err != nil {
		panic(err)
	}
	d2, err := time.Parse("2006-01-02 15:04:05", a[j].ReleaseDate)
	if err != nil {
		panic(err)
	}
	return d1.Before(d2)

}

type GameResponseDLCs struct {
	GameResponse
	DLCs []DLCDetails
}

type Fetcher struct {
	Url      string
	ApiKey   string
	consumer *apiconsumer.GbeConsumer
}

func NewFetcher(url string, apiK string) *Fetcher {
	return &Fetcher{Url: url, ApiKey: apiK}
}

func (fetcher *Fetcher) RawFetch(guid string) ([]json.RawMessage, error) {
	fetcher.consumer = apiconsumer.NewGBE(fetcher.Url)

	params := url.Values{}
	//params.Set("limit", "100")
	params.Set("api_key", fetcher.ApiKey)
	params.Set("format", "json")
	return fetcher.consumer.ApiTraverse(fmt.Sprintf("api/game/%s", guid), params.Encode())
}

func (fetcher *Fetcher) Fetch(guid string) (GameResponse, error) {

	rawJSONs, err := fetcher.RawFetch(guid)
	if err != nil {
		return GameResponse{}, fmt.Errorf("Trouble fetching data from api: %s", err.Error())
	}

	// result is guaranteed singleton or empty
	var Game = GameResponse{}
	//fmt.Printf("THE STUFFFF %v", string(rawJSONs[0]))
	err = json.Unmarshal([]byte(rawJSONs[0]), &Game)
	if err != nil {
		return GameResponse{}, fmt.Errorf("Trouble decoding data returned from api: %s", err.Error())

	}

	return Game, nil
}

func (fetcher *Fetcher) FetchWDLCs(guid string) (GameResponseDLCs, error) {
	rawJSONs, err := fetcher.RawFetch(guid)
	if err != nil {
		return GameResponseDLCs{}, fmt.Errorf("Trouble fetching data from api: %s", err.Error())
	}

	var Game = GameResponseDLCs{}

	err = json.Unmarshal([]byte(rawJSONs[0]), &Game)
	if err != nil {
		return GameResponseDLCs{}, fmt.Errorf("Trouble decoding data returned from api: %s", err.Error())

	}

	betterDLCs, err := fetcher.addReleaseDate(Game.DLCs)

	sort.Sort(ByReleaseD(betterDLCs))

	Game.DLCs = betterDLCs
	return Game, err
}

func (fetcher *Fetcher) addReleaseDate(dets []DLCDetails) ([]DLCDetails, error) {

	params := url.Values{}
	//params.Set("limit", "100")
	params.Set("api_key", fetcher.ApiKey)
	params.Set("format", "json")
	// by default ReleaseDate will be 0
	var wg sync.WaitGroup
	wg.Add(len(dets))

	errChan := make(chan error)
	dataChan := make(chan DLCDetails, len(dets))
	allData := []DLCDetails{}

	for _, input := range dets {
		go func(input DLCDetails) {
			defer wg.Done()

			rawJSON, err := fetcher.consumer.ApiTraverse(fmt.Sprintf("/api/dlc/%d", input.ID), params.Encode())
			if err != nil {
				errChan <- fmt.Errorf("Error fetch DLC details: %s", err.Error())
			}

			// println(string(rawJSON[0]))

			elem := DLCDetails{}
			err = json.Unmarshal([]byte(rawJSON[0]), &elem)
			if err != nil {
				errChan <- fmt.Errorf("Error decoding JSON details: %s", err.Error())
			}
			dataChan <- elem

		}(input)
	}

	go func() {
		wg.Wait()
		close(errChan)
		close(dataChan)

	}()

	for e := range errChan {

		return []DLCDetails{}, e
	}
	for dItem := range dataChan {

		allData = append(allData, dItem)
	}

	return allData, nil

}
