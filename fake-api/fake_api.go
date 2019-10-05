package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

//type fake_gbe_server struct {
//}

// minimal fake illustrating the subset of the api we care about and our assumptions

func errPrint(err error) {
	fmt.Fprintf(os.Stderr, err.Error())
}

func logQuery(r *http.Request) {
	log.Print(r.RequestURI)
}

func StartServer() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// server needs to have map of response to inputs
	// response types are: pages api response, simple response, errors
	// hard codes listeners  addResponse("term", GoodPagedResponse)

	const API_K = "ce4949c5a501cdc3b0cdfbca070fd53787ba59a1"

	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		logQuery(r)
		params := r.URL.Query()

		val := params.Get("api_key")
		if val == "" || val != API_K {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(401)

			// content, err := ioutil.ReadFile(fmt.Sprintf("%s/fixtures/no_api_key.json", dir))
			// if err != nil {
			// 	errPrint(err)
			// }
			//w.Write(content)
			w.Write([]byte(""))
			return
		}

		// api can have max limit of 100, so let's use that
		if params.Get("resources") != "game" || params.Get("format") != "json" || params.Get("limit") != "100" {
			w.WriteHeader(400)
			w.Write([]byte("NOT FAKED"))
			return
		}
		query := params.Get("query")
		if query == "" {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(200)
			content, err := ioutil.ReadFile(fmt.Sprintf("%s/fixtures/no_query.json", dir))
			if err != nil {
				errPrint(err)
			}
			w.Write(content)
			return
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

			content, err := ioutil.ReadFile(fmt.Sprintf("%s/fixtures/half-life/%s", dir, filename))
			if err != nil {
				errPrint(err)
			}
			w.Write(content)
			return
		}
		if query == "zzefseqg" {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(200)

			content, err := ioutil.ReadFile(fmt.Sprintf("%s/fixtures/zzefseqg/bad.json", dir))
			if err != nil {
				errPrint(err)
			}
			w.Write(content)
			return
		}

	})

	http.ListenAndServe(":8080", nil)
}

func main() {
	StartServer()
}
