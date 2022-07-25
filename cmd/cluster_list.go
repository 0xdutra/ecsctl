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

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
)

// listClustersCmd represents the listClusters command
var listClustersCmd = &cobra.Command{
	Use:   "list-clusters",
	Short: "Command to list all ECS clusters",
	Run:   listClusters,
}

func init() {
	clusterCmd.AddCommand(listClustersCmd)
}

func listClusters(cmd *cobra.Command, _ []string) {
	sess := provider.NewSession()
	svc := ecs.New(sess)

	var nextToken string

	for {
		result, err := svc.ListClusters(&ecs.ListClustersInput{
			NextToken: &nextToken,
		})

		if err != nil {
			log.Panic(err)
		}

		for _, clusterArn := range result.ClusterArns {
			fmt.Println(*clusterArn)
		}

		if result.NextToken == nil {
			break
		}

		nextToken = *result.NextToken
	}
}
