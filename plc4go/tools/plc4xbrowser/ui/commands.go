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
	plc4x_config "github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sruehl/tview"
	"net/url"
	"strings"
	"time"
)

const rootCommandIndicator = "rootCommand"

var commands = map[inputMode]Command{
	normalMode:        rootCommand,
	readEditMode:      rootCommand,
	writeEditMode:     rootCommand,
	subscribeEditMode: rootCommand,
}

var rootCommand = Command{
	Name: rootCommandIndicator,
	subCommands: []Command{
		{
			Name:        "discover",
			Description: "Discovers devices",
			action: func(_ Command, driverId string) error {
				if driver, ok := registeredDrivers[driverId]; ok {
					if !driver.SupportsDiscovery() {
						return errors.Errorf("%s doesn't support discovery", driverId)
					}
					return driver.Discover(func(event model.PlcDiscoveryItem) {
						_, _ = fmt.Fprintf(messageOutput, "%v\n", event)
					})
				} else {
					return errors.Errorf("%s not registered", driverId)
				}
			},
			parameterSuggestions: func(currentText string) (entries []string) {
				for _, protocol := range protocolList {
					entries = append(entries, protocol)
				}
				return
			},
		},
		{
			Name:        "connect",
			Description: "Connects to a device",
			action: func(_ Command, connectionString string) error {
				log.Info().Msgf("connect connectionString [%s]", connectionString)
				connectionUrl, err := url.Parse(connectionString)
				if err != nil {
					return errors.Wrapf(err, "can't parse connection url %s", connectionString)
				}
				addHostHistoryEntry(connectionUrl.Host)
				connectionId := fmt.Sprintf("%s://%s", connectionUrl.Scheme, connectionUrl.Host)
				if _, ok := connections[connectionId]; ok {
					return errors.Errorf("%s already connected", connectionId)
				}
				connectionResult := <-driverManager.GetConnection(connectionString)
				if err := connectionResult.GetErr(); err != nil {
					return errors.Wrapf(err, "%s can't connect to", connectionUrl.Host)
				}
				log.Info().Msgf("%s connected", connectionId)
				connections[connectionId] = connectionResult.GetConnection()
				connectionsChanged()
				return nil
			},
			parameterSuggestions: func(currentText string) (entries []string) {
				for _, protocol := range protocolList {
					if strings.HasPrefix(currentText, protocol) {
						for _, host := range config.History.Last10Hosts {
							entries = append(entries, protocol+"://"+host)
						}
						entries = append(entries, currentText)
					} else {
						entries = append(entries, protocol)
					}
				}
				return
			},
		},
		{
			Name:        "disconnect",
			Description: "Disconnect a connection",
			action: func(_ Command, connectionString string) error {
				if connection, ok := connections[connectionString]; !ok {
					return errors.Errorf("%s not connected", connectionString)
				} else {
					closeResult := <-connection.Close()
					log.Info().Msgf("%s disconnected", connectionString)
					delete(connections, connectionString)
					connectionsChanged()
					if err := closeResult.GetErr(); err != nil {
						return errors.Wrapf(err, "%s can't close", connectionString)
					}
				}
				return nil
			},
			parameterSuggestions: func(currentText string) (entries []string) {
				for connectionsString, _ := range connections {
					entries = append(entries, connectionsString)
				}
				return
			},
		},
		{
			Name:        "read",
			Description: "Starts a read request (switched mode to read edit)",
			action: func(_ Command, connectionsString string) error {
				if connection, ok := connections[connectionsString]; !ok {
					return errors.Errorf("%s not connected", connectionsString)
				} else {
					return errors.Errorf("%s mode switch not yet implemented", connection)
				}
			},
			parameterSuggestions: func(currentText string) (entries []string) {
				for connectionsString, _ := range connections {
					entries = append(entries, connectionsString)
				}
				return
			},
		},
		{
			Name:        "read-direct",
			Description: "Builds a read request with the supplied field",
			action: func(c Command, connectionsStringAndFieldQuery string) error {
				split := strings.Split(connectionsStringAndFieldQuery, " ")
				if len(split) != 2 {
					return errors.Errorf("%s expects exactly two arguments [connection url] [fieldQuery]", c)
				}
				connectionsString := split[0]
				if connection, ok := connections[connectionsString]; !ok {
					return errors.Errorf("%s not connected", connectionsString)
				} else {
					start := time.Now()
					readRequest, err := connection.ReadRequestBuilder().
						AddQuery("readField", split[1]).
						Build()
					if err != nil {
						return errors.Wrapf(err, "%s can't read", connectionsString)
					}
					readRequestResult := <-readRequest.Execute()
					if err := readRequestResult.GetErr(); err != nil {
						return errors.Wrapf(err, "%s can't read", connectionsString)
					}
					plc4xBrowserLog.Debug().Msgf("read took %f seconds", time.Now().Sub(start).Seconds())
					if err := readRequestResult.GetErr(); err != nil {
						return errors.Wrapf(err, "%s error reading", connectionsString)
					}
					numberOfMessagesReceived++
					messageReceived(numberOfMessagesReceived, time.Now(), readRequestResult.GetResponse())
				}
				return nil
			},
			parameterSuggestions: func(currentText string) (entries []string) {
				for connectionsString, _ := range connections {
					if strings.HasPrefix(currentText, connectionsString+"") {
						parse, _ := url.Parse(connectionsString)
						switch parse.Scheme {
						// TODO: add to protocol suggestor so it can be reused.

						}
					} else {
						entries = append(entries, connectionsString)
					}
				}
				return
			},
		},
		{
			Name:        "write",
			Description: "Starts a write request (switched mode to write edit)",
			action: func(_ Command, connectionsString string) error {
				if connection, ok := connections[connectionsString]; !ok {
					return errors.Errorf("%s not connected", connectionsString)
				} else {
					return errors.Errorf("%s mode switch not yet implemented", connection)
				}
			},
			parameterSuggestions: func(currentText string) (entries []string) {
				for connectionsString, _ := range connections {
					entries = append(entries, connectionsString)
				}
				return
			},
		},
		{
			Name:        "write-direct",
			Description: "Builds a write request with the supplied field",
			action: func(c Command, connectionsStringAndFieldQuery string) error {
				split := strings.Split(connectionsStringAndFieldQuery, " ")
				if len(split) != 3 {
					return errors.Errorf("%s expects exactly three arguments [connection url] [fieldQuery] [value]", c)
				}
				connectionsString := split[0]
				if connection, ok := connections[connectionsString]; !ok {
					return errors.Errorf("%s not connected", connectionsString)
				} else {
					start := time.Now()
					writeRequest, err := connection.WriteRequestBuilder().
						AddQuery("writeField", split[1], split[2]).
						Build()
					if err != nil {
						return errors.Wrapf(err, "%s can't write", connectionsString)
					}
					writeRequestResult := <-writeRequest.Execute()
					if err := writeRequestResult.GetErr(); err != nil {
						return errors.Wrapf(err, "%s can't write", connectionsString)
					}
					plc4xBrowserLog.Debug().Msgf("write took %f seconds", time.Now().Sub(start).Seconds())
					if err := writeRequestResult.GetErr(); err != nil {
						return errors.Wrapf(err, "%s error writing", connectionsString)
					}
					numberOfMessagesReceived++
					messageReceived(numberOfMessagesReceived, time.Now(), writeRequestResult.GetResponse())
				}
				return nil
			},
			parameterSuggestions: func(currentText string) (entries []string) {
				for connectionsString, _ := range connections {
					if strings.HasPrefix(currentText, connectionsString+"") {
						parse, _ := url.Parse(connectionsString)
						switch parse.Scheme {
						// TODO: add to protocol suggestor so it can be reused.

						}
					} else {
						entries = append(entries, connectionsString)
					}
				}
				return
			},
		},
		{
			Name:        "browse",
			Description: "Starts a browse request (switched mode to browse edit)",
			action: func(_ Command, connectionsString string) error {
				if connection, ok := connections[connectionsString]; !ok {
					return errors.Errorf("%s not connected", connectionsString)
				} else {
					return errors.Errorf("%s mode switch not yet implemented", connection)
				}
			},
			parameterSuggestions: func(currentText string) (entries []string) {
				for connectionsString, _ := range connections {
					entries = append(entries, connectionsString)
				}
				return
			},
		},
		{
			Name:        "browse-direct",
			Description: "Builds a browse request with the supplied field",
			action: func(c Command, connectionsStringAndFieldQuery string) error {
				split := strings.Split(connectionsStringAndFieldQuery, " ")
				if len(split) != 2 {
					return errors.Errorf("%s expects exactly three arguments [connection url] [fieldQuery]", c)
				}
				connectionsString := split[0]
				if connection, ok := connections[connectionsString]; !ok {
					return errors.Errorf("%s not connected", connectionsString)
				} else {
					start := time.Now()
					browseRequest, err := connection.BrowseRequestBuilder().
						AddQuery("browseField", split[1]).
						Build()
					if err != nil {
						return errors.Wrapf(err, "%s can't browse", connectionsString)
					}
					browseRequestResult := <-browseRequest.ExecuteWithInterceptor(func(result model.PlcBrowseEvent) bool {
						numberOfMessagesReceived++
						messageReceived(numberOfMessagesReceived, time.Now(), result)
						return true
					})
					if err := browseRequestResult.GetErr(); err != nil {
						return errors.Wrapf(err, "%s can't browse", connectionsString)
					}
					plc4xBrowserLog.Debug().Msgf("write took %f seconds", time.Now().Sub(start).Seconds())
					if err := browseRequestResult.GetErr(); err != nil {
						return errors.Wrapf(err, "%s error browse", connectionsString)
					}
					numberOfMessagesReceived++
					messageReceived(numberOfMessagesReceived, time.Now(), browseRequestResult.GetResponse())
				}
				return nil
			},
			parameterSuggestions: func(currentText string) (entries []string) {
				for connectionsString, _ := range connections {
					if strings.HasPrefix(currentText, connectionsString+"") {
						parse, _ := url.Parse(connectionsString)
						switch parse.Scheme {
						// TODO: add to protocol suggestor so it can be reused.
						case "c-bus":
							entries = append(entries, connectionsString+" info/*/*")
						}
					} else {
						entries = append(entries, connectionsString)
					}
				}
				return
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
			Name:        "subscribe",
			Description: "Starts a subscription request (switched mode to subscribe edit)",
			action: func(_ Command, connectionsString string) error {
				if connection, ok := connections[connectionsString]; !ok {
					return errors.Errorf("%s not connected", connectionsString)
				} else {
					return errors.Errorf("%s mode switch not yet implemented", connection)
				}
			},
			parameterSuggestions: func(currentText string) (entries []string) {
				for connectionsString, _ := range connections {
					entries = append(entries, connectionsString)
				}
				return
			},
		},
		{
			Name:        "subscribe-direct",
			Description: "Builds a subscriptions request with the supplied field",
			action: func(c Command, connectionsStringAndFieldQuery string) error {
				split := strings.Split(connectionsStringAndFieldQuery, " ")
				if len(split) != 2 {
					return errors.Errorf("%s expects exactly two arguments [connection url] [fieldQuery]", c)
				}
				connectionsString := split[0]
				if connection, ok := connections[connectionsString]; !ok {
					return errors.Errorf("%s not connected", connectionsString)
				} else {
					subscriptionRequest, err := connection.SubscriptionRequestBuilder().
						AddEventQuery("subscriptionField", split[1]).
						AddPreRegisteredConsumer("subscriptionField", func(event model.PlcSubscriptionEvent) {
							numberOfMessagesReceived++
							messageReceived(numberOfMessagesReceived, time.Now(), event)
						}).
						Build()
					if err != nil {
						return errors.Wrapf(err, "%s can't subscribe", connectionsString)
					}
					subscriptionRequestResult := <-subscriptionRequest.Execute()
					if err := subscriptionRequestResult.GetErr(); err != nil {
						return errors.Wrapf(err, "%s can't subscribe", connectionsString)
					}
					log.Info().Msgf("subscription result\n%s", subscriptionRequestResult.GetResponse())
				}
				return nil
			},
			parameterSuggestions: func(currentText string) (entries []string) {
				for connectionsString, _ := range connections {
					entries = append(entries, connectionsString)
					if strings.HasPrefix(currentText, connectionsString) {
						parse, _ := url.Parse(connectionsString)
						switch parse.Scheme {
						// TODO: add to protocol suggestor so it can be reused.
						case "c-bus":
							entries = append(entries, connectionsString+" salmonitor/*/*")
							entries = append(entries, connectionsString+" mmimonitor/*/*")
						}
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
					Name:        "plc4xbrowser-debug",
					Description: "Prints out debug information of the browser itself",
					subCommands: []Command{
						{
							Name:        "on",
							Description: "debug on",
							action: func(_ Command, _ string) error {
								plc4xBrowserLog = zerolog.New(zerolog.ConsoleWriter{Out: tview.ANSIWriter(consoleOutput)})
								return nil
							},
						},
						{
							Name:        "off",
							Description: "debug off",
							action: func(_ Command, _ string) error {
								plc4xBrowserLog = zerolog.Nop()
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
			} else {
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
	plc4xBrowserLog.Debug().Msgf("%s accepts %t", c, accepts)
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
	plc4xBrowserLog.Debug().Msgf("%s executes %s", c, commandText)
	if !c.acceptsCurrentText(commandText) {
		return errors.Errorf("%s doesn't understand %s", c.Name, commandText)
	}
	if c.doesCommandTextTargetSubCommand(commandText) {
		prepareForSubCommandForSubCommand := c.prepareForSubCommand(commandText)
		for _, command := range c.subCommands {
			if command.acceptsCurrentText(prepareForSubCommandForSubCommand) {
				plc4xBrowserLog.Debug().Msgf("%s delegates to sub %s", c, command)
				return command.Execute(prepareForSubCommandForSubCommand)
			}
		}
		return errors.Errorf("%s not accepted by any subcommands of %s", commandText, c.Name)
	} else {
		if c.action == nil {
			return NotDirectlyExecutable
		}
		plc4xBrowserLog.Debug().Msgf("%s executes %s directly", c, commandText)
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
