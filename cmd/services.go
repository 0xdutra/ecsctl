/*
Copyright © 2022 Gabriel M. Dutra <0xdutra@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

type servicesOpts struct {
	clusterName              string
	serviceName              string
	serviceTaskDefinition    string
	serviceEnableForceDeploy bool
	serviceInputJson         string
	serviceMinCapacity       int64
	serviceMaxCapacity       int64
	serviceDesiredCapacity   int64
}

var so = servicesOpts{}

func servicesRun(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		os.Exit(1)
	}
}

// servicesCmd represents the services command
var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Commands to manage ECS services",
	Run:   servicesRun,
}

func init() {
	rootCmd.AddCommand(servicesCmd)
	servicesCmd.PersistentFlags().StringVarP(&so.clusterName, "cluster", "c", "", "The name of the ECS cluster")

	if err := servicesCmd.MarkPersistentFlagRequired("cluster"); err != nil {
		log.Fatal(err)
	}
}
