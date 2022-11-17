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

type LoadBalancerConfig struct {
	ELBConfig `mapstructure:"loadBalancers" yaml:"loadBalancers"`
}

type LoadBalancerListener struct {
	DefaultActions map[string]string `mapstructure:"defaultActions" yaml:"defaultActions"`
	Port           int64             `mapstructure:"port" yaml:"port"`
	Protocol       string            `mapstructure:"protocol" yaml:"protocol"`
}

type ELBConfig struct {
	Name     string               `mapstructure:"name" yaml:"name"`
	Subnets  []string             `mapstructure:"subnets" yaml:"subnets"`
	Type     string               `mapstructure:"type" yaml:"type"`
	Scheme   string               `mapstructure:"scheme" yaml:"scheme"`
	Listener LoadBalancerListener `mapstructure:"listener" yaml:"listener"`
}

func createElb(configName string) {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var elb LoadBalancerConfig

	err = viper.Unmarshal(&elb)
	if err != nil {
		log.Fatal(err)
	}

	sess := provider.NewSession()
	svc := elbv2.New(sess)

	loadBalancerOutput, err := svc.CreateLoadBalancer(&elbv2.CreateLoadBalancerInput{
		Name:    aws.String(elb.Name),
		Subnets: aws.StringSlice(elb.Subnets),
		Type:    aws.String(elb.Type),
		Scheme:  aws.String(elb.Scheme),
	})

	if err != nil {
		log.Fatal(err)
	}

	loadBalancerArn := loadBalancerOutput.LoadBalancers[0].LoadBalancerArn
	loadBalancerListenerOutput, err := svc.CreateListener(&elbv2.CreateListenerInput{
		LoadBalancerArn: loadBalancerArn,
		Port:            aws.Int64(elb.Listener.Port),
		Protocol:        aws.String(elb.Listener.Protocol),

		DefaultActions: []*elbv2.Action{
			{
				TargetGroupArn: aws.String(elb.Listener.DefaultActions["targetgrouparn"]),
				Type:           aws.String(elb.Listener.DefaultActions["type"]),
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(loadBalancerListenerOutput)
}
