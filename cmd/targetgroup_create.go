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
	"github.com/spf13/viper"
)

type TargetGroupConfig struct {
	TGConfig `mapstructure:"targetGroup" yaml:"targetGroup"`
}

type TGConfig struct {
	Name                       string `mapstructure:"name" yaml:"name"`
	Port                       int64  `mapstructure:"port" yaml:"port"`
	Protocol                   string `mapstructure:"protocol" yaml:"protocol"`
	VpcId                      string `mapstructure:"vpcId" yaml:"vpcId"`
	TargetType                 string `mapstructure:"targetType" yaml:"targetType"`
	HealthCheckEnabled         bool   `mapstructure:"healthCheckEnabled" yaml:"healthCheckEnabled"`
	HealthCheckIntervalSeconds int64  `mapstructure:"healthCheckIntervalSeconds" yaml:"healthCheckIntervalSeconds"`
	HealthCheckPath            string `mapstructure:"healthCheckPath" yaml:"healthCheckPath"`
}

func createTargetGroup(configName string) {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var tg TargetGroupConfig

	err = viper.Unmarshal(&tg)
	if err != nil {
		log.Fatal(err)
	}

	sess := provider.NewSession()
	svc := elbv2.New(sess)

	input := &elbv2.CreateTargetGroupInput{
		Name:                       aws.String(tg.Name),
		Port:                       aws.Int64(tg.Port),
		Protocol:                   aws.String(tg.Protocol),
		VpcId:                      aws.String(tg.VpcId),
		TargetType:                 aws.String(tg.TargetType),
		HealthCheckEnabled:         aws.Bool(tg.HealthCheckEnabled),
		HealthCheckIntervalSeconds: aws.Int64(tg.HealthCheckIntervalSeconds),
		HealthCheckPath:            aws.String(tg.HealthCheckPath),
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
