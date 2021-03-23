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

// fileSizeCmd represents the fileSize command
var fileSizeCmd = &cobra.Command{
	Use:   "fileSize",
	Short: "Check the size of a directory",
	Long:  `The fileSize command checks the size of all files or with the --single command you can check if any file comes over your threshold.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("directory")
		entries, _ := os.ReadDir(dir)

		warning, _ := cmd.Flags().GetInt64("warning")
		critical, _ := cmd.Flags().GetInt64("critical")
		showPerf, _ := cmd.Flags().GetBool("showPerf")
		unit, _ := cmd.Flags().GetString("unit")

		count := int64(0)
		modifier := int64(1)
		exitString := "UNKNOWN - Something went wrong"
		exitCode := 3

		for _, entry := range entries {
			if !entry.IsDir() {
				fileInfo, _ := entry.Info()
				count += fileInfo.Size()
			}
		}

		switch unit {
		case "B":
		case "KB":
			modifier *= 1000
		case "MB":
			modifier *= 1000 * 1000
		case "GB":
			modifier *= 1000 * 1000 * 1000
		case "TB":
			modifier *= 1000 * 1000 * 1000 * 1000
		}

		countValue := count / modifier

		if countValue >= critical {
			exitString = fmt.Sprint("CRITICAL - The files has a size of ", countValue, " ", unit)
			exitCode = 2
		} else if countValue >= warning {
			exitString = fmt.Sprint("WARNING - The files has a size of ", countValue, " ", unit)
			exitCode = 1
		} else {
			exitString = fmt.Sprint("OK - The files has a size of ", countValue, " ", unit)
			exitCode = 0
		}

		if showPerf {
			exitString += fmt.Sprintf(" | fileSize=%d;%d;%d ", countValue, warning, critical)
		}

		fmt.Println(exitString)
		os.Exit(exitCode)
	},
}

func init() {
	RootCmd.AddCommand(fileSizeCmd)

	fileSizeCmd.Flags().StringP("unit", "u", "MB", "Set the unit for the warning value. Available units: (B, KB, MB, GB, TB)")
	fileSizeCmd.Flags().BoolP("showPerf", "p", false, "Turns on performance data")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fileSizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fileSizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
