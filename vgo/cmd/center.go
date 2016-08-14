// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/corego/vgo/vgo/center"
	"github.com/spf13/cobra"
)

// centerCmd represents the center command
var centerCmd = &cobra.Command{
	Use:   "center",
	Short: "A brief description of your command",
	Long:  ``,
	Run:   centerrun,
}

func init() {
	RootCmd.AddCommand(centerCmd)
}

func centerrun(cmd *cobra.Command, args []string) {
	c := center.New()
	// start center server
	go c.Start()
	// wait server stop signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("center service received Signal: ", <-chSig)
	fmt.Println("center service is going to stop")
	// close center server
	c.Close()
	fmt.Println("center service is stopped")
}
