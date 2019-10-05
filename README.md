# space-ape-gbe
Spaceape Technical Test: Giant Bomb Edition

We discovered this api has pagination and only allows a maximum limit value of 100



NOTE:
Despite giantbomb's api has an invalid api key error the api just returns 401 with an empty body.
TODO:
add config for api_key (currently hardcoded)


// currently this will dump the json as is, improvement could be to use text/template


We decide on using json as output for composability

jq pattern used to extract correct structure for search:

cat 1.json 2.json 3.json | jq  '[.results[] | {name : .name, deck: .deck, release_date:(.expected_release_year // "not yet released") , platforms: [.platforms[]? | .name ]}] ' > output