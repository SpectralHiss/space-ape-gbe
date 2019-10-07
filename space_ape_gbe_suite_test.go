package main_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSpaceApeGbe(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SpaceApeGbe Suite")
}

var cd, _ = os.Getwd()

func cmdWithArgs(ss ...string) *exec.Cmd {
	return exec.Command(fmt.Sprintf("%s/%s", cd, "space-ape-gbe"), ss...)
}

var _ = Describe("GiantBomb cli", func() {

	var serverURL = "http://localhost:8080"

	var cd, _ = os.Getwd()
	var fakeapicmd *exec.Cmd
	var serverOut io.ReadCloser

	var args []string

	var command *exec.Cmd
	var err error
	var output []byte

	BeforeSuite(func() {
		fakeapicmd = exec.Command(fmt.Sprintf("%s/fake-api/fake-api", cd))
		var err error
		Expect(err).To(BeNil())

		serverOut, err = fakeapicmd.StdoutPipe()

		Expect(err).To(BeNil())

		err = fakeapicmd.Start()
		Expect(err).To(BeNil())

		Eventually(func() error {
			_, err := http.Get(serverURL)
			return err
		}).Should(BeNil())

		cmd := exec.Command("go", "build", cd)
		err = cmd.Run()
		Expect(err).To(BeNil())

	})

	AfterSuite(func() {

		// time.Sleep(400000000000)
		err := fakeapicmd.Process.Kill()
		if err != nil {
			panic(err)
		}

		output, err := ioutil.ReadAll(serverOut)
		Expect(err).To(BeNil())
		fmt.Printf("FULL OUTPUT: %s", string(output))

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
				Expect(err).To(BeNil())
				Expect(string(output)).To(ContainSubstring("Searches GiantBomb's api for a game"))
			})
		})

		Context("When the search api is live and a valid simple search query is issued", func() {

			BeforeEach(func() {
				args = []string{"search", "half-life"}
			})

			FIt("Returns the right search results for that query term", func() {
				Expect(command.ProcessState.ExitCode()).To(Equal(0))
				//println(output)
				content, err := ioutil.ReadFile("./test_output/full.json")
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
				It("returns the game details without  the DLC", func() {

				})
			})

			When("the DLC is option is supplied", func() {
				It("returns the details including a summary of DLCs", func() {})
			})
		})

	})
})
