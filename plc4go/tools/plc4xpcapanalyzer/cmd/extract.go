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
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/extractor"
	"github.com/pkg/errors"
	"os"

	"github.com/spf13/cobra"
)

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract [protocolType] [pcapfile]",
	Short: "extract a pcap file using a driver supplied driver",
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
		extractor.Extract(pcapFile, protocolType)
		println("Done")
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)

	extractCmd.Flags().StringVarP(&config.ExtractConfigInstance.Client, "client", "c", "", "The client ip (this is useful for protocols where request/response is different e.g. modbus, cbus)")
	extractCmd.Flags().BoolVarP(&config.ExtractConfigInstance.ShowDirectionalIndicators, "show-directional-indicators", "", true, "Indicates if directional markers should be printed")
}
