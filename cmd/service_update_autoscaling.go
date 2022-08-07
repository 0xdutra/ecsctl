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
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
)

// updateServiceCmd represents the updateService command
var updateServiceCmd = &cobra.Command{
	Use:   "update-autoscaling",
	Short: "Commands to update ECS services",
	Run:   updateServiceRun,
}

var (
	minCapacity     int64
	maxCapacity     int64
	desiredCapacity int64
	serviceName     string
)

func init() {
	servicesCmd.AddCommand(updateServiceCmd)
	servicesCmd.PersistentFlags().Int64VarP(&minCapacity, "min", "", 2, "The lower boundary to which Service Auto Scaling can adjust your service’s desired count")
	servicesCmd.PersistentFlags().Int64VarP(&maxCapacity, "max", "", 2, "The upper boundary to which Service Auto Scaling can adjust your service’s desired count")
	servicesCmd.PersistentFlags().Int64VarP(&desiredCapacity, "desired", "", 2, "The initial desired count to start with before Service Auto Scaling begins adjustment")
	servicesCmd.PersistentFlags().StringVarP(&serviceName, "service-name", "", "", "The name of the ECS service")
}

func updateServiceRun(cmd *cobra.Command, args []string) {
	sess := provider.NewSession()
	svc := ecs.New(sess)

	err := updateAutoScalingRequirements(servicesClusterName, serviceName, maxCapacity, minCapacity)
	if err != nil {
		log.Panic(err)
	}

	_, err = svc.UpdateService(&ecs.UpdateServiceInput{
		Cluster:      aws.String(servicesClusterName),
		DesiredCount: aws.Int64(desiredCapacity),
		Service:      aws.String(serviceName),
	})

	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Updating capacity %s to minimum %d, desired %d and maximum %d\n", serviceName, minCapacity, desiredCapacity, maxCapacity)
}

func updateAutoScalingRequirements(clusterName string, serviceName string, maxCapacity int64, minCapacity int64) error {
	sess := provider.NewSession()
	svc := applicationautoscaling.New(sess)

	resourceId := fmt.Sprintf("service/%s/%s", clusterName, serviceName)

	_, err := svc.RegisterScalableTarget(&applicationautoscaling.RegisterScalableTargetInput{
		MaxCapacity:       aws.Int64(maxCapacity),
		MinCapacity:       aws.Int64(minCapacity),
		ResourceId:        aws.String(resourceId),
		ScalableDimension: aws.String("ecs:service:DesiredCount"),
		ServiceNamespace:  aws.String("ecs"),
	})

	if err != nil {
		return err
	}

	return err
}
