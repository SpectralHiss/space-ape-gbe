# space-ape-gbe
Spaceape Technical Test: Giant Bomb Edition

We have packaged tests with docker-compose
to run, need to first compile binary:
`docker build . -t space-ape-gbe_binary`
then
run `docker-compose run test` to run tests


set env `API_KEY` to your api key
`API_URL` should also be set to "https://giantbomb.com"

run the binary with : 
`docker run  -e API_URL=https://giantbomb.com -e API_KEY=YOUR_KEY IMAGE`

## Running confusion!!
sorry about the trouble I had devised a testing dockerfile and it was not checked in, I had to rush at the end
I have committed it now so you can fairly appraise it.. One last thing as well, make sure the 
--from step to matche whatever label you built the binary with
the binary image with, it needs the binary because some tests shell out to test as black box
oddly the From is not a syntax error and runs on my machine! 

## API GOTCHAS:
We discovered this api has pagination and only allows a maximum limit value of 100
Despite giantbomb's api has an invalid api key error the api just returns 401 with an empty body.
empty api_key returns INVALID API KEY response whereas an incorrect key returns empty body

## TODO:
- use query filtering to reduce payload
- parallelize page fetch also (only done dlc fetching in parallel)

- presentation:
  - sort out null -> not yet released
  - flatten out Name -> string
  - format flag ? yaml json text/template


