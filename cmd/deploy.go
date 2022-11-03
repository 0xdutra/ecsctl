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

	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Commands to deploy ECS infrastructure",
	Run:   deployRun,
}

var configName string
var resource string

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.PersistentFlags().StringVarP(&configName, "config", "c", "", "The name of the yaml config")
	deployCmd.PersistentFlags().StringVarP(&resource, "resource", "r", "", "The name of the resource to deploy")

	if err := deployCmd.MarkPersistentFlagRequired("config"); err != nil {
		log.Fatal(err)
	}
}

func deployRun(cmd *cobra.Command, args []string) {

	switch resource {
	case "task-definition":
		registerTaskDefinition(configName)
	case "service":
		createService(configName)
	case "load-balancer":
		createElb(configName)
	}
}
