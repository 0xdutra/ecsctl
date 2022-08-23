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
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/spf13/cobra"
)

// createTargetGroupCmd represents the create command
var createTargetGroupCmd = &cobra.Command{
	Use:   "create",
	Short: "Commands to create target groups",
	Run:   createTargetGroupRun,
}

var (
	targetGroupName                string
	targetGroupPort                int64
	targetGroupProtocol            string
	targetGroupVpcID               string
	targetGroupType                string
	targetGroupEnableHealthCheck   bool
	targetGroupHealthCheckPath     string
	targetGroupHealthCheckInterval int64
)

func init() {
	targetgroupCmd.AddCommand(createTargetGroupCmd)
	createTargetGroupCmd.PersistentFlags().StringVarP(&targetGroupName, "name", "", "", "The name of the target group")
	createTargetGroupCmd.PersistentFlags().Int64VarP(&targetGroupPort, "port", "", 80, "The port of the target group")
	createTargetGroupCmd.PersistentFlags().StringVarP(&targetGroupProtocol, "protocol", "", "HTTP", "The protocol of the target group")
	createTargetGroupCmd.PersistentFlags().StringVarP(&targetGroupVpcID, "vpcid", "", "", "The vpcid of the target group")
	createTargetGroupCmd.PersistentFlags().StringVarP(&targetGroupType, "type", "", "ip", "The type of the target group")
	createTargetGroupCmd.PersistentFlags().BoolVarP(&targetGroupEnableHealthCheck, "healthcheck", "", true, "Enable health check in target group")
	createTargetGroupCmd.PersistentFlags().StringVarP(&targetGroupHealthCheckPath, "healthcheck-path", "", "/", "Health check path")
	createTargetGroupCmd.PersistentFlags().Int64VarP(&targetGroupHealthCheckInterval, "healthcheck-interval", "", 30, "Health check interval in seconds")

	if err := createTargetGroupCmd.MarkPersistentFlagRequired("name"); err != nil {
		log.Fatal(err)
	}

	if err := createTargetGroupCmd.MarkPersistentFlagRequired("vpcid"); err != nil {
		log.Fatal(err)
	}
}

func createTargetGroupRun(cmd *cobra.Command, args []string) {
	sess := provider.NewSession()
	svc := elbv2.New(sess)

	input := &elbv2.CreateTargetGroupInput{
		Name:                       aws.String(targetGroupName),
		Port:                       aws.Int64(targetGroupPort),
		Protocol:                   aws.String(targetGroupProtocol),
		VpcId:                      aws.String(targetGroupVpcID),
		TargetType:                 aws.String(targetGroupType),
		HealthCheckEnabled:         aws.Bool(targetGroupEnableHealthCheck),
		HealthCheckIntervalSeconds: aws.Int64(targetGroupHealthCheckInterval),
		HealthCheckPath:            aws.String(targetGroupHealthCheckPath),
	}

	result, err := svc.CreateTargetGroup(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elbv2.ErrCodeDuplicateTargetGroupNameException:
				fmt.Println(elbv2.ErrCodeDuplicateTargetGroupNameException, aerr.Error())
			case elbv2.ErrCodeTooManyTargetGroupsException:
				fmt.Println(elbv2.ErrCodeTooManyTargetGroupsException, aerr.Error())
			case elbv2.ErrCodeInvalidConfigurationRequestException:
				fmt.Println(elbv2.ErrCodeInvalidConfigurationRequestException, aerr.Error())
			case elbv2.ErrCodeTooManyTagsException:
				fmt.Println(elbv2.ErrCodeTooManyTagsException, aerr.Error())
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
