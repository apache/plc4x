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

package ui

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/spi"
	plc4x_config "github.com/apache/plc4x/plc4go/pkg/api/config"
	cliConfig "github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/config"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/analyzer"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/extractor"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sruehl/tview"
	"os"
	"path"
	"strings"
	"time"
)

const rootCommandIndicator = "rootCommand"

var rootCommand = Command{
	Name: rootCommandIndicator,
	subCommands: []Command{
		{
			Name:        "ls",
			Description: "list directories",
			action: func(_ Command, dir string) error {
				if dir == "" {
					dir = currentDir
				}
				_, _ = fmt.Fprintf(commandOutput, "dir cotents of %s\n", dir)
				readDir, err := os.ReadDir(dir)
				if err != nil {
					return err
				}
				for _, dirEntry := range readDir {
					isDir := dirEntry.IsDir()
					name := dirEntry.Name()
					name = strings.TrimPrefix(name, dir)
					if isDir {
						name = fmt.Sprintf("[#0000ff]%s[white]", name)
					} else if strings.HasSuffix(name, ".pcap") || strings.HasSuffix(name, ".pcapng") {
						name = fmt.Sprintf("[#00ff00]%s[white]", name)
					}
					_, _ = fmt.Fprintf(commandOutput, "%s\n", name)
				}
				return nil
			},
			// TODO: add parameter suggestions
		},
		{
			Name:        "cd",
			Description: "changes directory",
			action: func(_ Command, newDir string) error {
				var proposedCurrentDir string
				if newDir == "" {
					var err error
					proposedCurrentDir, err = os.UserHomeDir()
					if err != nil {
						return err
					}
				} else if strings.HasPrefix(newDir, "."+string(os.PathSeparator)) {
					proposedCurrentDir = currentDir + strings.TrimPrefix(newDir, ".")
				} else if strings.HasPrefix(newDir, ""+string(os.PathSeparator)) {
					proposedCurrentDir = newDir
				}
				stat, err := os.Stat(proposedCurrentDir)
				if err != nil {
					return err
				}
				if !stat.IsDir() {
					return errors.Errorf("%s is not a dir", newDir)
				}
				currentDir = proposedCurrentDir
				_, _ = fmt.Fprintf(commandOutput, "current directory: %s\n", currentDir)
				return nil
			},
			parameterSuggestions: func(currentText string) (entries []string) {
				if strings.HasPrefix(currentText, string(os.PathSeparator)) {
					dirEntries, err := os.ReadDir(currentText)
					if err != nil {
						plc4xpcapanalyzerLog.Warn().Err(err).Msg("Error suggesting directories")
						return
					}
					for _, dirEntry := range dirEntries {
						entry := path.Join(currentText, dirEntry.Name())
						entries = append(entries, entry)
					}
				} else {
					dirEntries, err := os.ReadDir(currentDir)
					if err != nil {
						plc4xpcapanalyzerLog.Warn().Err(err).Msg("Error suggesting directories")
						return
					}
					for _, dirEntry := range dirEntries {
						entry := path.Join(".", dirEntry.Name())
						entries = append(entries, entry)
					}
				}
				return
			},
		},
		{
			Name:        "pwd",
			Description: "shows current directory",
			action: func(_ Command, _ string) error {
				_, _ = fmt.Fprintf(commandOutput, "current directory: %s\n", currentDir)
				return nil
			},
		},
		{
			Name:        "open",
			Description: "open file",
			action: func(_ Command, pcapFile string) error {
				return OpenFile(pcapFile)
			},
			parameterSuggestions: func(currentText string) (entries []string) {
				entries = append(entries, config.History.Last10Files...)
				readDir, err := os.ReadDir(currentDir)
				if err != nil {
					return
				}
				for _, dirEntry := range readDir {
					name := dirEntry.Name()
					name = strings.TrimPrefix(name, currentDir)
					if strings.HasSuffix(dirEntry.Name(), ".pcap") || strings.HasSuffix(name, ".pcapng") {
						entries = append(entries, name)
					}
				}
				return
			},
		},
		{
			Name:        "analyze",
			Description: "Analyzes a pcap file using a driver",
			action: func(_ Command, protocolTypeAndPcapFile string) error {
				split := strings.Split(protocolTypeAndPcapFile, " ")
				if len(split) != 2 {
					return errors.Errorf("expect protocol and pcapfile")
				}
				protocolType := split[0]
				pcapFile := strings.TrimPrefix(protocolTypeAndPcapFile, protocolType+" ")
				cliConfig.PcapConfigInstance.Client = config.HostIp
				cliConfig.RootConfigInstance.HideProgressBar = true
				cliConfig.RootConfigInstance.Verbosity = 4
				return analyzer.AnalyzeWithOutputAndCallback(pcapFile, protocolType, tview.ANSIWriter(messageOutput), tview.ANSIWriter(messageOutput), func(parsed spi.Message) {
					spiNumberOfMessagesReceived++
					spiMessageReceived(spiNumberOfMessagesReceived, time.Now(), parsed)
				})
			},
			parameterSuggestions: func(currentText string) (entries []string) {
				for _, file := range loadedPcapFiles {
					for _, protocol := range protocolList {
						entries = append(entries, protocol+" "+file.path)
					}
				}
				return
			},
		},
		{
			Name:        "extract",
			Description: "Extract a pcap file using a driver",
			action: func(_ Command, protocolTypeAndPcapFile string) error {
				split := strings.Split(protocolTypeAndPcapFile, " ")
				if len(split) != 2 {
					return errors.Errorf("expect protocol and pcapfile")
				}
				protocolType := split[0]
				pcapFile := strings.TrimPrefix(protocolTypeAndPcapFile, protocolType+" ")
				cliConfig.PcapConfigInstance.Client = config.HostIp
				cliConfig.RootConfigInstance.HideProgressBar = true
				cliConfig.RootConfigInstance.Verbosity = 4
				return extractor.ExtractWithOutput(pcapFile, protocolType, tview.ANSIWriter(messageOutput), tview.ANSIWriter(messageOutput))
			},
			parameterSuggestions: func(currentText string) (entries []string) {
				for _, file := range loadedPcapFiles {
					for _, protocol := range protocolList {
						entries = append(entries, protocol+" "+file.path)
					}
				}
				return
			},
		},
		{
			Name:        "host",
			Description: "The host which is assumed to be the sender (important for protocols that are directional)",
			subCommands: []Command{
				{
					Name: "set",
					action: func(_ Command, host string) error {
						config.HostIp = host
						return nil
					},
				},
				{
					Name: "get",
					action: func(_ Command, host string) error {
						_, _ = fmt.Fprintf(consoleOutput, "current set host %s", config.HostIp)
						return nil
					},
				},
			},
		},
		{
			Name:        "register",
			Description: "register a driver in the subsystem",
			action: func(_ Command, driver string) error {
				return registerDriver(driver)
			},
			parameterSuggestions: func(currentText string) (entries []string) {
				for _, protocol := range protocolList {
					if strings.HasPrefix(protocol, currentText) {
						entries = append(entries, protocol)
					}
				}
				return
			},
		},
		{
			Name:        "quit",
			Description: "Quits the application",
		},
		{
			Name:        "log",
			Description: "Log related operations",
			subCommands: []Command{
				{
					Name:        "get",
					Description: "Get a log level",
					action: func(_ Command, _ string) error {
						_, _ = fmt.Fprintf(commandOutput, "Current log level %s", log.Logger.GetLevel())
						return nil
					},
				},
				{
					Name:        "set",
					Description: "Sets a log level",
					action: func(_ Command, level string) error {
						parseLevel, err := zerolog.ParseLevel(level)
						if err != nil {
							return errors.Wrapf(err, "Error setting log level")
						}
						setLevel(parseLevel)
						log.Logger = log.Logger.Level(parseLevel)
						return nil
					},
					parameterSuggestions: func(currentText string) (entries []string) {
						levels := []string{
							zerolog.LevelTraceValue,
							zerolog.LevelDebugValue,
							zerolog.LevelInfoValue,
							zerolog.LevelWarnValue,
							zerolog.LevelErrorValue,
							zerolog.LevelFatalValue,
							zerolog.LevelPanicValue,
						}
						for _, level := range levels {
							entries = append(entries, level)
						}
						return
					},
				},
			},
		},
		{
			Name:        "plc4x-conf",
			Description: "plc4x related settings",
			subCommands: []Command{
				{
					Name:        "TraceTransactionManagerWorkers",
					Description: "print information about transaction manager workers",
					subCommands: []Command{
						{
							Name:        "on",
							Description: "trace on",
							action: func(_ Command, _ string) error {
								plc4x_config.TraceTransactionManagerWorkers = true
								return nil
							},
						},
						{
							Name:        "off",
							Description: "trace off",
							action: func(_ Command, _ string) error {
								plc4x_config.TraceTransactionManagerWorkers = false
								return nil
							},
						},
					},
				},
				{
					Name:        "TraceTransactionManagerTransactions",
					Description: "print information about transaction manager transactions",
					subCommands: []Command{
						{
							Name:        "on",
							Description: "trace on",
							action: func(_ Command, _ string) error {
								plc4x_config.TraceTransactionManagerTransactions = true
								return nil
							},
						},
						{
							Name:        "off",
							Description: "trace off",
							action: func(_ Command, _ string) error {
								plc4x_config.TraceTransactionManagerTransactions = false
								return nil
							},
						},
					},
				},
				{
					Name:        "TraceDefaultMessageCodecWorker",
					Description: "print information about message codec workers",
					subCommands: []Command{
						{
							Name:        "on",
							Description: "trace on",
							action: func(_ Command, _ string) error {
								plc4x_config.TraceDefaultMessageCodecWorker = true
								return nil
							},
						},
						{
							Name:        "off",
							Description: "trace off",
							action: func(_ Command, _ string) error {
								plc4x_config.TraceDefaultMessageCodecWorker = false
								return nil
							},
						},
					},
				},
				{
					Name:        "plc4xpcapanalyzer-debug",
					Description: "Prints out debug information of the pcap analyzer itself",
					subCommands: []Command{
						{
							Name:        "on",
							Description: "debug on",
							action: func(_ Command, _ string) error {
								plc4xpcapanalyzerLog = zerolog.New(zerolog.ConsoleWriter{Out: tview.ANSIWriter(consoleOutput)})
								return nil
							},
						},
						{
							Name:        "off",
							Description: "debug off",
							action: func(_ Command, _ string) error {
								plc4xpcapanalyzerLog = zerolog.Nop()
								return nil
							},
						},
					},
				},
				{
					Name:        "auto-register",
					Description: "autoregister driver at startup",
					subCommands: []Command{
						{
							Name: "list",
							action: func(currentCommand Command, argument string) error {
								_, _ = fmt.Fprintf(commandOutput, "Auto-register enabled drivers:\n  %s\n", strings.Join(config.AutoRegisterDrivers, "\n  "))
								return nil
							},
						},
						{
							Name: "enable",
							action: func(_ Command, argument string) error {
								return enableAutoRegister(argument)
							},
							parameterSuggestions: func(currentText string) (entries []string) {
								for _, protocol := range protocolList {
									if strings.HasPrefix(protocol, currentText) {
										entries = append(entries, protocol)
									}
								}
								return
							},
						},
						{
							Name: "disable",
							action: func(_ Command, argument string) error {
								return disableAutoRegister(argument)
							},
							parameterSuggestions: func(currentText string) (entries []string) {
								for _, protocol := range protocolList {
									if strings.HasPrefix(protocol, currentText) {
										entries = append(entries, protocol)
									}
								}
								return
							},
						},
					},
				},
			},
		},
		{
			Name:        "history",
			Description: "outputs the last commands",
			action: func(_ Command, _ string) error {
				outputCommandHistory()
				return nil
			},
		},
		{
			Name:        "clear",
			Description: "clear all outputs",
			action: func(_ Command, _ string) error {
				messageOutputClear()
				consoleOutputClear()
				commandOutputClear()
				return nil
			},
			subCommands: []Command{
				{
					Name:        "message",
					Description: "clears message output",
					action: func(_ Command, _ string) error {
						messageOutputClear()
						return nil
					},
				},
				{
					Name:        "console",
					Description: "clears console output",
					action: func(_ Command, _ string) error {
						consoleOutputClear()
						return nil
					},
				},
				{
					Name:        "command",
					Description: "clears command output",
					action: func(_ Command, _ string) error {
						commandOutputClear()
						return nil
					},
				},
			},
		},
	},
}

