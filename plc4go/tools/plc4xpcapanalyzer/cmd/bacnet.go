/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package cmd

import (
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/config"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/analyzer"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"os"

	"github.com/spf13/cobra"
)

// bacnetCmd represents the bacnet command
var bacnetCmd = &cobra.Command{
	Use:   "bacnet [pcapfile]",
	Short: "analyzes a pcap file using a bacnet driver",
	Long: `Analyzes a pcap file using a bacnet driver
TODO: document me
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires exactly one arguments")
		}
		pcapFile := args[0]
		if _, err := os.Stat(pcapFile); errors.Is(err, os.ErrNotExist) {
			return errors.Errorf("Pcap file not found %s", pcapFile)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		pcapFile := args[0]
		if !config.BacnetConfigInstance.NoFilter {
			if config.BacnetConfigInstance.Filter == "" && config.BacnetConfigInstance.BacnetFilter != "" {
				log.Debug().Str("filter", config.BacnetConfigInstance.Filter).Msg("Setting bacnet filter")
				config.BacnetConfigInstance.Filter = config.BacnetConfigInstance.BacnetFilter
			}
		} else {
			log.Info().Msg("All filtering disabled")
		}
		if err := analyzer.Analyze(pcapFile, "bacnet"); err != nil {
			panic(err)
		}
		println("Done")
	},
}

func init() {
	analyzeCmd.AddCommand(bacnetCmd)

	bacnetCmd.PersistentFlags().StringVarP(&config.BacnetConfigInstance.BacnetFilter, "default-bacnet-filter", "", "udp port 47808 and udp[4:2] > 29", "Defines the default filter when bacnet is selected")

	addAnalyzeFlags(bacnetCmd)
}
