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

	"github.com/aiyun/openapm/mecury/agent"
	_ "github.com/aiyun/openapm/mecury/plugins/input/all"
	_ "github.com/aiyun/openapm/mecury/plugins/output/all"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mecury",
	Short: "the agent of vgo platform",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: start,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}

func start(cmd *cobra.Command, args []string) {
	isReload := true

	for isReload {
		isReload = false

		// init config
		agent.LoadConfig()
		// init logger

		ag := agent.New()
		ag.Init()

		// agent shutdown signal
		shutdown := make(chan struct{})

		// catch system signals
		signals := make(chan os.Signal)
		signal.Notify(signals, syscall.SIGHUP, syscall.SIGTERM)

		// config reload signal
		reload := make(chan struct{})
		// go agent.Reload(reload)

		// wait for system exit and config reload signal
		go func() {
			select {
			case <-signals:
			case <-reload:
				isReload = true
			}
			close(shutdown)
		}()

		ag.Start(shutdown)
	}
}
