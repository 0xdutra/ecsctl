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

// servicesEventsCmd represents the status command
var servicesEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Get events of the ECS service",
	Run:   getervicesEventsRun,
}

var maxEventsResults int

func init() {
	servicesCmd.AddCommand(servicesEventsCmd)
	servicesEventsCmd.PersistentFlags().StringVarP(&so.clusterName, "cluster", "", "", "The name of the ECS cluster")
	servicesEventsCmd.PersistentFlags().StringVarP(&so.serviceName, "service", "", "", "The name of the ECS service")
	servicesEventsCmd.PersistentFlags().IntVarP(&maxEventsResults, "max-result", "", 10, "Maximum results to deplay")

	if err := servicesEventsCmd.MarkPersistentFlagRequired("cluster"); err != nil {
		log.Panic(err)
	}

	if err := servicesEventsCmd.MarkPersistentFlagRequired("service"); err != nil {
		log.Fatal(err)
	}
}

func getervicesEventsRun(cmd *cobra.Command, args []string) {
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

	events := result.Services[0].Events

	for _, event := range events {
		if maxEventsResults == 0 {
			break
		}

		fmt.Printf("CreatedAt: %s\n", event.CreatedAt)
		fmt.Printf("Id: %s\n", *event.Id)
		fmt.Printf("Message: %s\n\n", *event.Message)
		maxEventsResults--
	}
}
