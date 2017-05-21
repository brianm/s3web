// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/brianm/s3web/fix"
	"github.com/spf13/cobra"
	"log"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(fixCmd)
	fixCmd.Flags().BoolP("detect-unknown-mimetypes", "d", false, "Sniff for mimetypes if unknown suffix")
	viper.BindPFlags(fixCmd.Flags())
}

// fixCmd represents the fix command
var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Ensure contents of a bucket are reasonably correct",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		bucket := viper.GetString("bucket")
		if bucket == "" {
			bucket = args[0]
		}

		f := fix.Fix{
			Bucket:   bucket,
			Verbose:  viper.GetBool("verbose"),
			Simulate: viper.GetBool("simulate"),
			DetectMimeTypesForUnknownNames: viper.GetBool("detect-unknown-mimetypes"),
		}

		if err := f.Fix(); err != nil {
			log.Fatalf("unable to fix %s: %+v", f.Bucket, err)
		}
	},
}
