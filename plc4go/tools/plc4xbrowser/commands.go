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

package main

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/ads"
	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	"github.com/apache/plc4x/plc4go/internal/cbus"
	"github.com/apache/plc4x/plc4go/internal/s7"
	plc4x_config "github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/transports"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/url"
	"strings"
)

var plc4xBrowserLog = zerolog.Nop()

const rootCommandIndicator = "rootCommand"

const protocols = "ads,bacnetip,c-bus,s7"

var protocolsSuggestions = strings.Split(protocols, ",")

var commands = map[inputMode]Command{
	normalMode:        rootCommand,
	subscribeEditMode: rootCommand,
}

var rootCommand = Command{
	Name: rootCommandIndicator,
	subCommands: []Command{
		{
			Name:        "connect",
			Description: "Connects to a device",
			action: func(_ Command, connectionString string) error {
				log.Info().Msgf("connect connectionString [%s]", connectionString)
				connectionUrl, err := url.Parse(connectionString)
				if err != nil {
					return errors.Wrapf(err, "can't parse connection url %s", connectionString)
				}
				addHost(connectionUrl.Host)
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
				for _, protocol := range protocolsSuggestions {
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
			Name: "read",
		},
		{
			Name: "write",
		},
		{
			Name:        "register",
			Description: "register a driver in the subsystem",
			action: func(_ Command, protocol string) error {
				switch protocol {
				case "ads":
					driverManager.RegisterDriver(ads.NewDriver())
					transports.RegisterTcpTransport(driverManager)
				case "bacnetip":
					driverManager.RegisterDriver(bacnetip.NewDriver())
					transports.RegisterUdpTransport(driverManager)
				case "c-bus":
					driverManager.RegisterDriver(cbus.NewDriver())
					transports.RegisterTcpTransport(driverManager)
				case "s7":
					driverManager.RegisterDriver(s7.NewDriver())
					transports.RegisterTcpTransport(driverManager)
				default:
					return errors.Errorf("Unknown protocol %s", protocol)
				}
				driverAdded(protocol)
				return nil
			},
			parameterSuggestions: func(currentText string) (entries []string) {
				for _, protocol := range protocolsSuggestions {
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
						AddItemHandler(func(event model.PlcSubscriptionEvent) {
							messagesReceived++
							_, _ = fmt.Fprintf(messageOutput, "[\"%d\"]\n%s[\"\"]", messagesReceived, event)
						}).
						Build()
					if err != nil {
						return errors.Wrapf(err, "%s can't subscribe", connectionsString)
					}
					subscriptionRequestResult := <-subscriptionRequest.Execute()
					if err := subscriptionRequestResult.GetErr(); err != nil {
						return errors.Wrapf(err, "%s can't subscribe", connectionsString)
					}
					log.Info().Msgf("subscription result %s", subscriptionRequestResult.GetResponse())
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
					Name: "TraceTransactionManagerWorkers",
					action: func(_ Command, argument string) error {
						switch argument {
						case "on":
							plc4x_config.TraceTransactionManagerWorkers = true
						case "off":
							plc4x_config.TraceTransactionManagerWorkers = false
						default:
							return errors.Errorf("illegal argument %s", argument)
						}
						return nil
					},
					parameterSuggestions: func(currentText string) (entries []string) {
						entries = append(entries, "on", "off")
						return
					},
				},
				{
					Name: "TraceTransactionManagerTransactions",
					action: func(_ Command, argument string) error {
						switch argument {
						case "on":
							plc4x_config.TraceTransactionManagerTransactions = true
						case "off":
							plc4x_config.TraceTransactionManagerTransactions = false
						default:
							return errors.Errorf("illegal argument %s", argument)
						}
						return nil
					},
					parameterSuggestions: func(currentText string) (entries []string) {
						entries = append(entries, "on", "off")
						return
					},
				},
				{
					Name: "TraceDefaultMessageCodecWorker",
					action: func(_ Command, argument string) error {
						switch argument {
						case "on":
							plc4x_config.TraceDefaultMessageCodecWorker = true
						case "off":
							plc4x_config.TraceDefaultMessageCodecWorker = false
						default:
							return errors.Errorf("illegal argument %s", argument)
						}
						return nil
					},
					parameterSuggestions: func(currentText string) (entries []string) {
						entries = append(entries, "on", "off")
						return
					},
				},
				{
					Name: "plc4xbrowser-debug",
					action: func(_ Command, argument string) error {
						switch argument {
						case "on":
							plc4xBrowserLog = zerolog.New(zerolog.ConsoleWriter{Out: tview.ANSIWriter(consoleOutput)})
						case "off":
							plc4xBrowserLog = zerolog.Nop()
						default:
							return errors.Errorf("illegal argument %s", argument)
						}
						return nil
					},
					parameterSuggestions: func(currentText string) (entries []string) {
						entries = append(entries, "on", "off")
						return
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
				_, _ = fmt.Fprintf(commandOutput, "%s [#00ff00]%s[white]: %s\n", indentString, command.Name, command.Description)
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
		addCommand(commandText)
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
		f(i+1, subCommand)
	}
}

func (c Command) String() string {
	return c.Name
}
