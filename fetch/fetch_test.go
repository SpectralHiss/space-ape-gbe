package fetch_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	//. "UNKNOWN_PACKAGE_PATH"
)

var _ = Describe("Fetch", func() {
	When("fetch is called for a game", func(){
		When("And DLC is not supplied", func(){
			It("returns more information about that game, without dlc")
		})
		
		When("And DlC is supplied", func (){
			When("The game has DLCs", func (){
				It("returns a description with dlcs, sorted by release date (ascending)", func(){

				})
			})

			When("The game has no DLCs", func(){
				It("returns a description with an empty DLC entry")
			})
		})
	})

	When("search is called but the api is down", func() {
		It("returns an informative error", func() {})
	})
})
