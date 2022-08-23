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
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/spf13/cobra"
)

// elbCreateCmd represents the create command
var elbCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Commands to create Elastic Load Balancer",
	Run:   createElbRun,
}

var (
	subnets            []string
	loadBalancerName   string
	loadBalancerScheme string
	loadBalancerType   string
)

func init() {
	elbCmd.AddCommand(elbCreateCmd)
	elbCreateCmd.PersistentFlags().StringVarP(&loadBalancerName, "name", "", "", "The name of the Elastic Load Balancer")
	elbCreateCmd.PersistentFlags().StringVarP(&loadBalancerScheme, "scheme", "", "internet-facing", "Specify a scheme for a Elastic Load Balancer")
	elbCreateCmd.PersistentFlags().StringArrayVarP(&subnets, "subnet", "", nil, "The list of subnets ids")
	elbCreateCmd.PersistentFlags().StringVarP(&loadBalancerType, "type", "", "application", "The type of Elastic Load Balancer")

	if err := elbCreateCmd.MarkPersistentFlagRequired("name"); err != nil {
		log.Fatal(err)
	}

	if err := elbCreateCmd.MarkPersistentFlagRequired("subnet"); err != nil {
		log.Fatal(err)
	}
}

func createElbRun(cmd *cobra.Command, args []string) {
	var awsSubnets []*string
	for _, v := range subnets {
		awsSubnets = append(awsSubnets, aws.String(v))
	}

	sess := provider.NewSession()
	svc := elbv2.New(sess)

	input := &elbv2.CreateLoadBalancerInput{
		Name:    aws.String(loadBalancerName),
		Subnets: awsSubnets,
		Type:    aws.String(loadBalancerType),
		Scheme:  aws.String(loadBalancerScheme),
	}

	result, err := svc.CreateLoadBalancer(input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
