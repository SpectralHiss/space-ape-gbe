package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	 "github.com/gorilla/mux"
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

const API_K = "ce4949c5a501cdc3b0cdfbca070fd53787ba59a1"


var dir string

func StartServer() {
	var err error
	dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// server needs to have map of response to inputs
	// response types are: pages api response, simple response, errors
	// hard codes listeners  addResponse("term", GoodPagedResponse)



	// bad gameid

	rtr := mux.NewRouter()

	rtr.HandleFunc("/api/dlc/{guid}", func(w http.ResponseWriter, r *http.Request){
		logQuery(r)
		params := r.URL.Query()

		val := params.Get("api_key")
		if !checkAPIK(w, val) {
			return
		}
		massEDLCs := []string{"3100","3101", "3102", "3184", "3185", "3186", "3288", "3289", "3290", "3431", "3432",
		 "3433", "3561",  "3562", "3563", "3706",  "3742",  "3866", "3867", "3868", "3869" ,"3870", "3871",
		  "3919", "3920", "3921", "3922" ,"3929"}
		  vars := mux.Vars(r)
		  guid, ok := vars["guid"]
		  // mass effect 29935
		  if !ok {
			  notFound := fmt.Sprintf("%s/fixtures/fetch/no_guid_not_found.html")
			  writeResponseHTML(w, notFound , 404)
			  return
		  }

		  
		  if contains(massEDLCs, guid) {
			writeResponse(w, fmt.Sprintf("%s/fixtures/fetch/mass-effect-dlcs/%s.json", dir, guid), 200)
			return
		  } else {
			  println(dir)
			writeResponse(w, fmt.Sprintf("%s/fixtures/fetch/bad.json",dir),200)
			return
		  }

		})

	rtr.HandleFunc("/api/game/{guid}", func(w http.ResponseWriter, r *http.Request) {

		logQuery(r)
		params := r.URL.Query()

		val := params.Get("api_key")
		if !checkAPIK(w, val) {
			return
		}
		// api can have max limit of 100, so let's use that
		if params.Get("format") != "json" || params.Get("limit") != "100" {
			w.WriteHeader(400)
			w.Write([]byte("NOT FAKED"))
			return
		}
		vars := mux.Vars(r)
		guid, ok := vars["guid"]
		// mass effect 29935
		if !ok {
			notFound := fmt.Sprintf("%s/fixtures/fetch/no_guid_not_found.html",dir)
			writeResponseHTML(w, notFound , 404)
			return
		}
		
		if guid == "29935" {
			println(dir)
			writeResponse(w, fmt.Sprintf("%s/fixtures/fetch/mass-effect-3.json",dir), 200)
			return
		} 
		
		if guid == "999999999" {
			writeResponse(w, fmt.Sprintf("%s/fixtures/fetch/bad.json",dir),200)
			return
		// 10000000
		}
	})

	rtr.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		logQuery(r)
		params := r.URL.Query()

		val := params.Get("api_key")
		checkAPIK(w, val)
		if !checkAPIK(w, val) {
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

			writeResponse(w, fmt.Sprintf("%s/fixtures/no_query.json", dir), 200)
			return
		}

		if query == "half-life" {

			var filename string
			page := params.Get("page")
			if page == "" {
				filename = "first.json"
			} else {
				filename = fmt.Sprintf("%s.json", page)
			}

			writeResponse(w, fmt.Sprintf("%s/fixtures/search/half-life/%s", dir, filename), 200)
			return
		}
		if query == "zzefseqg" {

			writeResponse(w, fmt.Sprintf("%s/fixtures/search/zzefseqg/bad.json", dir), 200)
			return
		}

	})

	http.Handle("/", rtr)

	http.ListenAndServe(":8080", nil)
}

func main() {
	StartServer()
}



func writeResponseHTML(w http.ResponseWriter, fixturePath string, statusCode int) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(statusCode)
	content, err := ioutil.ReadFile(fixturePath)
	if err != nil {
		errPrint(err)
	}
	w.Write(content)
}

func writeResponse(w http.ResponseWriter, fixturePath string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	content, err := ioutil.ReadFile(fixturePath)
	if err != nil {
		errPrint(err)
	}
	w.Write(content)
}

func checkAPIK(w http.ResponseWriter, val string) bool {
	if val == "" || val != API_K {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(401)
		var content = []byte("")
		if val == "" {
			var err error
			content, err = ioutil.ReadFile(fmt.Sprintf("%s/fixtures/no_api_key.json", dir))
			if err != nil {
				errPrint(err)
			}
			w.Write(content)
		}

		return false
	}
	return true
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}