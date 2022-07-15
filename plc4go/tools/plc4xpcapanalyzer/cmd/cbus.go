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

// cbusCmd represents the cbus command
var cbusCmd = &cobra.Command{
	Use:   "c-bus [pcapfile]",
	Short: "analyzes a pcap file using a c-bus driver",
	Long: `Analyzes a pcap file using a c-bus driver
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
		if !config.CBusConfigInstance.NoFilter {
			if config.CBusConfigInstance.Filter == "" && config.CBusConfigInstance.CBusFilter != "" {
				log.Debug().Str("filter", config.CBusConfigInstance.Filter).Msg("Setting cbus filter")
				config.CBusConfigInstance.Filter = config.CBusConfigInstance.CBusFilter
			}
		} else {
			log.Info().Msg("All filtering disabled")
		}
		analyzer.Analyze(pcapFile, "c-bus")
		println("Done")
	},
}

func init() {
	analyzeCmd.AddCommand(cbusCmd)

	cbusCmd.PersistentFlags().StringVarP(&config.CBusConfigInstance.CBusFilter, "default-cbus-filter", "", "tcp port 10001", "Defines the default filter when c-bus is selected")

	cbusCmd.Flags().BoolVarP(&config.CBusConfigInstance.Connect, "cbus-connect", "", false, "Defines that SAL messages can occur at any time")
	cbusCmd.Flags().BoolVarP(&config.CBusConfigInstance.Smart, "cbus-smart", "", false, "Disable echo of characters. When used with connect SAL have a long option. Select long from of most CAL replies")
	cbusCmd.Flags().BoolVarP(&config.CBusConfigInstance.Idmon, "cbus-idmon", "", false, "only works with smart. Select long form of CAL messages")
	cbusCmd.Flags().BoolVarP(&config.CBusConfigInstance.Exstat, "cbus-exstat", "", false, "useful with smart. Select long form, extended format for all monitored and initiated status requests")
	cbusCmd.Flags().BoolVarP(&config.CBusConfigInstance.Monitor, "cbus-monitor", "", false, "monitors all traffic for status requests. Status requests will be returned as CAL. Replies are modified by exstat. Usually used in conjunction with connect.")
	cbusCmd.Flags().BoolVarP(&config.CBusConfigInstance.Monall, "cbus-monall", "", false, "Same as connect. In addition it will return remote network SAL")
	cbusCmd.Flags().BoolVarP(&config.CBusConfigInstance.Pun, "cbus-pun", "", false, "Serial interface will emit a power up notification")
	cbusCmd.Flags().BoolVarP(&config.CBusConfigInstance.Pcn, "cbus-pcn", "", false, "causes parameter change notifications to be emitted.")
	cbusCmd.Flags().BoolVarP(&config.CBusConfigInstance.Srchk, "cbus-srchk", "", false, "enabled the crc checks")

	addAnalyzeFlags(cbusCmd)
}
