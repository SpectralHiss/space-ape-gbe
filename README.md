# space-ape-gbe
Spaceape Technical Test: Giant Bomb Edition


TODO:
add config for api_key (currently hardcoded)

We decide on using json as output for composability

jq pattern used to extract correct structure for search:

cat 1.json 2.json 3.json | jq  '[.results[] | {name : .name, deck: .deck, release_date:(.expected_release_year // "not yet released") , platforms: [.platforms[]? | .name ]}] ' > output