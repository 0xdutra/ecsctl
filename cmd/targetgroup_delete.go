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

// deleteTargetGroupCmd represents the delete command
var deleteTargetGroupCmd = &cobra.Command{
	Use:   "delete",
	Short: "Commands to delete target groups",
	Run:   deleteTargetGroupRun,
}

func init() {
	targetgroupCmd.AddCommand(deleteTargetGroupCmd)
	deleteTargetGroupCmd.PersistentFlags().StringVarP(&tgo.targetGroupArn, "arn", "", "", "The arn of the target group")

	if err := deleteTargetGroupCmd.MarkPersistentFlagRequired("arn"); err != nil {
		log.Fatal(err)
	}
}

func deleteTargetGroupRun(cmd *cobra.Command, args []string) {
	sess := provider.NewSession()
	svc := elbv2.New(sess)

	input := &elbv2.DeleteTargetGroupInput{
		TargetGroupArn: aws.String(tgo.targetGroupArn),
	}

	_, err := svc.DeleteTargetGroup(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elbv2.ErrCodeResourceInUseException:
				fmt.Println(elbv2.ErrCodeResourceInUseException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Printf("%s successfully deleted\n", tgo.targetGroupArn)
}
