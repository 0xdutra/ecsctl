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

func createService(configName string) {
	sess := provider.NewSession()
	svc := ecs.New(sess)

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var service ecs.CreateServiceInput

	service.Cluster = aws.String(viper.GetString("service.clusterName"))
	service.DesiredCount = aws.Int64(viper.GetInt64("service.desiredCount"))
	service.EnableECSManagedTags = aws.Bool(viper.GetBool("service.enableECSManagedTags"))
	service.EnableExecuteCommand = aws.Bool(viper.GetBool("service.enableExecuteCommand"))
	service.LaunchType = aws.String(viper.GetString("service.launchType"))

	service.DeploymentConfiguration = &ecs.DeploymentConfiguration{
		MaximumPercent:        aws.Int64(viper.GetInt64("service.deploymentConfiguration.maximumPercent")),
		MinimumHealthyPercent: aws.Int64(viper.GetInt64("service.deploymentConfiguration.minimumPercent")),
	}

	service.DeploymentController = &ecs.DeploymentController{
		Type: aws.String(viper.GetString("service.deploymentController.type")),
	}

	if viper.IsSet("service.loadBalancers") {
		service.LoadBalancers = []*ecs.LoadBalancer{
			{
				ContainerName:  aws.String(viper.GetString("service.loadBalancers.containerName")),
				ContainerPort:  aws.Int64(viper.GetInt64("service.loadBalancers.containerPort")),
				TargetGroupArn: aws.String(viper.GetString("service.loadBalancers.targetGroupArn")),
			},
		}
	}

	if viper.IsSet("service.networkConfiguration") {
		service.NetworkConfiguration = &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				AssignPublicIp: aws.String(viper.GetString("service.networkConfiguration.awsVpcConfiguration.assignPublicIp")),
				SecurityGroups: aws.StringSlice(viper.GetStringSlice("service.networkConfiguration.awsVpcConfiguration.securityGroups")),
				Subnets:        aws.StringSlice(viper.GetStringSlice("service.networkConfiguration.awsVpcConfiguration.subnets")),
			},
		}
	}

	service.SchedulingStrategy = aws.String(viper.GetString("service.schedulingStrategy"))
	service.ServiceName = aws.String(viper.GetString("service.serviceName"))
	service.TaskDefinition = aws.String(viper.GetString("service.taskDefinition"))

	result, err := svc.CreateService(&service)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(result)

	/*
		sess := provider.NewSession()
		svc := ecs.New(sess)

		serviceFile, err := os.Open(so.serviceInputJson)
		if err != nil {
			log.Panic(err)
		}

		defer serviceFile.Close()

		readServiceFile, err := ioutil.ReadAll(serviceFile)
		if err != nil {
			log.Panic(err)
		}

		var serviceInput ecs.CreateServiceInput

		err = json.Unmarshal([]byte(readServiceFile), &serviceInput)
		if err != nil {
			log.Panic(err)
		}
	*/
}
