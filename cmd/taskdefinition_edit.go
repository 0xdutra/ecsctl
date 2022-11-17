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
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/andreyvit/diff"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
)

// editTaskDefinitionCmd represents the editTaskDefinition command
var editTaskDefinitionCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a task definition and create a new revision",
	Run:   editTaskdefinitionRun,
}

func init() {
	taskDefinitionCmd.AddCommand(editTaskDefinitionCmd)
	editTaskDefinitionCmd.PersistentFlags().StringVarP(&tdfo.editor, "editor", "", "vim", "The name of the text editor to use to edit task definition")
	editTaskDefinitionCmd.PersistentFlags().StringVarP(&tdfo.taskName, "name", "", "", "The name of the task definition")
	editTaskDefinitionCmd.PersistentFlags().Int64VarP(&tdfo.taskRevision, "revision", "", 1, "The revision of the task definition")
}

func createDiff(currentTaskDefinition string, editedTaskDefinition string) string {
	return diff.LineDiff(currentTaskDefinition, editedTaskDefinition)
}

func getCurrentTaskDefinition(taskDefinitionName string, taskDefinitionRevision int64) *ecs.DescribeTaskDefinitionOutput {
	sess := provider.NewSession()
	svc := ecs.New(sess)

	taskDefinitionWithRevision := fmt.Sprintf("%s:%d", tdfo.taskName, tdfo.taskRevision)

	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinitionWithRevision),
	}

	currentTaskDefinition, err := svc.DescribeTaskDefinition(input)
	if err != nil {
		log.Panic(err)
	}

	return currentTaskDefinition
}

func waitUserAcceptChanges() bool {
	var resp string

	fmt.Println("\nDo you want to make these changes?")
	fmt.Println("ecsctl will create a new revision of the task definition.")
	fmt.Println("Only 'yes' will be accepted to approve.")

	fmt.Printf("Enter value: ")
	fmt.Scanf("%s", &resp)

	return resp == "yes"
}

func editTaskdefinitionRun(cmd *cobra.Command, _ []string) {
	sess := provider.NewSession()
	svc := ecs.New(sess)

	currentTaskDefinition := getCurrentTaskDefinition(tdfo.taskName, tdfo.taskRevision)

	currentTaskDefinitionJSON, _ := json.MarshalIndent(currentTaskDefinition.TaskDefinition, "", "  ")
	editedTaskDefinitionFile := fmt.Sprintf("tmp_%s_%d.json", tdfo.taskName, tdfo.taskRevision)

	newTaskDefinition, err := os.Create(editedTaskDefinitionFile)
	if err != nil {
		log.Panic(err)
	}

	defer newTaskDefinition.Close()

	_, err = newTaskDefinition.WriteString(string(currentTaskDefinitionJSON))
	if err != nil {
		log.Panic(err)
	}

	editorCommand := exec.Command(tdfo.editor, editedTaskDefinitionFile)
	editorCommand.Stdin = os.Stdin
	editorCommand.Stdout = os.Stdout

	err = editorCommand.Run()
	if err != nil {
		log.Panic(err)
	}

	readEditedTaskDefinitionFile, err := os.Open(editedTaskDefinitionFile)
	if err != nil {
		log.Panic(err)
	}

	defer readEditedTaskDefinitionFile.Close()

	newTaskDefinitionJSON, err := io.ReadAll(readEditedTaskDefinitionFile)
	if err != nil {
		log.Panic(err)
	}

	var newTasktaskInput ecs.RegisterTaskDefinitionInput

	fmt.Println(createDiff(string(currentTaskDefinitionJSON), string(newTaskDefinitionJSON)))

	if waitUserAcceptChanges() {
		err = json.Unmarshal([]byte(newTaskDefinitionJSON), &newTasktaskInput)
		if err != nil {
			log.Panic(err)
		}

		result, err := svc.RegisterTaskDefinition(&newTasktaskInput)
		if err != nil {
			log.Panic(err)
		}

		err = os.Remove(editedTaskDefinitionFile)
		if err != nil {
			log.Panic(err)
		}

		revision := result.TaskDefinition

		fmt.Printf("\n\nCreated revision %d for the task definition %s\n", *revision.Revision, tdfo.taskName)
	}
}
