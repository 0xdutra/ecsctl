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

func logConfiguration(logDriver string, options map[string]string) *ecs.LogConfiguration {
	if logDriver != "" {
		return &ecs.LogConfiguration{
			LogDriver: aws.String(logDriver),
			Options:   aws.StringMap(options),
		}
	}

	return nil
}

type TaskDefinitionConfig struct {
	TaskDefConfig `mapstructure:"taskDefinition" yaml:"taskDefinition"`
}

type LogConfiguration struct {
	LogDriver string            `mapstructure:"logDriver" yaml:"logDriver"`
	Options   map[string]string `mapstructure:"options" yaml:"options"`
}

type PortMappings struct {
	ContainerPort int64  `mapstructure:"containerPort" yaml:"containerPort"`
	HostPort      int64  `mapstructure:"hostPort" yaml:"hostPort"`
	Protocol      string `mapstructure:"protocol" yaml:"protocol"`
}

type ContainerDefinitions struct {
	DnsSearchDomains  []string          `mapstructure:"dnsSearchDomains" yaml:"dnsSearchDomains"`
	EnvironmentFiles  map[string]string `mapstructure:"environmentFiles" yaml:"environmentFiles"`
	LogConfiguration  LogConfiguration  `mapstructure:"logConfiguration" yaml:"logConfiguration"`
	Name              string            `mapstructure:"name" yaml:"name"`
	Image             string            `mapstructure:"image" yaml:"image"`
	Cpu               int64             `mapstructure:"cpu" yaml:"cpu"`
	Memory            int64             `mapstructure:"memory" yaml:"memory"`
	MemoryReservation int64             `mapstructure:"memoryReservation" yaml:"memoryReservation"`
	Commands          []string          `mapstructure:"commands" yaml:"commands"`
	Environment       map[string]string `mapstructure:"environment" yaml:"environment"`
	PortsMappings     PortMappings      `mapstructure:"portsMappings" yaml:"portsMappings"`
}

type TaskDefConfig struct {
	ExecutionRoleArn        string `mapstructure:"executionRoleArn" yaml:"executionRoleArn"`
	ContainerDefinitions    ContainerDefinitions
	Cpu                     string   `mapstructure:"cpu" yaml:"cpu"`
	Memory                  string   `mapstructure:"memory" yaml:"memory"`
	TaskRoleArn             string   `mapstructure:"taskRoleArn" yaml:"taskRoleArn"`
	RequiresCompatibilities []string `mapstructure:"requiresCompatibilities" yaml:"requiresCompatibilities"`
	Family                  string   `mapstructure:"family" yaml:"family"`
	NetworkMode             string   `mapstructure:"networkMode" yaml:"networkMode"`
}

func registerTaskDefinition(configName string) {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var taskDefConfig TaskDefinitionConfig

	err = viper.Unmarshal(&taskDefConfig)
	if err != nil {
		log.Panic(err)
	}

	var taskDefinition ecs.RegisterTaskDefinitionInput

	taskDefinition.ExecutionRoleArn = &taskDefConfig.ExecutionRoleArn

	if viper.IsSet("taskDefinition.containerDefinitions") {
		taskDefinition.ContainerDefinitions = []*ecs.ContainerDefinition{
			{
				DnsSearchDomains: aws.StringSlice(taskDefConfig.ContainerDefinitions.DnsSearchDomains),
				EnvironmentFiles: environmentFiles(taskDefConfig.ContainerDefinitions.EnvironmentFiles),

				LogConfiguration: logConfiguration(
					taskDefConfig.ContainerDefinitions.LogConfiguration.LogDriver,
					taskDefConfig.ContainerDefinitions.LogConfiguration.Options,
				),

				Name:              aws.String(taskDefConfig.ContainerDefinitions.Name),
				Image:             aws.String(taskDefConfig.ContainerDefinitions.Image),
				Cpu:               aws.Int64(taskDefConfig.ContainerDefinitions.Cpu),
				Memory:            aws.Int64(taskDefConfig.ContainerDefinitions.Memory),
				MemoryReservation: aws.Int64(taskDefConfig.ContainerDefinitions.MemoryReservation),
				PortMappings: []*ecs.PortMapping{
					{
						ContainerPort: aws.Int64(taskDefConfig.ContainerDefinitions.PortsMappings.ContainerPort),
						HostPort:      aws.Int64(taskDefConfig.ContainerDefinitions.PortsMappings.HostPort),
						Protocol:      aws.String(taskDefConfig.ContainerDefinitions.PortsMappings.Protocol),
					},
				},
				Command:     aws.StringSlice(taskDefConfig.ContainerDefinitions.Commands),
				Environment: environment(taskDefConfig.ContainerDefinitions.Environment),
			},
		}
	}

	taskDefinition.Cpu = aws.String(taskDefConfig.Cpu)
	taskDefinition.Memory = aws.String(taskDefConfig.Memory)
	taskDefinition.TaskRoleArn = aws.String(taskDefConfig.TaskRoleArn)
	taskDefinition.RequiresCompatibilities = aws.StringSlice(taskDefConfig.RequiresCompatibilities)
	taskDefinition.Family = aws.String(taskDefConfig.Family)
	taskDefinition.NetworkMode = aws.String(taskDefConfig.NetworkMode)

	fmt.Println(taskDefinition)

	sess := provider.NewSession()
	svc := ecs.New(sess)

	result, err := svc.RegisterTaskDefinition(&taskDefinition)
	if err != nil {
		fmt.Println("[-] - Task definition registration failed")
		log.Fatal(err)
	}

	r := result.TaskDefinition

	fmt.Printf("[+] - Register task definition with revision %d\n", *r.Revision)
}