func init() {
	// Because of the cycle we need to define the help command here as it needs access to the to command
	rootCommand.subCommands = append(rootCommand.subCommands, Command{
		Name:        "help",
		Description: "prints out this help",
		action: func(_ Command, _ string) error {
			_, _ = fmt.Fprintf(commandOutput, "[#0000ff]Available commands[white]\n")
			rootCommand.visit(0, func(currentIndent int, command Command) {
				indentString := strings.Repeat("  ", currentIndent)
				description := command.Description
				if description == "" {
					description = command.Name + "s"
				}
				_, _ = fmt.Fprintf(commandOutput, "%s [#00ff00]%s[white]: %s\n", indentString, command.Name, description)
			})
			return nil
		},
	})
}

var NotDirectlyExecutable = errors.New("Not directly executable")

type Command struct {
	Name                 string
	Description          string
	action               func(currentCommand Command, argument string) error
	subCommands          []Command
	parameterSuggestions func(currentText string) (entries []string)
}

func (c Command) Completions(currentCommandText string) (entries []string) {
	if c.Name == rootCommandIndicator && len(currentCommandText) == 0 {
		// We don't return anything here to not pollute the command text by default
		return
	}
	if c.acceptsCurrentText(currentCommandText) {
		currentCommandPrefix := c.currentCommandPrefix()
		doesCommandTextTargetSubCommand := c.doesCommandTextTargetSubCommand(currentCommandPrefix)
		if c.hasDirectExecution() && !doesCommandTextTargetSubCommand {
			if c.parameterSuggestions != nil {
				preparedForParameters := c.prepareForParameters(currentCommandText)
				for _, parameterSuggestion := range c.parameterSuggestions(preparedForParameters) {
					entries = append(entries, currentCommandPrefix+parameterSuggestion)
				}
			} else if currentCommandText == "" {
				entries = append(entries, c.Name)
			}
		}
		if doesCommandTextTargetSubCommand {
			remainder := c.prepareForSubCommand(currentCommandText)
			for _, command := range c.subCommands {
				for _, subCommandCompletions := range command.Completions(remainder) {
					entries = append(entries, currentCommandPrefix+subCommandCompletions)
				}
			}
		}
	} else if strings.HasPrefix(c.Name, currentCommandText) {
		// Suggest ourselves if we start with the current letter
		entries = append(entries, c.Name)
	}
	return
}

