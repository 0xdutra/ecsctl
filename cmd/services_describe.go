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
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
)

// describeServicesCmd represents the describeServices command
var describeServicesCmd = &cobra.Command{
	Use:   "describe",
	Short: "Commands to describe ECS services",
	Run:   describeServicesRun,
}

var (
	serviceArn string
)

func init() {
	servicesCmd.AddCommand(describeServicesCmd)
	describeServicesCmd.Flags().StringVarP(&serviceArn, "service-arn", "", "", "The ARN of the ECS service")
}

func describeServicesRun(cmd *cobra.Command, args []string) {
	sess := provider.NewSession()
	svc := ecs.New(sess)

	result, err := svc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster: aws.String(servicesClusterName),
		Services: []*string{
			aws.String(serviceArn),
		},
	})

	if err != nil {
		log.Panic(err)
	}

	fmt.Println(result.Services)

	for _, svc := range result.Services {
		fmt.Printf("\nClusterArn: %s\n", *svc.ClusterArn)
		fmt.Printf("CreatedAt: %s\n", *svc.CreatedAt)
		fmt.Printf("CreatedBy: %s\n", *svc.CreatedBy)
		fmt.Printf("DeploymentConfiguration: %s\n", *svc.DeploymentConfiguration)
		fmt.Printf("Deployments: %s\n", svc.Deployments)
		fmt.Printf("DesiredCount: %d\n", *svc.DesiredCount)
		fmt.Printf("PendingCount: %d\n", *svc.PendingCount)
		fmt.Printf("PropagateTags: %s\n", *svc.PropagateTags)
		fmt.Printf("EnableECSManagedTags: %v\n", *svc.EnableECSManagedTags)
		fmt.Printf("LaunchType: %s\n", *svc.LaunchType)
		fmt.Printf("LoadBalancers: %s\n", svc.LoadBalancers)
		fmt.Printf("PlacementConstraints: %s\n", svc.PlacementConstraints)
		fmt.Printf("PlacementStrategy: %s\n", svc.PlacementStrategy)
		fmt.Printf("PlacementConstraints: %s\n", svc.PlacementConstraints)
		fmt.Printf("SchedulingStrategy: %s\n", *svc.SchedulingStrategy)
		fmt.Printf("ServiceArn: %s\n", *svc.ServiceArn)
		fmt.Printf("ServiceName: %s\n", *svc.ServiceName)
		fmt.Printf("ServiceRegistries: %s\n", svc.ServiceRegistries)
		fmt.Printf("Status: %s\n", *svc.Status)
		fmt.Printf("TaskDefinition: %s\n", *svc.TaskDefinition)
	}
}
