package fetch_test

import (
	"net/http"
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"

	"os"

	"github.com/SpectralHiss/space-ape-gbe/fetch"

	//"path/filepath"
	//"log"
	"fmt"
	"io/ioutil"

	"encoding/json"
	"testing"
)

func TestFetch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fetch Suite")
}

var root string

var _ = Describe("Fetch", func() {
	format.TruncatedDiff = true

	// var cd, _ = os.Getwd()
	// var cmd *exec.Cmd
	// var serverOut io.ReadCloser

	BeforeSuite(func() {
		// 	cmd = exec.Command(fmt.Sprintf("%s/../fake-api/fake-api", cd))
		// 	var err error
		// 	Expect(err).To(BeNil())

		// 	serverOut, err = cmd.StdoutPipe()

		// 	Expect(err).To(BeNil())

		// 	err = cmd.Start()
		// 	Expect(err).To(BeNil())

		Eventually(func() error {
			_, err := http.Get(os.Getenv("API_URL"))
			return err
		}, 5).Should(BeNil())

	})

	// AfterSuite(func() {

	// 	// time.Sleep(400000000000)
	// 	err := cmd.Process.Kill()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	output, err := ioutil.ReadAll(serverOut)
	// 	Expect(err).To(BeNil())
	// 	fmt.Printf("FULL OUTPUT: %s", string(output))

	// })

	var fetcher *fetch.Fetcher
	var out fetch.GameResponse

	var err error

	root = ".."
	//var outDLC []fetch.GameResponseDLC
	var url string
	var key string
	BeforeEach(func() {
		url = os.Getenv("API_URL")
		key = os.Getenv("API_KEY")

		fetcher = fetch.NewFetcher(url, key)
	})

	When("fetch is called for a real game", func() {
		var out2 fetch.GameResponseDLCs

		When("And DLC is not supplied", func() {
			BeforeEach(func() {
				out, err = fetcher.Fetch("29935")
			})

			It("returns more information about that game, without dlc", func() {
				Expect(err).To(BeNil())
				content, err := ioutil.ReadFile(fmt.Sprintf("%s/test_output/full-fetch-no-dlc.json", root))

				Expect(err).To(BeNil())
				parsed := fetch.GameResponse{}
				err = json.Unmarshal(content, &parsed)
				Expect(err).To(BeNil())
				Expect(out).To(Equal(parsed))
			})
		})

		When("And DlC is supplied", func() {
			BeforeEach(func() {
				out2, err = fetcher.FetchWDLCs("29935")
			})
			When("The game has DLCs", func() {
				It("returns a description with dlcs, sorted by release date (ascending)", func() {
					Expect(err).To(BeNil())
					Expect(len(out2.DLCs)).NotTo(BeZero())
					Expect(sort.IsSorted(fetch.ByReleaseD(out2.DLCs))).To(BeTrue())
					//Expect(out2.DLCs[0].ReleaseDate).NotTo(BeEmpty())
				})
			})

			// TODO check
			When("The game has no DLCs", func() {
				It("returns a description with an empty DLC entry", func() {})
			})
		})

	})

	When("fetch is called with a non-existent guid", func() {
		BeforeEach(func() {
			out, err = fetcher.Fetch("999999999")
		})

		It("returns an informative error", func() {
			Expect(err).NotTo(BeNil())
			Expect(err).To(MatchError("Trouble fetching data from api: Object Not Found"))
		})
	})

	When("search is called but the api is down", func() {
		// BeforeEach(func() {
		// 	//url = "http://example.com"
		// 	//key = "garbage"
		// })

		// //TODO test setup url rather than api_k
		// It("returns an informative error", func() {
		// 	Expect(err).NotTo(BeNil())

		// 	Expect(err).To(MatchError("Trouble fetching data from api: The API endpoint is currently unreachable"))
		// })
	})
})
