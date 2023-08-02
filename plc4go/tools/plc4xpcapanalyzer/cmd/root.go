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
	"bytes"
	"fmt"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "plc4xpcapanalyzer",
	Short: "analyzes pcaps with plc4x",
	Long:  `todo describe me`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&config.RootConfigInstance.CfgFile, "config", "", "config file (default is $HOME/.plc4xpcapanalyzer.yaml)")
	rootCmd.PersistentFlags().StringVar(&config.RootConfigInstance.LogType, "log-type", "text", "define how the log will be evaluated")
	rootCmd.PersistentFlags().StringVar(&config.RootConfigInstance.LogLevel, "log-level", "error", "define the log Level")
	rootCmd.PersistentFlags().CountVarP(&config.RootConfigInstance.Verbosity, "verbose", "v", "counted verbosity")
	rootCmd.PersistentFlags().BoolVarP(&config.RootConfigInstance.HideProgressBar, "hide-progress-bar", "", false, "hides the progress bar")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if config.RootConfigInstance.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(config.RootConfigInstance.CfgFile)
	} else {
		// Find user config directory.
		home, err := os.UserConfigDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("plc4xpcapanalyzer-viper")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	if config.RootConfigInstance.LogType == "text" {
		log.Logger = log.
			//// Enable below if you want to see the filenames
			//With().Caller().Logger().
			Output(zerolog.NewConsoleWriter(
				func(w *zerolog.ConsoleWriter) {
					w.Out = os.Stderr
				},
				func(w *zerolog.ConsoleWriter) {
					w.FormatFieldValue = func(i interface{}) string {
						if aString, ok := i.(string); ok && strings.Contains(aString, "\\n") {
							return fmt.Sprintf("\x1b[%dm%v\x1b[0m", 31, "see below")
						}
						return fmt.Sprintf("%s", i)
					}
					w.FormatExtra = func(m map[string]interface{}, buffer *bytes.Buffer) error {
						for key, i := range m {
							if aString, ok := i.(string); ok && strings.Contains(aString, "\n") {
								buffer.WriteString("\n")
								buffer.WriteString(fmt.Sprintf("\x1b[%dm%v\x1b[0m", 32, "field "+key))
								buffer.WriteString(":\n" + aString)
							}
						}
						return nil
					}
				},
			),
			).
			Level(parseLogLevel())
	}
}

func parseLogLevel() zerolog.Level {
	level, err := zerolog.ParseLevel(config.RootConfigInstance.LogLevel)
	if err != nil {
		log.Fatal().Err(err).Str("level", config.RootConfigInstance.LogLevel).Msg("Unknown log level")
	}
	return level
}
