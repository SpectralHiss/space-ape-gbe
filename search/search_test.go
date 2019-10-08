package search_test

import (
	"net/http"
	"os"

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

	BeforeSuite(func() {

		Eventually(func() error {
			_, err := http.Get(os.Getenv("API_URL"))
			return err
		}, 5).Should(BeNil())

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
