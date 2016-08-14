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

	"github.com/corego/vgo/vgo/stream"
	"github.com/spf13/cobra"
)

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: " ",
	Run:   streamrun,
}

func init() {
	RootCmd.AddCommand(streamCmd)
}

func streamrun(cmd *cobra.Command, args []string) {
	s := stream.New()
	// start stream server
	go s.Start()
	// wait server stop signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("stream service received Signal: ", <-chSig)
	fmt.Println("stream service is going to stop")
	// stop stream server
	s.Close()
	fmt.Println("stream service is stopped")
}
