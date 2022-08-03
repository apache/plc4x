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
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
)

var driverManager plc4go.PlcDriverManager
var driverAdded func(string)
var connections map[string]plc4go.PlcConnection
var connectionsChanged func()

var commandsExecuted int

var messagesReceived int
var messageOutput io.Writer

var consoleOutput io.Writer

var commandOutput io.Writer

func init() {
	hasShutdown = false
	connections = make(map[string]plc4go.PlcConnection)
}

func initSubsystem() {
	driverManager = plc4go.NewPlcDriverManager()

	logLevel := zerolog.InfoLevel
	if configuredLevel := config.LogLevel; configuredLevel != "" {
		if parsedLevel, err := zerolog.ParseLevel(configuredLevel); err != nil {
			panic(err)
		} else {
			logLevel = parsedLevel
		}
	}

	log.Logger = log.
		//// Enable below if you want to see the filenames
		//With().Caller().Logger().
		Output(zerolog.ConsoleWriter{Out: tview.ANSIWriter(consoleOutput)}).
		Level(logLevel)
}

var shutdownMutex sync.Mutex
var hasShutdown bool

func shutdown() {
	shutdownMutex.Lock()
	defer shutdownMutex.Unlock()
	if hasShutdown {
		return
	}
	for _, connection := range connections {
		connection.Close()
	}
	saveConfig()
	hasShutdown = true
}

func main() {
	application := tview.NewApplication()

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	connectionArea := buildConnectionArea(newPrimitive, application)
	outputArea := buildOutputArea(newPrimitive, application)
	commandArea := buildCommandArea(newPrimitive, application)

	grid := tview.NewGrid().
		SetRows(3, 0, 1).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive("PLC4X Browser"), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("https://github.com/apache/plc4x"), 2, 0, 1, 3, 0, 0, false)

	// Layout for screens narrower than 100 cells (connectionArea and side bar are hidden).
	grid.AddItem(connectionArea, 0, 0, 0, 0, 0, 0, false).
		AddItem(outputArea, 1, 0, 1, 3, 0, 0, false).
		AddItem(commandArea, 0, 0, 0, 0, 0, 0, true)

	// Layout for screens wider than 100 cells.
	grid.AddItem(connectionArea, 1, 0, 1, 1, 0, 100, false).
		AddItem(outputArea, 1, 1, 1, 1, 0, 100, false).
		AddItem(commandArea, 1, 2, 1, 1, 0, 100, false)

	application.SetRoot(grid, true).EnableMouse(true)

	loadConfig()

	initSubsystem()

	application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			shutdown()
		}
		return event
	})

	if err := application.Run(); err != nil {
		panic(err)
	}
	shutdown()
}

func buildConnectionArea(newPrimitive func(text string) tview.Primitive, application *tview.Application) tview.Primitive {
	connectionAreaHeader := newPrimitive("Connections")
	connectionArea := tview.NewGrid().
		SetRows(3, 0, 10).
		SetColumns(0).
		AddItem(connectionAreaHeader, 0, 0, 1, 1, 0, 0, false)
	{
		connectionList := tview.NewList()
		connectionsChanged = func() {
			application.QueueUpdateDraw(func() {
				connectionList.Clear()
				for connectionString, connection := range connections {
					connectionList.AddItem(connectionString, "", 0x0, func() {
						//TODO: disconnect popup
						_ = connection
					})
				}
			})
		}
		connectionArea.AddItem(connectionList, 1, 0, 1, 1, 0, 0, false)
		{
			registeredDriverAreaHeader := newPrimitive("Registered drivers")
			registeredDriverArea := tview.NewGrid().
				SetRows(3, 0).
				SetColumns(0).
				AddItem(registeredDriverAreaHeader, 0, 0, 1, 1, 0, 0, false)
			{
				driverList := tview.NewList()
				driverAdded = func(driver string) {
					application.QueueUpdateDraw(func() {
						driverList.AddItem(driver, "", 0x0, func() {
							//TODO: disconnect popup
						})
					})
				}
				registeredDriverArea.AddItem(driverList, 1, 0, 1, 1, 0, 0, false)
			}
			connectionArea.AddItem(registeredDriverArea, 2, 0, 1, 1, 0, 0, false)
		}

	}
	return connectionArea
}

func buildCommandArea(newPrimitive func(text string) tview.Primitive, application *tview.Application) tview.Primitive {
	commandAreaHeader := newPrimitive("Commands")
	commandArea := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(0).
		AddItem(commandAreaHeader, 0, 0, 1, 1, 0, 0, false)
	{
		enteredCommands := tview.NewTextView().
			SetDynamicColors(true).
			SetRegions(true).
			SetWordWrap(true).
			SetChangedFunc(func() {
				application.Draw()
			})
		commandOutput = enteredCommands

		commandArea.AddItem(enteredCommands, 1, 0, 1, 1, 0, 0, false)

		commandInputField := tview.NewInputField().
			SetLabel("$").
			SetFieldWidth(30)
		commandInputField.
			SetDoneFunc(func(key tcell.Key) {
				commandText := commandInputField.GetText()
				if commandText == "quit" {
					// TODO: maybe add a modal here
					application.Stop()
					return
				}
				commandsExecuted++
				_, _ = fmt.Fprintf(enteredCommands, "%s [\"%d\"]%s[\"\"]\n", time.Now().Format("04:05"), commandsExecuted, commandText)
				go func() {
					if err := Execute(commandText); err != nil {
						_, _ = fmt.Fprintf(enteredCommands, "[#ff0000]%s %s[white]\n", time.Now().Format("04:05"), err)
						return
					}
					application.QueueUpdateDraw(func() {
						commandInputField.SetText("")
					})
				}()
			})
		commandInputField.SetAutocompleteFunc(rootCommand.Completions)
		commandArea.AddItem(commandInputField, 2, 0, 1, 1, 0, 0, true)
	}
	return commandArea
}

func buildOutputArea(newPrimitive func(text string) tview.Primitive, application *tview.Application) *tview.Grid {
	outputAreaHeader := newPrimitive("Output")
	outputArea := tview.NewGrid().
		SetRows(3, 0, 10).
		SetColumns(0).
		AddItem(outputAreaHeader, 0, 0, 1, 1, 0, 0, false)
	{
		{
			outputView := tview.NewTextView().
				SetDynamicColors(true).
				SetRegions(true).
				SetWordWrap(true).
				SetChangedFunc(func() {
					application.Draw()
				})
			messageOutput = outputView

			outputView.SetDoneFunc(func(key tcell.Key) {
				currentSelection := outputView.GetHighlights()
				if key == tcell.KeyEnter {
					if len(currentSelection) > 0 {
						outputView.Highlight()
					} else {
						outputView.Highlight("0").ScrollToHighlight()
					}
				} else if len(currentSelection) > 0 {
					index, _ := strconv.Atoi(currentSelection[0])
					if key == tcell.KeyTab {
						index = (index + 1) % messagesReceived
					} else if key == tcell.KeyBacktab {
						index = (index - 1 + messagesReceived) % messagesReceived
					} else {
						return
					}
					outputView.Highlight(strconv.Itoa(index)).ScrollToHighlight()
				}
			})
			outputView.SetBorder(false)
			outputArea.AddItem(outputView, 1, 0, 1, 1, 0, 0, false)
		}

		{
			consoleView := tview.NewTextView().
				SetDynamicColors(true).
				SetChangedFunc(func() {
					application.Draw()
				})
			consoleOutput = consoleView

			consoleView.SetBorder(false)
			outputArea.AddItem(consoleView, 2, 0, 1, 1, 0, 0, false)
		}
	}
	return outputArea
}
