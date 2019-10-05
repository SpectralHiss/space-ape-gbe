package apiconsumer_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
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

var _ = Describe("ApiConsumer", func() {
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

	var consumer *apiconsumer.GbeConsumer
	var outRes []json.RawMessage
	var outErr error
	var params url.Values

	BeforeEach(func() {
		consumer = apiconsumer.NewGBE(serverURL)
		params = url.Values{}
		params.Add("api_key", "ce4949c5a501cdc3b0cdfbca070fd53787ba59a1")
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
