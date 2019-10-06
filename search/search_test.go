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
			searcher = search.NewSearcher("http://localhost:8080", "ce4949c5a501cdc3b0cdfbca070fd53787ba59a1")
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
		It("returns no results", func() {})
	})

	When("search is called but there is some error in the apis", func() {
		It("returns an informative error", func() {})
	})
})
