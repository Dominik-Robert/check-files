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

	"github.com/spf13/cobra"
)

// fileCountCmd represents the fileCount command
var fileCountCmd = &cobra.Command{
	Use:   "fileCount",
	Short: "Count your files in a directory",
	Long:  `With fileCount you can count all files in a directory and compare it with your threshold.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("directory")
		entries, _ := os.ReadDir(dir)

		warning, _ := cmd.Flags().GetInt64("warning")
		critical, _ := cmd.Flags().GetInt64("critical")
		showPerf, _ := cmd.Flags().GetBool("showPerf")

		count := int64(0)
		exitString := "UNKNOWN - Something went wrong"
		exitCode := 3

		for _, entry := range entries {
			if !entry.IsDir() {
				count++
			}
		}

		if count >= critical {
			exitString = fmt.Sprint("CRITICAL - There are ", count, " files in the directory")
			exitCode = 2
		} else if count >= warning {
			exitString = fmt.Sprint("WARNING - There are ", count, " files in the directory")
			exitCode = 1
		} else {
			exitString = fmt.Sprint("OK - There are ", count, " files in the directory")
			exitCode = 0
		}

		if showPerf {
			exitString += fmt.Sprintf(" | fileCount=%d;%d;%d", count, warning, critical)
		}

		fmt.Println(exitString)
		os.Exit(exitCode)
	},
}

func init() {
	RootCmd.AddCommand(fileCountCmd)
	fileCountCmd.Flags().BoolP("showPerf", "p", false, "Turns on performance data")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fileCountCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fileCountCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
