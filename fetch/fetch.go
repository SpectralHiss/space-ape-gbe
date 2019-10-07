package fetch

type GameResponse struct {
	Name        string `json:"name"`
	Deck        string `json:"deck"`
	ReleaseDate int    `json:"expected_release_year"`
	Genres      string `json:"genres"`
	Platforms   []struct {
		Name string
	} `json:"platforms"`

	ID          int    `json:"id"`
	Description string `json:"description"`
	GameRating  string `json:"original_game_rating"`
	Characters  string `json:"characters"`
	Concepts    string `json:"concepts"`
	Publishers  string `json:"publishers`
	Developers  string `json:"developers"`
	Franchises  string `json:"franchises"`
}

type DLCDetails struct {
}

type GameResponseDLCs struct {
	GameResponse
	DLCs []DLCDetails
}
