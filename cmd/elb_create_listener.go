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

/*
// createElbListenerCmd represents the createListener command
var createElbListenerCmd = &cobra.Command{
	Use:   "create-listener",
	Short: "Commands to create Elastic Load Balancer listeners",
	Run:   createElbListenerRun,
}

func init() {
	elbCmd.AddCommand(createElbListenerCmd)
	createElbListenerCmd.PersistentFlags().StringVarP(&eo.loadBalancerListenerArn, "elb-arn", "", "", "The arn of the Elastic Load Balancer")
	createElbListenerCmd.PersistentFlags().StringVarP(&eo.loadBalancerListenerTgArn, "tg-arn", "", "", "The arn of the Target Group")
	createElbListenerCmd.PersistentFlags().Int64VarP(&eo.loadBalancerListenerPort, "tg-port", "", 80, "The port of the Target Group")
	createElbListenerCmd.PersistentFlags().StringVarP(&eo.loadBalancerListenerProtocol, "tg-protocol", "", "HTTP", "The protocol of the Target Group")

	if err := createElbListenerCmd.MarkPersistentFlagRequired("elb-arn"); err != nil {
		log.Fatal(err)
	}

	if err := createElbListenerCmd.MarkPersistentFlagRequired("tg-arn"); err != nil {
		log.Fatal(err)
	}
}
*/
