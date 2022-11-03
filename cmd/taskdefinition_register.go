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
	"ecsctl/pkg/colors"
	"ecsctl/pkg/provider"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/enescakir/emoji"
	"github.com/spf13/viper"
)

func environmentFiles(env map[string]string) []*ecs.EnvironmentFile {
	if env["type"] != "" && env["value"] != "" {
		return []*ecs.EnvironmentFile{
			{
				Type:  aws.String(env["type"]),
				Value: aws.String(env["value"]),
			},
		}
	}

	return []*ecs.EnvironmentFile{}
}

func logConfiguration(logDriver string, options map[string]string) *ecs.LogConfiguration {
	if logDriver != "" {
		return &ecs.LogConfiguration{
			LogDriver: aws.String(logDriver),
			Options:   aws.StringMap(options),
		}
	}

	return &ecs.LogConfiguration{}
}

func environment(envs map[string]string) []*ecs.KeyValuePair {
	var environment_vars []*ecs.KeyValuePair

	for key, value := range envs {
		environment_vars = append(environment_vars, &ecs.KeyValuePair{
			Name:  aws.String(key),
			Value: aws.String(value),
		})
	}

	return environment_vars
}

func registerTaskDefinition(configName string) {
	sess := provider.NewSession()
	svc := ecs.New(sess)

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var taskDefinition ecs.RegisterTaskDefinitionInput

	taskDefinition.ExecutionRoleArn = aws.String(viper.GetString("taskDefinition.executionRoleArn"))

	if viper.IsSet("taskDefinition.containerDefinitions") {
		taskDefinition.ContainerDefinitions = []*ecs.ContainerDefinition{
			{
				DnsSearchDomains:  aws.StringSlice(viper.GetStringSlice("taskDefinition.containerDefinitions.dnsSearchDomain")),
				EnvironmentFiles:  environmentFiles(viper.GetStringMapString("taskDefinition.containerDefinitions.environmentFiles")),
				LogConfiguration:  logConfiguration(viper.GetString("taskDefinition.containerDefinitions.logConfiguration.logDriver"), viper.GetStringMapString("taskDefinition.containerDefinitions.logConfiguration.options")),
				Name:              aws.String(viper.GetString("taskDefinition.containerDefinitions.name")),
				Image:             aws.String(viper.GetString("taskDefinition.containerDefinitions.image")),
				Cpu:               aws.Int64(viper.GetInt64("taskDefinition.containerDefinitions.cpu")),
				Memory:            aws.Int64(viper.GetInt64("taskDefinition.containerDefinitions.memory")),
				MemoryReservation: aws.Int64(viper.GetInt64("taskDefinition.containerDefinitions.memoryReservation")),
				PortMappings: []*ecs.PortMapping{
					{
						ContainerPort: aws.Int64(viper.GetInt64("taskDefinition.containerDefinitions.portsMappings.containerPort")),
						HostPort:      aws.Int64(viper.GetInt64("taskDefinition.containerDefinitions.portsMappings.hostPort")),
						Protocol:      aws.String(viper.GetString("taskDefinition.containerDefinitions.portsMappings.protocol")),
					},
				},
				Command:     aws.StringSlice(viper.GetStringSlice("taskDefinition.containerDefinitions.command")),
				Environment: environment(viper.GetStringMapString("taskDefinition.containerDefinitions.environment")),
			},
		}
	}

	taskDefinition.Cpu = aws.String(viper.GetString("taskDefinition.cpu"))
	taskDefinition.Memory = aws.String(viper.GetString("taskDefinition.memory"))
	taskDefinition.TaskRoleArn = aws.String(viper.GetString("taskDefinition.taskRoleArn"))
	taskDefinition.RequiresCompatibilities = aws.StringSlice(viper.GetStringSlice("taskDefinition.requiresCompatibilities"))
	taskDefinition.Family = aws.String(viper.GetString("taskDefinition.family"))
	taskDefinition.NetworkMode = aws.String(viper.GetString("taskDefinition.networkMode"))

	result, err := svc.RegisterTaskDefinition(&taskDefinition)
	if err != nil {
		fmt.Fprintf(colors.Out, "%s\n",
			colors.Red(fmt.Sprintf("%s  Register task definition\n", emoji.CrossMark)))

		log.Fatal(err)
	}

	r := result.TaskDefinition

	fmt.Fprintf(colors.Out, "%s\n",
		colors.Green(fmt.Sprintf("%s  Register task definition %s with revision %d\n",
			emoji.CheckMark,
			viper.GetString("taskDefinition.containerDefinitions.name"),
			*r.Revision)))
}
