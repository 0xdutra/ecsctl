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
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
)

// createTaskdefinitionCmd represents the createTaskdefinition command
var createTaskdefinitionCmd = &cobra.Command{
	Use:   "create-task-definition",
	Short: "Commands to create task definition",
	Run:   createTaskdefinitionRun,
}

func init() {
	taskDefinitionCmd.AddCommand(createTaskdefinitionCmd)
}

func parserLogOptions(logOptions *string) map[string]*string {
	newStr := strings.Split(*logOptions, "&")
	options := make(map[string]*string)

	for _, pair := range newStr {
		z := strings.Split(pair, "=")
		options[z[0]] = &z[1]
	}

	return options
}

func createTaskdefinitionRun(cmd *cobra.Command, _ []string) {
	sess := provider.NewSession()
	svc := ecs.New(sess)

	input := &ecs.RegisterTaskDefinitionInput{

		ContainerDefinitions: []*ecs.ContainerDefinition{
			{
				Cpu:       aws.Int64(cpu),
				Essential: aws.Bool(true),
				Image:     aws.String(image),
				Memory:    aws.Int64(memory),
				Name:      aws.String(name),
				LogConfiguration: &ecs.LogConfiguration{
					LogDriver: aws.String(logDriver),
					Options:   parserLogOptions(&logOptions),
				},
			},
		},
		Family:      aws.String(family),
		TaskRoleArn: aws.String(taskRoleArn),
	}

	result, err := svc.RegisterTaskDefinition(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)

}
