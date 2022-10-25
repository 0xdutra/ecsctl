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
	"ecsctl/pkg/provider"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
)

// registerTaskdefinitionCmd represents the registerTaskdefinition command
var registerTaskdefinitionCmd = &cobra.Command{
	Use:   "register",
	Short: "Register task definition using JSON file",
	Run:   registerTaskdefinitionRun,
}

func init() {
	taskDefinitionCmd.AddCommand(registerTaskdefinitionCmd)
	registerTaskdefinitionCmd.PersistentFlags().StringVarP(&tdfo.taskJsonInput, "input-json", "", "", "Input task definition with JSON format")
}

func registerTaskdefinitionRun(cmd *cobra.Command, _ []string) {
	sess := provider.NewSession()
	svc := ecs.New(sess)

	taskFile, err := os.Open(tdfo.taskJsonInput)
	if err != nil {
		log.Panic(err)
	}

	defer taskFile.Close()

	rawTaskInput, err := ioutil.ReadAll(taskFile)
	if err != nil {
		log.Panic(err)
	}

	var taskInput ecs.RegisterTaskDefinitionInput

	err = json.Unmarshal([]byte(rawTaskInput), &taskInput)
	if err != nil {
		log.Panic(err)
	}

	result, err := svc.RegisterTaskDefinition(&taskInput)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(result)
}
