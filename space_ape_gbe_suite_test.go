package main_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"

	"github.com/SpectralHiss/space-ape-gbe/fetch"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSpaceApeGbe(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SpaceApeGbe Suite")
}

var cd string

func cmdWithArgs(ss ...string) *exec.Cmd {
	return exec.Command(fmt.Sprintf("%s/%s", cd, "space-ape-gbe"), ss...)
}

var _ = Describe("GiantBomb cli", func() {

	var serverURL = os.Getenv("API_URL")

	var err error
	cd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var args []string

	var command *exec.Cmd

	var output []byte

	BeforeSuite(func() {

		Eventually(func() error {
			_, err := http.Get(serverURL)
			return err
		}).Should(BeNil())

	})

	JustBeforeEach(func() {
		command = cmdWithArgs(args...)

		output, err = command.CombinedOutput()
	})

	Describe("Searching for game titles", func() {

		Context("when the search command is invoked with no args", func() {
			BeforeEach(func() {
				args = []string{"search"}
			})
			It("returns a command description message", func() {
				Expect(err).NotTo(BeNil())
				Expect(string(output)).To(ContainSubstring("Error: requires at least 1 arg(s), only received 0"))
			})
		})

		Context("When the search api is live and a valid simple search query is issued", func() {

			BeforeEach(func() {
				args = []string{"search", "half-life"}
			})

			It("Returns the right search results for that query term", func() {
				Expect(command.ProcessState.ExitCode()).To(Equal(0))
				//println(output)
				content, err := ioutil.ReadFile("./test_output/full-search.json")
				Expect(err).To(BeNil())
				//fmt.Printf(string(content))
				Expect(output).To(MatchJSON(string(content)))
			})
		})

	})

	Describe("Getting a specific game by ID", func() {
		When("A game is fetched by ID", func() {
			BeforeEach(func() {
				args = []string{"fetch", "29935"}
			})

			When("the DLC option is not supplied", func() {
				It("returns the game details without  the DLCs", func() {
					content, err := ioutil.ReadFile("./test_output/full-fetch-no-dlc.json")
					Expect(err).To(BeNil())
					//fmt.Printf(string(content))
					Expect(output).To(MatchJSON(string(content)))

				})
			})

			When("the DLC is option is supplied", func() {

				BeforeEach(func() {
					args = []string{"fetch", "29935", "--dlcs"}
				})

				It("returns the details including the DLCs", func() {
					data := fetch.GameResponseDLCs{}

					err := json.Unmarshal(output, &data)
					Expect(err).To(BeNil())

					Expect(len(data.DLCs)).To(Equal(30))
					//Expect(data.DLCs[0].ReleaseDate).NotTo(Equal(""))
				})
			})
		})

		// TODO case where game has  no DLC

	})
})
