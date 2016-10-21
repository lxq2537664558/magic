/***
The code in this project source from https://github.com/influxdata/telegraf
and contains modifications so that no other dependency from that project is needed. Other modifications included
1.reconstructing the whole project structure
2. removing uneccessary code and dependencies
3. Many code and performance optimizing
4. delete many useless inputs and outputs ,eg: java_metrics_http
5. optimize some inputs,eg.: log_collect
6. optimize configurations way, eg : removing many useless structures

Itis licensed under Apache Version 2.0, http://www.apache.org/licenses/LICENSE-2.0.html
***/

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

package main

import "github.com/aiyun/openapm/mecury/cmd"

func main() {
	cmd.Execute()
}