func (c Command) acceptsCurrentText(currentCommandText string) bool {
	if c.Name == rootCommandIndicator {
		return true
	}
	hasThePrefix := strings.HasPrefix(currentCommandText, c.Name)
	hasNoMatchingAlternative := !strings.HasPrefix(currentCommandText, c.Name+"-")
	accepts := hasThePrefix && hasNoMatchingAlternative
	plc4xpcapanalyzerLog.Debug().Msgf("%s accepts %t", c, accepts)
	return accepts
}

func (c Command) doesCommandTextTargetSubCommand(currentCommandText string) bool {
	if c.Name == rootCommandIndicator {
		return true
	}
	if len(c.subCommands) == 0 {
		return false
	}
	return strings.HasPrefix(currentCommandText, c.currentCommandPrefix())
}

func (c Command) prepareForParameters(currentCommandText string) string {
	if currentCommandText == c.Name {
		return ""
	}
	return strings.TrimPrefix(currentCommandText, c.currentCommandPrefix())
}
func (c Command) prepareForSubCommand(currentCommandText string) string {
	return strings.TrimPrefix(currentCommandText, c.currentCommandPrefix())
}

func (c Command) currentCommandPrefix() string {
	if c.Name == rootCommandIndicator {
		return ""
	}
	return c.Name + " "
}

