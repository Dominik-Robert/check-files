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
	"time"

	"github.com/spf13/cobra"
)

// fileAgeCmd represents the fileAge command
var fileAgeCmd = &cobra.Command{
	Use:   "fileAge",
	Short: "Check the age of all files in a directory",
	Long: `The fileAge command checks the age of all files in a directory. For that you can specify which unit (DAYS, HOURS, MINUTES) 
you want to have and the warning and critical treshold.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("directory")
		entries, _ := os.ReadDir(dir)

		fileMap := make(map[string][]saveFile)
		warning, _ := cmd.Flags().GetInt64("warning")
		critical, _ := cmd.Flags().GetInt64("critical")
		unit, _ := cmd.Flags().GetString("unit")
		ignoreFiles, _ := cmd.Flags().GetStringArray("ignore")

		// outputFormat, _ := cmd.Flags().GetString("outputFormat")
		descriptionFormat, _ := cmd.Flags().GetString("descriptionFormat")

		exitString := "UNKNOWN - Something went wrong"
		exitCode := 3

		CacheRegexPatterns(ignoreFiles)

		var unitValue time.Duration

		switch unit {
		case "DAY", "DAYS":
			unitValue = time.Hour * 24
		case "HOUR", "HOURS":
			unitValue = time.Hour
		case "MINUTE", "MINUTES":
			unitValue = time.Minute
		}

		for _, entry := range entries {
			if !entry.IsDir() && !RegMatchArr(entry.Name()) {
				fileInfo, _ := entry.Info()
				modTime := fileInfo.ModTime()
				isCrit := false
				if modTime.Unix() < time.Now().Add(-1*unitValue*time.Duration(critical)).Unix() {
					isCrit = true
					fileMap["CRITICAL"] = append(fileMap["CRITICAL"], saveFile{
						Name: entry.Name(),
						Date: modTime,
					})
				}
				if modTime.Unix() < time.Now().Add(-1*unitValue*time.Duration(warning)).Unix() && !isCrit {
					fileMap["WARNING"] = append(fileMap["WARNING"], saveFile{
						Name: entry.Name(),
						Date: modTime,
					})
				}
			}
		}

		if len(fileMap["CRITICAL"]) > 0 {
			exitString = fmt.Sprint("CRITICAL - Some files are older than ", critical, " ", unit)
			exitCode = 2
		} else if len(fileMap["WARNING"]) > 0 {
			exitString = fmt.Sprint("WARNING - Some files are older than ", warning, " ", unit)
			exitCode = 1
		} else {
			exitString = fmt.Sprint("OK - All files are younger than ", warning, " ", unit)
			exitCode = 0

			fmt.Println(exitString)
			os.Exit(0)
		}

		if descriptionFormat == "MARKDOWN" {
			exitString += `
			| Filename | Fileage |
			| --- | --- |
			| CRITICAL | |
		`
			for _, value := range fileMap["CRITICAL"] {
				exitString += "|" + value.Name + "|" + value.Date.Format(time.ANSIC) + "|" + "\n"
			}

			exitString += `| WARNING | | `

			for _, value := range fileMap["WARNING"] {
				exitString += "|" + value.Name + "|" + value.Date.Format(time.ANSIC) + "|" + "\n"
			}
		} else {
			exitString += `
			<table>
				<tr>
					<th>Filename</th>
					<th>Fileage</th>
				</tr>
				<tr>
					<td colspan="2">CRITICAL</td>
				</tr>`

			for _, value := range fileMap["CRITICAL"] {
				exitString += "<tr><td>" + value.Name + "</td><td>" + value.Date.Format(time.ANSIC) + "</td></tr>"
			}

			exitString += `
			<tr>
				<td colspan="2">WARNING</td>
			</tr>
		`

			for _, value := range fileMap["WARNING"] {
				exitString += "<tr><td>" + value.Name + "</td><td>" + value.Date.Format(time.ANSIC) + "</td></tr>"
			}

			exitString += `
			<table>`
		}

		fmt.Println(exitString)
		os.Exit(exitCode)
	},
}

func init() {
	RootCmd.AddCommand(fileAgeCmd)
	fileAgeCmd.Flags().StringP("unit", "u", "DAY", "Set the unit for the warning value. Available units: (DAY, HOUR, MINUTE)")
	fileAgeCmd.Flags().String("outputFormat", "NAGIOS", "The format for the output. Available formats: NAGIOS and JSON")
	fileAgeCmd.Flags().String("descriptionFormat", "HTML", "The format for the output of the description. Available formats: HTML and MARKDWON")
}
