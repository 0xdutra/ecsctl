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
	"github.com/spf13/viper"
)

/*
func createElbListener(configName string) {
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
		log.Fatal(err)
	}

	fmt.Println(result)
}
*/

func createElb(configName string) {
	sess := provider.NewSession()
	svc := elbv2.New(sess)

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var loadBalancer elbv2.CreateLoadBalancerInput
	//var loadBalancerListener elbv2.CreateListenerInput

	loadBalancer.Name = aws.String(viper.GetString("loadBalancers.name"))
	loadBalancer.Subnets = aws.StringSlice(viper.GetStringSlice("loadBalancers.subnets"))
	loadBalancer.Type = aws.String(viper.GetString("loadBalancers.type"))
	loadBalancer.Scheme = aws.String(viper.GetString("loadBalancers.scheme"))

	loadBalancerOutput, err := svc.CreateLoadBalancer(&loadBalancer)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(loadBalancerOutput)

	/*
		loadBalancerListener.DefaultActions = []*elbv2.Action{
			{
				TargetGroupArn: aws.String(loadBalancerOutput.LoadBalancers),
				Type:           aws.String("forward"),
			},
		}

		loadBalancerListener.LoadBalancerArn = aws.String("")
		loadBalancerListener.Port = aws.Int64(80)
		loadBalancerListener.Protocol = aws.String("")

		loadBalancerListenerOutput, err := svc.CreateListener(&loadBalancerListener)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(loadBalancerListenerOutput)
	*/
}
