package apiconsumer_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/SpectralHiss/space-ape-gbe/apiconsumer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestApiConsumer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Consumer Suite")
}

func IsJSON(s json.RawMessage) bool {
	var js interface{}

	return json.Unmarshal([]byte(s), &js) == nil
}

var cd string

var _ = Describe("ApiConsumer", func() {
	var serverURL = os.Getenv("API_URL")

	BeforeSuite(func() {

		Eventually(func() error {
			_, err := http.Get(serverURL)
			return err
		}).Should(BeNil())

	})

	var consumer *apiconsumer.GbeConsumer
	var outRes []json.RawMessage
	var outErr error
	var params url.Values

	BeforeEach(func() {
		consumer = apiconsumer.NewGBE(serverURL)
		params = url.Values{}
		params.Add("api_key", os.Getenv("API_KEY"))
		params.Add("resources", "game")
		params.Add("format", "json")
		params.Add("limit", "100")
	})

	JustBeforeEach(func() {

		outRes, outErr = consumer.ApiTraverse("api/search", params.Encode())
	})

	When("there are no search results", func() {
		BeforeEach(func() {

			params.Add("query", "zzefseqg")
		})

		It("returns an empty response, no error", func() {
			Expect(outErr).To(BeNil())
			Expect(outRes).To(Equal([]json.RawMessage{}))
		})
	})

	When("A valid api query is passed that requires pagination (half-life)", func() {
		BeforeEach(func() {
			params.Add("query", "half-life")
		})
		// we test the prefix of each one of the results returned
		It("returns the concatenation of the results", func() {
			Expect(outErr).To(BeNil())
			Expect(len(outRes)).To(Equal(3))

			Expect(IsJSON(outRes[0]))

			Expect(IsJSON(outRes[1]))

			Expect(IsJSON(outRes[2]))

			// page 1 results output string
			Expect(outRes[0]).To(HavePrefix(`[{"aliases":"HL1\r\nHL\r\nHalf-Life: Source"`))

			// // page 2 results output string
			Expect(outRes[1]).To(HavePrefix(`[{"aliases":"Ryu ga Gotoku 6"`))
			// // page 3
			Expect(outRes[2]).To(HavePrefix(`[{"aliases":null,"api_detail_url":"https:\/\/www.giantbomb.com\/api\/game\/3030-47927`))

		})
	})

	When("there are api returned errors in a call, due to bad api_key", func() {
		BeforeEach(func() {
			params.Add("query", "half-life")
			params.Set("api_key", "")
		})

		It("returns an informative api error", func() {

			Expect(outErr).NotTo(BeNil())
			Expect(outErr).To(MatchError("Invalid API Key"))
		})
	})

	When("there are api returned errors in a call, due to no query", func() {
		BeforeEach(func() {
			params.Add("query", "")
		})

		It("returns an informative api error", func() {

			Expect(outErr).NotTo(BeNil())
			Expect(outErr).To(MatchError("Object Not Found"))
		})
	})

	// When("there the api_endpoit is unreachable", func() {
	// 	JustBeforeEach(func() {
	// 		err := cmd.Process.Kill()
	// 		Expect(err).To(BeNil())
	// 	})

	// 	It("returns an informative error", func() {
	// 		Expect(outErr).NotTo(BeNil())
	// 		Expect(outErr).To(MatchError("The API endpoint is currently unreachable."))
	// 	})
	// })
})
