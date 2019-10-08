# space-ape-gbe
Spaceape Technical Test: Giant Bomb Edition

We have packaged tests and the binary with docker-compose
run `docker-compose run tests` to run tests


set env `API_KEY` to your api key
`API_URL` should also be set to "https://giantbomb.com"

run the binary with : 
`docker run  -e API_URL=https://giantbomb.com -e API_KEY=YOUR_KEY IMAGE`

# API GOTCHAS:
We discovered this api has pagination and only allows a maximum limit value of 100
Despite giantbomb's api has an invalid api key error the api just returns 401 with an empty body.
empty api_key returns INVALID API KEY response whereas an incorrect key returns empty body

# TODO:
- use query filtering to reduce payload
- parallelize page fetch also (only done dlc fetching in parallel)

- presentation:
  - sort out null -> not yet released
  - flatten out Name -> string
  - format flag ? yaml json text/template


