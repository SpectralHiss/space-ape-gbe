package main_test

import (
	"fmt"
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
	var args []string

	var command *exec.Cmd

	JustBeforeEach(func() {
		cmd := exec.Command("go", "build", cd)
		err := cmd.Run()
		Expect(err).To(BeNil())

		command = cmdWithArgs(args...)
		err = command.Run()
		Expect(err).To(BeNil())
		command.Run()

	})

	Context("when the search command is invoked with no args", func() {
		BeforeEach(func() {
			args = []string{"search"}
		})
		It("returns a usage message", func() {

		})
	})

	Context("When the search api is live and a valid simple search query is issued", func() {

		BeforeEach(func() {
			args = []string{"search", "pubg"}
		})

		It("Returns the right search results for that query term", func() {
			Expect(command.ProcessState.ExitCode()).To(Equal(0))
			//Expect(command.CombinedOutput()).To(MatchJSON())

		})
	})
})
