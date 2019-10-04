package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//type fake_gbe_server struct {
//}

// minimal fake illustrating the subset of the api we care about and our assumptions
func StartServer() {

	// server needs to have map of response to inputs
	// response types are: pages api response, simple response, errors
	// hard codes listeners  addResponse("term", GoodPagedResponse)

	const API_K = "ce4949c5a501cdc3b0cdfbca070fd53787ba59a1"

	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		cwd, _ := os.Getwd()

		val := params.Get("api_key")
		if val == "" || val != API_K {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(401)

			content, _ := ioutil.ReadFile(fmt.Sprintf("%s/fixtures/no_api_key.json", cwd))
			w.Write(content)
		}

		// api can have max limit of 100, so let's use that
		if params.Get("resources") != "game" || params.Get("format") != "json" || params.Get("limit") != "100" {
			w.WriteHeader(400)
			w.Write([]byte("NOT FAKED"))
		}
		query := params.Get("query")
		if query == "" {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(200)
			content, _ := ioutil.ReadFile(fmt.Sprintf("%s/fixtures/no_query.json", cwd))

			w.Write(content)

		}

		if query == "half-life" {

			w.Header().Set("Content-Type", "application/json; charset=utf-8")

			w.WriteHeader(200)

			var filename string
			page := params.Get("page")
			if page == "" {
				filename = "first.json"
			} else {
				filename = fmt.Sprintf("%s.json", page)
			}

			content, _ := ioutil.ReadFile(fmt.Sprintf("%s/fixtures/half-life/%s", cwd, filename))

			w.Write(content)
		}
		if query == "zzefseqg" {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(200)

			content, _ := ioutil.ReadFile(fmt.Sprintf("%s/fixtures/zzefseqg/bad.json", cwd))
			w.Write(content)
		}

	})

	http.ListenAndServeTLS(":8080", nil)
}

func main() {
	StartServer()
}
