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
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
)

// editTaskDefinitionCmd represents the editTaskDefinition command
var editTaskDefinitionCmd = &cobra.Command{
	Use:   "edit-task-definition",
	Short: "Edit a task definition and create a new revision",
	Run:   editTaskdefinitionRun,
}

var (
	editor       string
	taskName     string
	taskRevision int64
)

func init() {
	taskDefinitionCmd.AddCommand(editTaskDefinitionCmd)
	taskDefinitionCmd.PersistentFlags().StringVarP(&editor, "editor", "", "vim", "The name of the text editor to use to edit task definition")
	taskDefinitionCmd.PersistentFlags().StringVarP(&taskName, "name", "", "", "The name of the task definition")
	taskDefinitionCmd.PersistentFlags().Int64VarP(&taskRevision, "revision", "", 1, "The revision of the task definition")
}

func editTaskdefinitionRun(cmd *cobra.Command, _ []string) {
	sess := provider.NewSession()
	svc := ecs.New(sess)

	taskDefinitionWithRevision := fmt.Sprintf("%s:%d", taskName, taskRevision)

	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinitionWithRevision),
	}

	taskDefinition, err := svc.DescribeTaskDefinition(input)
	if err != nil {
		log.Panic(err)
	}

	newTaskDefinitionJSON, _ := json.MarshalIndent(taskDefinition.TaskDefinition, "", "  ")
	tmpTaskDefinitionFile := fmt.Sprintf("tmp_%s_%d.json", taskName, taskRevision)

	newTaskDefinition, err := os.Create(tmpTaskDefinitionFile)
	if err != nil {
		log.Panic(err)
	}

	defer newTaskDefinition.Close()

	_, err = newTaskDefinition.WriteString(string(newTaskDefinitionJSON))
	if err != nil {
		log.Panic(err)
	}

	editorCommand := exec.Command(editor, tmpTaskDefinitionFile)
	editorCommand.Stdin = os.Stdin
	editorCommand.Stdout = os.Stdout
	editorCommand.Run()

	readTmpTaskDefinitionFile, err := os.Open(tmpTaskDefinitionFile)
	if err != nil {
		log.Panic(err)
	}

	defer readTmpTaskDefinitionFile.Close()

	rawTaskInput, err := ioutil.ReadAll(readTmpTaskDefinitionFile)
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

	fmt.Println(result.TaskDefinition)

	err = os.Remove(tmpTaskDefinitionFile)
	if err != nil {
		log.Panic(err)
	}
}
