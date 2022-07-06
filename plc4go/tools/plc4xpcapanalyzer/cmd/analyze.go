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
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/analyzer"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"math"
	"os"

	"github.com/spf13/cobra"
)

var validProtocolType = map[string]interface{}{
	"bacnet": nil,
	"c-bus":  nil,
}

// flags here
var (
	filter                              string
	noFilter, onlyParse, noBytesCompare bool
	client                              string
	bacnetFilter                        string
	startPackageNumber                  uint
	packageNumberLimit                  uint
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze [protocolType] [pcapfile]",
	Short: "analyzes a pcap file using a driver supplied driver",
	Long: `Analyzes a pcap file using a driver
TODO: document me
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires exactly two arguments")
		}
		if _, ok := validProtocolType[args[0]]; !ok {
			return errors.Errorf("Only following protocols are supported %v", validProtocolType)
		}
		pcapFile := args[1]
		if _, err := os.Stat(pcapFile); errors.Is(err, os.ErrNotExist) {
			return errors.Errorf("Pcap file not found %s", pcapFile)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		protocolType := args[0]
		pcapFile := args[1]
		if !noFilter {
			switch protocolType {
			case "bacnet":
				if filter != "" && bacnetFilter != "" {
					log.Debug().Str("filter", filter).Msg("Setting bacnet filter")
					filter = bacnetFilter
				}
			}
		} else {
			log.Info().Msg("All filtering disabled")
		}
		analyzer.Analyze(pcapFile, protocolType, filter, onlyParse, noBytesCompare, client, startPackageNumber, packageNumberLimit, verbosity)
		println("Done")
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().StringVarP(&filter, "filter", "f", "", "BFF filter to apply")
	analyzeCmd.Flags().BoolVarP(&noFilter, "no-filter", "n", false, "disable filter")
	analyzeCmd.Flags().BoolVarP(&onlyParse, "onlyParse", "o", false, "only parse messaged")
	analyzeCmd.Flags().BoolVarP(&noBytesCompare, "noBytesCompare", "b", false, "don't compare original bytes with serialized bytes")
	analyzeCmd.Flags().StringVarP(&client, "client", "c", "", "The client ip (this is useful for protocols where request/response is different e.g. modbus, cbus)")
	analyzeCmd.Flags().UintVarP(&startPackageNumber, "startPackageNumber", "s", 0, "Defines with what package number should be started")
	analyzeCmd.Flags().UintVarP(&packageNumberLimit, "packageNumberLimit", "l", math.MaxUint, "Defines how many packages should be parsed")
	// TODO: maybe it is smarter to convert this into subcommands because this option is only relevant to bacnet
	analyzeCmd.PersistentFlags().StringVarP(&bacnetFilter, "default-bacnet-filter", "", "udp port 47808 and udp[4:2] > 29", "Defines the default filter when bacnet is selected")
	// TODO: support other protocols
}
