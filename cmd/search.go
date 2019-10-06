/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"os"

	"github.com/SpectralHiss/space-ape-gbe/search"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func jsonPrint(res []search.SearchResponse) {
	bytes, err := json.Marshal(res)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem formatting result output")
	}

	fmt.Println(string(bytes))
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Searches GiantBomb's api for a game",
	Long:  `Search searches titles blablab`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := viper.GetString("API_URL")
		key := viper.GetString("API_KEY")

		searcher := search.NewSearcher(url, key)
		out, err := searcher.SearchTitles(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "We having problem searching games: %s", err.Error())
		}

		jsonPrint(out)
		//fmt.Println("search called")

	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
