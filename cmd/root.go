/*
Copyright © 2021 Dominik Robert dmnkrobert@gmail.com

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
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/spf13/cobra"
)

const VERSION = "0.9.0"

var (
	CachedRegexp []*regexp.Regexp
)

type saveFile struct {
	Name string
	Date time.Time
	Size int64
}

func CacheRegexPatterns(matcher []string) {
	for _, value := range matcher {
		cachedReg := regexp.MustCompile(value)
		CachedRegexp = append(CachedRegexp, cachedReg)
	}
}

// RegMatchArr returns true if the value matches an array full of regex expressions
func RegMatchArr(match string) bool {
	for _, value := range CachedRegexp {
		if value.MatchString(match) {
			return true
		}
	}
	return false
}

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "check-files",
	Short: "A nagios compliant check for different file operations",
	Long:  `With this check you can perform check operations with your files to prove if they are correct in age, size and count`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringP("directory", "d", "./", "Specify the directory which contains the files")
	RootCmd.PersistentFlags().StringArray("ignore", []string{}, "Ignore specific files. Can contain single files or a regex")
	RootCmd.PersistentFlags().Int64P("warning", "w", 10, "Specify the Warning value. All files which are older than 10 units")
	RootCmd.PersistentFlags().Int64P("critical", "c", 14, "Specify the Critical value. All files which are older than 10 units")

	CachedRegexp = make([]*regexp.Regexp, 0)
}
