package search_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/SpectralHiss/space-ape-gbe/search"
	//. "github.com/SpectralHiss/space-ape-gbe/search"
)

func TestSearch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Search Suite")
}

// func MustEnv(envString string) string {
// 	val, err := os.Env(envString)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return val
// }

var _ = Describe("Search", func() {

	var serverURL = "http://localhost:8080"

	var cd, _ = os.Getwd()
	var cmd *exec.Cmd
	var serverOut io.ReadCloser

	BeforeSuite(func() {
		cmd = exec.Command(fmt.Sprintf("%s/../fake-api/fake-api", cd))
		var err error
		Expect(err).To(BeNil())

		serverOut, err = cmd.StdoutPipe()

		Expect(err).To(BeNil())

		err = cmd.Start()
		Expect(err).To(BeNil())

		Eventually(func() error {
			_, err := http.Get(serverURL)
			return err
		}).Should(BeNil())

	})

	AfterSuite(func() {

		// time.Sleep(400000000000)
		err := cmd.Process.Kill()
		if err != nil {
			panic(err)
		}

		output, err := ioutil.ReadAll(serverOut)
		Expect(err).To(BeNil())
		fmt.Printf("FULL OUTPUT: %s", string(output))

	})

	var searcher *search.Searcher

	When("search is called with a param that yields results", func() {
		BeforeEach(func() {
			url := os.Getenv("API_URL")
			key := os.Getenv("API_KEY")

			searcher = search.NewSearcher(url, key)
		})
		It("returns a condensed list of results", func() {

			res, err := searcher.SearchTitles("half-life")
			Expect(err).To(BeNil())
			Expect(len(res) == 282)
			// first last to triple check

			Expect(res[0]).To(Equal(search.SearchResponse{
				Name:        "Half-Life",
				Deck:        "Take on the role of Gordon Freeman as he escapes the disastrous aftermath of an experiment gone wrong in the Black Mesa Research Facility.",
				ReleaseDate: 1998,
				ID:          2980,
				Platforms: []struct{ Name string }{
					{
						Name: "Mac",
					},
					{
						Name: "PlayStation 2",
					},
					{
						Name: "PC",
					},
					{
						Name: "Linux",
					},
				},
			}))
		})

	})
	When("search is called with a param that has no results", func() {
		It("returns no results", func() {
			res, err := searcher.SearchTitles("zzefseqg")
			Expect(err).To(BeNil())
			Expect(len(res)).To(Equal(0))
		})
	})

	//TODO
	When("search is called but the api is down", func() {
		It("returns an informative error", func() {})
	})
})
