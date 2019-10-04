package apiconsumer_test

import (
	"bytes"
	"fmt"
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

var _ = Describe("ApiConsumer", func() {

	var cd, _ = os.Getwd()
	var cmd *exec.Cmd
	BeforeSuite(func() {
		cmd = exec.Command(fmt.Sprintf("%s/../fake-api/fake-api", cd))
		err := cmd.Start()
		Expect(err).To(BeNil())
	})

	AfterSuite(func() {
		//err := cmd.Wait()
		//Expect(err).To(BeNil())
	})

	var consumer *apiconsumer.GbeConsumer
	var outRes []string
	var outErr error
	var params url.Values

	BeforeEach(func() {
		consumer = apiconsumer.NewGBE("http://localhost:8080")
		params = url.Values{}
		params.Add("api_key", "ce4949c5a501cdc3b0cdfbca070fd53787ba59a1")
		params.Add("resources", "game")
		params.Add("format", "json")
		params.Add("limit", "100")
	})

	JustBeforeEach(func() {

		req, err := http.NewRequest("GET", params.Encode(), &bytes.Buffer{})
		Expect(err).To(BeNil())
		outRes, outErr = consumer.ApiTraverse(req)
	})

	When("there are no search results", func() {
		BeforeEach(func() {

			params.Add("query", "zzefseqg")
		})

		FIt("returns an empty response, no error", func() {
			Expect(outErr).To(BeNil())
			Expect(outRes).To(Equal([]string{}))
		})
	})

	When("A valid api query is passed that requires pagination (half-life)", func() {
		// we test the prefix of each one of the results returned
		It("returns the concatenation of the results", func() {

		})
	})

	When("there are api returned errors in a call", func() {
		It("returns the api errors", func() {

		})
	})

	When("there the api_endpoit is unreachable", func() {
		It("returns a custom errors", func() {

		})
	})
})
