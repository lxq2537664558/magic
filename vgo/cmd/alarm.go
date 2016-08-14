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

	"github.com/corego/vgo/vgo/alarm"
	"github.com/spf13/cobra"
)

// alarmCmd represents the alarm command
var alarmCmd = &cobra.Command{
	Use:   "alarm",
	Short: "A brief description of your command",
	Long:  ``,
	Run:   alarmrun,
}

func init() {
	RootCmd.AddCommand(alarmCmd)
}

func alarmrun(cmd *cobra.Command, args []string) {
	a := alarm.New()
	// start alarm server
	go a.Start()
	// wait server stop signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("alarm service received Signal: ", <-chSig)
	fmt.Println("alarm service is going to stop")
	// close alarm server
	a.Close()
	fmt.Println("alarm service is stopped")
}