func (c Command) hasDirectExecution() bool {
	return c.action != nil
}

func Execute(commandText string) error {
	err := rootCommand.Execute(commandText)
	if err == nil {
		addCommandHistoryEntry(commandText)
	}
	return err
}

func (c Command) Execute(commandText string) error {
	plc4xpcapanalyzerLog.Debug().Msgf("%s executes %s", c, commandText)
	if !c.acceptsCurrentText(commandText) {
		return errors.Errorf("%s doesn't understand %s", c.Name, commandText)
	}
	if c.doesCommandTextTargetSubCommand(commandText) {
		prepareForSubCommandForSubCommand := c.prepareForSubCommand(commandText)
		for _, command := range c.subCommands {
			if command.acceptsCurrentText(prepareForSubCommandForSubCommand) {
				plc4xpcapanalyzerLog.Debug().Msgf("%s delegates to sub %s", c, command)
				return command.Execute(prepareForSubCommandForSubCommand)
			}
		}
		return errors.Errorf("%s not accepted by any subcommands of %s", commandText, c.Name)
	} else {
		if c.action == nil {
			return NotDirectlyExecutable
		}
		plc4xpcapanalyzerLog.Debug().Msgf("%s executes %s directly", c, commandText)
		preparedForParameters := c.prepareForParameters(commandText)
		return c.action(c, preparedForParameters)
	}
}

func (c Command) visit(i int, f func(currentIndent int, command Command)) {
	f(i, c)
	for _, subCommand := range c.subCommands {
		subCommand.visit(i+1, f)
	}
}

func (c Command) String() string {
	return c.Name
}
