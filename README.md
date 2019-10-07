# space-ape-gbe
Spaceape Technical Test: Giant Bomb Edition

We discovered this api has pagination and only allows a maximum limit value of 100



NOTE:
Despite giantbomb's api has an invalid api key error the api just returns 401 with an empty body.
empty api_key returns INVALID API KEY response whereas an incorrect key returns empty body

TODO:
remove existing api_key refs and change to env vars 
use query filtering to reduce payload
sort out null -> not yet released

// currently this will dump the json as is, improvement could be to use "text/template" to beautify or give option

jq pattern used to extract correct structure for search:

cat 1.json 2.json 3.json | jq  '[.results[] | {name : .name, deck: .deck, release_date:(.expected_release_year // "not yet released") , platforms: [.platforms[]? | .name ], id: .id}] ' > output.json