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
	"bytes"
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/ads"
	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	"github.com/apache/plc4x/plc4go/internal/cbus"
	"github.com/apache/plc4x/plc4go/internal/s7"
	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/transports/pcap"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"strings"
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
		Output(zerolog.NewConsoleWriter(
			func(w *zerolog.ConsoleWriter) {
				w.Out = tview.ANSIWriter(consoleOutput)
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
		)).
		Level(logLevel)

	// We offset the commands executed with the last commands
	commandsExecuted = len(config.History.Last10Commands)
	outputCommandHistory()

	for _, driver := range config.AutoRegisterDrivers {
		log.Info().Str("driver", driver).Msg("Auto register driver")
		if err := validateDriverParam(driver); err != nil {
			log.Err(err).Msg("Invalid configuration")
			continue
		}
		_ = registerDriver(driver)
	}
}

func OpenFile(pcapFile string) error {
	if !strings.HasPrefix(pcapFile, string(os.PathSeparator)) {
		pcapFile = path.Join(currentDir, pcapFile)
	}
	var name string
	if stat, err := os.Stat(pcapFile); err != nil {
		return err
	} else if stat.IsDir() {
		stat.Name()
		return errors.Errorf("%s is a dir", pcapFile)
	} else {
		name = stat.Name()
	}
	for _, loadedPcapFile := range loadedPcapFiles {
		if loadedPcapFile.path == pcapFile {
			return errors.Errorf("%s already loaded", name)
		}
	}
	addRecentFilesEntry(pcapFile)
	loadedPcapFiles = append(loadedPcapFiles, loadedPcapFile{
		name: name,
		path: pcapFile,
	})
	loadedPcapFilesChanged()
	return nil
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

func registerDriver(driver string) error {
	switch driver {
	case "ads":
		driverManager.RegisterDriver(ads.NewDriver())
	case "bacnetip":
		driverManager.RegisterDriver(bacnetip.NewDriver())
	case "c-bus":
		driverManager.RegisterDriver(cbus.NewDriver())
	case "s7":
		driverManager.RegisterDriver(s7.NewDriver())
	default:
		return errors.Errorf("Unknown driver %s", driver)
	}
	driverManager.(spi.TransportAware).RegisterTransport(pcap.NewTransport())
	go driverAdded(driver)
	return nil
}
