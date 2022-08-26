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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Commands to update ECS services",
	Run:   updateServiceRun,
}

var (
	updateServiceName              string
	updateServiceTaskDef           string
	updateServiceEnableForceDeploy bool
)

func init() {
	servicesCmd.AddCommand(updateCmd)
	servicesCmd.PersistentFlags().StringVarP(&updateServiceName, "service", "s", "", "The name of the ECS service")
	servicesCmd.PersistentFlags().StringVarP(&updateServiceTaskDef, "task-definition", "t", "", "The name of the task definition")
	servicesCmd.PersistentFlags().BoolVarP(&updateServiceEnableForceDeploy, "force-new-deployment", "f", false, "Force new deployment")
}

func updateServiceRun(cmd *cobra.Command, args []string) {
	sess := provider.NewSession()
	svc := ecs.New(sess)

	input := &ecs.UpdateServiceInput{
		Cluster:            aws.String(servicesClusterName),
		Service:            aws.String(updateServiceName),
		TaskDefinition:     aws.String(updateServiceTaskDef),
		ForceNewDeployment: aws.Bool(updateServiceEnableForceDeploy),
	}

	_, err := svc.UpdateService(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
			case ecs.ErrCodeServiceNotFoundException:
				fmt.Println(ecs.ErrCodeServiceNotFoundException, aerr.Error())
			case ecs.ErrCodeServiceNotActiveException:
				fmt.Println(ecs.ErrCodeServiceNotActiveException, aerr.Error())
			case ecs.ErrCodePlatformUnknownException:
				fmt.Println(ecs.ErrCodePlatformUnknownException, aerr.Error())
			case ecs.ErrCodePlatformTaskDefinitionIncompatibilityException:
				fmt.Println(ecs.ErrCodePlatformTaskDefinitionIncompatibilityException, aerr.Error())
			case ecs.ErrCodeAccessDeniedException:
				fmt.Println(ecs.ErrCodeAccessDeniedException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Printf("Successful update of %s service for task definition %s\n", updateServiceName, updateServiceTaskDef)
}
