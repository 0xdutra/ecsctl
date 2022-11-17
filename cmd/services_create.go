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
	"github.com/spf13/viper"
)

type DeploymentConfiguration struct {
	MaximumPercent        int64 `mapstructure:"maximumPercent" yaml:"maximumPercent"`
	MinimumHealthyPercent int64 `mapstructure:"minimumHealthyPercent" yaml:"minimumHealthyPercent"`
}

type DeploymentController struct {
	Type string `mapstructure:"type" yaml:"type"`
}

type LoadBalancer struct {
	ContainerName  string `mapstructure:"containerName" yaml:"containerName"`
	ContainerPort  int64  `mapstructure:"containerPort" yaml:"containerPort"`
	TargetGroupArn string `mapstructure:"targetGroupArn" yaml:"targetGroupArn"`
}

type AwsVpcConfiguration struct {
	AssignPublicIp string   `mapstructure:"assignPublicIp" yaml:"assignPublicIp"`
	SecurityGroups []string `mapstructure:"securityGroups" yaml:"securityGroups"`
	Subnets        []string `mapstructure:"subnets" yaml:"subnets"`
}

type SvcConfig struct {
	Cluster                 string `mapstructure:"cluster" yaml:"cluster"`
	DesiredCount            int64  `mapstructure:"desiredCount" yaml:"desiredCount"`
	EnableECSManagedTags    bool   `mapstructure:"enableECSManagedTags" yaml:"enableECSManagedTags"`
	EnableExecuteCommand    bool   `mapstructure:"enableExecuteCommand" yaml:"enableExecuteCommand"`
	LaunchType              string `mapstructure:"launchType" yaml:"launchType"`
	SchedulingStrategy      string `mapstructure:"schedulingStrategy" yaml:"schedulingStrategy"`
	DeploymentController    DeploymentController
	DeploymentConfiguration DeploymentConfiguration
	AwsVpcConfiguration     AwsVpcConfiguration
	LoadBalancer            LoadBalancer
	ServiceName             string `mapstructure:"serviceName" yaml:"serviceName"`
	TaskDefinition          string `mapstructure:"taskDefinition" yaml:"taskDefinition"`
	Test                    string `mapstructure:"test" yaml:"test"`
}

type ServiceConfig struct {
	SvcConfig `mapstructure:"service" yaml:"service"`
}

func createService(configName string) {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var ecsService ServiceConfig

	err = viper.Unmarshal(&ecsService)
	if err != nil {
		log.Fatal(err)
	}

	sess := provider.NewSession()
	svc := ecs.New(sess)

	var service ecs.CreateServiceInput

	service.Cluster = aws.String(ecsService.Cluster)
	service.DesiredCount = aws.Int64(ecsService.DesiredCount)
	service.EnableECSManagedTags = aws.Bool(ecsService.EnableECSManagedTags)
	service.EnableExecuteCommand = aws.Bool(ecsService.EnableExecuteCommand)
	service.LaunchType = aws.String(ecsService.LaunchType)

	if viper.IsSet("service.awsVpcConfiguration") {
		service.NetworkConfiguration = &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				AssignPublicIp: aws.String(ecsService.AwsVpcConfiguration.AssignPublicIp),
				SecurityGroups: aws.StringSlice(ecsService.AwsVpcConfiguration.SecurityGroups),
				Subnets:        aws.StringSlice(ecsService.AwsVpcConfiguration.Subnets),
			},
		}
	}

	service.DeploymentConfiguration = &ecs.DeploymentConfiguration{
		MaximumPercent:        aws.Int64(ecsService.DeploymentConfiguration.MaximumPercent),
		MinimumHealthyPercent: aws.Int64(ecsService.DeploymentConfiguration.MinimumHealthyPercent),
	}

	service.DeploymentController = &ecs.DeploymentController{
		Type: aws.String(ecsService.DeploymentController.Type),
	}

	if viper.IsSet("service.loadBalancer") {
		service.LoadBalancers = []*ecs.LoadBalancer{
			{
				ContainerName:  aws.String(ecsService.LoadBalancer.ContainerName),
				ContainerPort:  aws.Int64(ecsService.LoadBalancer.ContainerPort),
				TargetGroupArn: aws.String(ecsService.LoadBalancer.TargetGroupArn),
			},
		}
	}

	service.SchedulingStrategy = aws.String(ecsService.SchedulingStrategy)
	service.ServiceName = aws.String(ecsService.ServiceName)
	service.TaskDefinition = aws.String(ecsService.TaskDefinition)

	result, err := svc.CreateService(&service)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(result)
}
