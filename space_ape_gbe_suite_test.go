package main_test

import (
	"fmt"
	"io/ioutil"
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

var _ = Describe("searching GiantBomb for a game", func() {

	BeforeSuite(func() {
		cmd := exec.Command("go", "build", cd)
		err := cmd.Run()
		Expect(err).To(BeNil())
	})

	AfterSuite(func() {
		//time.Sleep(100000000000)

	})

	var args []string

	var command *exec.Cmd
	var err error
	var output []byte

	JustBeforeEach(func() {
		command = cmdWithArgs(args...)

		output, err = command.CombinedOutput()
	})

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

		It("Returns the right search results for that query term", func() {
			Expect(command.ProcessState.ExitCode()).To(Equal(0))
			println(output)
			content, err := ioutil.ReadFile("./test_output/full.json")
			Expect(err).To(BeNil())
			Expect(output).To(MatchJSON(string(content)))
		})
	})
})
