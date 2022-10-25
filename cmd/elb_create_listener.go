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

// createElbListenerCmd represents the createListener command
var createElbListenerCmd = &cobra.Command{
	Use:   "create-listener",
	Short: "Commands to create Elastic Load Balancer listeners",
	Run:   createElbListenerRun,
}

func init() {
	elbCmd.AddCommand(createElbListenerCmd)
	createElbListenerCmd.PersistentFlags().StringVarP(&eo.loadBalancerListenerArn, "elb-arn", "", "", "The arn of the Elastic Load Balancer")
	createElbListenerCmd.PersistentFlags().StringVarP(&eo.loadBalancerListenerTgArn, "tg-arn", "", "", "The arn of the Target Group")
	createElbListenerCmd.PersistentFlags().Int64VarP(&eo.loadBalancerListenerPort, "tg-port", "", 80, "The port of the Target Group")
	createElbListenerCmd.PersistentFlags().StringVarP(&eo.loadBalancerListenerProtocol, "tg-protocol", "", "HTTP", "The protocol of the Target Group")

	if err := createElbListenerCmd.MarkPersistentFlagRequired("elb-arn"); err != nil {
		log.Fatal(err)
	}

	if err := createElbListenerCmd.MarkPersistentFlagRequired("tg-arn"); err != nil {
		log.Fatal(err)
	}
}

func createElbListenerRun(cmd *cobra.Command, args []string) {
	sess := provider.NewSession()
	svc := elbv2.New(sess)

	input := &elbv2.CreateListenerInput{
		DefaultActions: []*elbv2.Action{
			{
				TargetGroupArn: aws.String(eo.loadBalancerListenerTgArn),
				Type:           aws.String("forward"),
			},
		},
		LoadBalancerArn: aws.String(eo.loadBalancerListenerArn),
		Port:            aws.Int64(eo.loadBalancerListenerPort),
		Protocol:        aws.String(eo.loadBalancerListenerProtocol),
	}

	result, err := svc.CreateListener(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elbv2.ErrCodeDuplicateListenerException:
				fmt.Println(elbv2.ErrCodeDuplicateListenerException, aerr.Error())
			case elbv2.ErrCodeTooManyListenersException:
				fmt.Println(elbv2.ErrCodeTooManyListenersException, aerr.Error())
			case elbv2.ErrCodeTooManyCertificatesException:
				fmt.Println(elbv2.ErrCodeTooManyCertificatesException, aerr.Error())
			case elbv2.ErrCodeLoadBalancerNotFoundException:
				fmt.Println(elbv2.ErrCodeLoadBalancerNotFoundException, aerr.Error())
			case elbv2.ErrCodeTargetGroupNotFoundException:
				fmt.Println(elbv2.ErrCodeTargetGroupNotFoundException, aerr.Error())
			case elbv2.ErrCodeTargetGroupAssociationLimitException:
				fmt.Println(elbv2.ErrCodeTargetGroupAssociationLimitException, aerr.Error())
			case elbv2.ErrCodeInvalidConfigurationRequestException:
				fmt.Println(elbv2.ErrCodeInvalidConfigurationRequestException, aerr.Error())
			case elbv2.ErrCodeIncompatibleProtocolsException:
				fmt.Println(elbv2.ErrCodeIncompatibleProtocolsException, aerr.Error())
			case elbv2.ErrCodeSSLPolicyNotFoundException:
				fmt.Println(elbv2.ErrCodeSSLPolicyNotFoundException, aerr.Error())
			case elbv2.ErrCodeCertificateNotFoundException:
				fmt.Println(elbv2.ErrCodeCertificateNotFoundException, aerr.Error())
			case elbv2.ErrCodeUnsupportedProtocolException:
				fmt.Println(elbv2.ErrCodeUnsupportedProtocolException, aerr.Error())
			case elbv2.ErrCodeTooManyRegistrationsForTargetIdException:
				fmt.Println(elbv2.ErrCodeTooManyRegistrationsForTargetIdException, aerr.Error())
			case elbv2.ErrCodeTooManyTargetsException:
				fmt.Println(elbv2.ErrCodeTooManyTargetsException, aerr.Error())
			case elbv2.ErrCodeTooManyActionsException:
				fmt.Println(elbv2.ErrCodeTooManyActionsException, aerr.Error())
			case elbv2.ErrCodeInvalidLoadBalancerActionException:
				fmt.Println(elbv2.ErrCodeInvalidLoadBalancerActionException, aerr.Error())
			case elbv2.ErrCodeTooManyUniqueTargetGroupsPerLoadBalancerException:
				fmt.Println(elbv2.ErrCodeTooManyUniqueTargetGroupsPerLoadBalancerException, aerr.Error())
			case elbv2.ErrCodeALPNPolicyNotSupportedException:
				fmt.Println(elbv2.ErrCodeALPNPolicyNotSupportedException, aerr.Error())
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
