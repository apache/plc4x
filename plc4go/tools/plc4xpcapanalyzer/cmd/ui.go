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
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/ui"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"os"

	"github.com/spf13/cobra"
)

// uiCmd represents the ui command
var uiCmd = &cobra.Command{
	Use:   "ui [pcapfile]",
	Short: "Start the ui with optional pcapfile",
	Long: `Analyzes a pcap file using a bacnet driver
TODO: document me
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return nil
		}
		pcapFile := args[0]
		if _, err := os.Stat(pcapFile); errors.Is(err, os.ErrNotExist) {
			return errors.Errorf("Pcap file not found %s", pcapFile)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ui.LoadConfig()
		application := ui.SetupApplication()
		ui.InitSubsystem()
		if len(args) > 0 {
			pcapFile := args[0]
			go func() {
				err := ui.OpenFile(pcapFile)
				if err != nil {
					log.Error().Err(err).Msg("Error opening argument file")
				}
			}()
		}

		defer ui.Shutdown()
		if err := application.Run(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
}
