# space-ape-gbe
Spaceape Technical Test: Giant Bomb Edition

We have packaged tests and the binary with docker-compose
run `docker-compose run tests` to run tests


set env `API_KEY` to your api key
`API_URL` should also be set to "https://giantbomb.com"
and run `docker-compose run binary` to run the binary

NOTE:
We discovered this api has pagination and only allows a maximum limit value of 100
Despite giantbomb's api has an invalid api key error the api just returns 401 with an empty body.
empty api_key returns INVALID API KEY response whereas an incorrect key returns empty body

TODO:
use query filtering to reduce payload

sort out null -> not yet released
flatten out Name -> string
spawn and waitgroup synchronization for each page request and dlc query

// currently this will dump the json as is, improvement could be to use "text/template" to beautify and flatten name or give option

