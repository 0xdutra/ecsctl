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
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
)

// servicesStatusCmd represents the status command
var servicesStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Commands to check the status of services",
	Run:   statusServicesRun,
}

func init() {
	servicesCmd.AddCommand(servicesStatusCmd)
	servicesStatusCmd.PersistentFlags().StringVarP(&so.clusterName, "cluster", "", "", "The name of the ECS cluster")
	servicesStatusCmd.PersistentFlags().StringVarP(&so.serviceName, "service", "", "", "The arn or name of the ECS service")

	if err := servicesStatusCmd.MarkPersistentFlagRequired("cluster"); err != nil {
		log.Panic(err)
	}

	if err := servicesStatusCmd.MarkPersistentFlagRequired("service"); err != nil {
		log.Panic(err)
	}
}

func statusServicesRun(cmd *cobra.Command, args []string) {
	sess := provider.NewSession()
	svc := ecs.New(sess)

	input := &ecs.DescribeServicesInput{
		Cluster: aws.String(so.clusterName),
		Services: []*string{
			aws.String(so.serviceName),
		},
	}

	result, err := svc.DescribeServices(input)
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
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	for _, status := range result.Services {
		fmt.Printf("Services status: %s\n\tRunning %d/%d tasks (%d pending)\n\n",
			*status.Status,
			*status.DesiredCount,
			*status.RunningCount,
			*status.PendingCount)

		for _, deployment := range status.Deployments {
			fmt.Printf("Rollout state: %s (Id: %s)\n\tFailed tasks: %d\n\tLaunch type: %s\n\tTask Definition: %s\n\tUpdated At: %s\n",
				*deployment.RolloutState,
				*deployment.Id,
				*deployment.FailedTasks,
				*deployment.LaunchType,
				*deployment.TaskDefinition,
				*deployment.UpdatedAt)
		}

		latest_event := status.Events[0]
		fmt.Printf("Latest event: %s\nId: %s\nMessage: %s\n", *latest_event.CreatedAt, *latest_event.Id, *latest_event.Message)
	}
}
