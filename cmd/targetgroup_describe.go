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

// describeTargetGroupsCmd represents the describe command
var describeTargetGroupsCmd = &cobra.Command{
	Use:   "describe",
	Short: "Commands to describe target groups",
	Run:   describeTargetGroupsRun,
}

func init() {
	targetgroupCmd.AddCommand(describeTargetGroupsCmd)
	describeTargetGroupsCmd.PersistentFlags().StringVarP(&tgo.targetGroupName, "name", "", "", "The name of the target group")

	if err := describeTargetGroupsCmd.MarkPersistentFlagRequired("name"); err != nil {
		log.Fatal(err)
	}

}

func describeTargetGroupsRun(cmd *cobra.Command, args []string) {
	sess := provider.NewSession()
	svc := elbv2.New(sess)

	input := &elbv2.DescribeTargetGroupsInput{
		Names: []*string{
			aws.String(tgo.targetGroupName),
		},
	}

	result, err := svc.DescribeTargetGroups(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elbv2.ErrCodeLoadBalancerNotFoundException:
				fmt.Println(elbv2.ErrCodeLoadBalancerNotFoundException, aerr.Error())
			case elbv2.ErrCodeTargetGroupNotFoundException:
				fmt.Println(elbv2.ErrCodeTargetGroupNotFoundException, aerr.Error())
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
