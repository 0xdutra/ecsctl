/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"ecsctl/pkg/provider"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
)

// deleteServiceCmd represents the delete command
var deleteServiceCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Commands to delete ECS service",
	Example: "ecsctl service delete --name <name>",
	Run:     deleteServiceRun,
}

func init() {
	servicesCmd.AddCommand(deleteServiceCmd)
	deleteServiceCmd.PersistentFlags().StringVarP(&so.serviceName, "service", "", "", "The name of the ECS service")
}

func deleteServiceRun(cmd *cobra.Command, args []string) {
	sess := provider.NewSession()
	svc := ecs.New(sess)

	input := &ecs.DeleteServiceInput{
		Service: aws.String(so.serviceName),
	}

	result, err := svc.DeleteService(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
			case ecs.ErrCodeServiceNotFoundException:
				fmt.Println(ecs.ErrCodeServiceNotFoundException, aerr.Error())
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
