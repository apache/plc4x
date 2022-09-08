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
	"github.com/apache/plc4x/plc4go/internal/ads"
	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	"github.com/apache/plc4x/plc4go/internal/cbus"
	"github.com/apache/plc4x/plc4go/internal/s7"
	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/transports"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitSubsystem() {
	logLevel := zerolog.InfoLevel
	if configuredLevel := config.LogLevel; configuredLevel != "" {
		if parsedLevel, err := zerolog.ParseLevel(configuredLevel); err != nil {
			panic(err)
		} else {
			logLevel = parsedLevel
		}
	}
	driverManager = plc4go.NewPlcDriverManager()

	log.Logger = log.
		//// Enable below if you want to see the filenames
		//With().Caller().Logger().
		Output(zerolog.ConsoleWriter{Out: tview.ANSIWriter(consoleOutput)}).
		Level(logLevel)

	// We offset the commands executed with the last commands
	commandsExecuted = len(config.History.Last10Commands)
	outputCommandHistory()

	for _, driver := range config.AutoRegisterDrivers {
		log.Info().Msgf("Auto register driver %s", driver)
		if err := validateDriverParam(driver); err != nil {
			log.Err(err).Msgf("Invalid configuration")
			continue
		}
		_ = registerDriver(driver)
	}
}

func outputCommandHistory() {
	_, _ = fmt.Fprintln(commandOutput, "[#0000ff]Last 10 commands[white]")
	for i, command := range config.History.Last10Commands {
		_, _ = fmt.Fprintf(commandOutput, "   [#00ff00]%d[white]: [\"%d\"]%s[\"\"]\n", i, i, tview.Escape(command))
	}
}

func validateDriverParam(driver string) error {
	for _, protocol := range protocolList {
		if protocol == driver {
			return nil
		}
	}
	return errors.Errorf("protocol %s not found", driver)
}

var tcpRegistered, udpRegistered bool

func registerDriver(driverId string) error {
	if _, ok := registeredDrivers[driverId]; ok {
		return errors.Errorf("%s already registered", driverId)
	}
	var driver plc4go.PlcDriver
	switch driverId {
	case "ads":
		driver = ads.NewDriver()
		driverManager.RegisterDriver(driver)
		if !tcpRegistered {
			transports.RegisterTcpTransport(driverManager)
			tcpRegistered = true
		}
	case "bacnetip":
		driver = bacnetip.NewDriver()
		driverManager.RegisterDriver(driver)
		if !udpRegistered {
			transports.RegisterUdpTransport(driverManager)
			udpRegistered = true
		}
	case "c-bus":
		driver = cbus.NewDriver()
		driverManager.RegisterDriver(driver)
		if !tcpRegistered {
			transports.RegisterTcpTransport(driverManager)
			tcpRegistered = true
		}
	case "s7":
		driver = s7.NewDriver()
		driverManager.RegisterDriver(driver)
		if !tcpRegistered {
			transports.RegisterTcpTransport(driverManager)
			tcpRegistered = true
		}
	default:
		return errors.Errorf("Unknown driver %s", driverId)
	}
	registeredDrivers[driverId] = driver
	go driverAdded(driver)
	return nil
}
