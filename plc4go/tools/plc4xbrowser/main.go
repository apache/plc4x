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
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/apache/plc4x/plc4go/internal/cbus"
	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/transports"
)

// TODO: replace with real commands
const plc4xCommands = "connect,disconnect,read,write,register,subscribe,quit"

var connections map[string]plc4go.PlcConnection
var connectionsChanged func()

var commandsExecuted int

var messagesReceived int
var messageOutput io.Writer

func init() {
	connections = make(map[string]plc4go.PlcConnection)
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
		SetRows(3, 0, 3).
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

	if err := application.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
	for _, connection := range connections {
		connection.Close()
	}
}

func buildConnectionArea(newPrimitive func(text string) tview.Primitive, application *tview.Application) tview.Primitive {
	connectionAreaHeader := newPrimitive("Connections")
	connectionArea := tview.NewGrid().
		SetRows(3, 0).
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
		connectionArea.AddItem(connectionList, 1, 0, 1, 1, 0, 0, true)
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
		commandArea.AddItem(enteredCommands, 1, 0, 1, 1, 0, 0, false)

		words := strings.Split(plc4xCommands, ",")
		commandInputField := tview.NewInputField().
			SetLabel("PLC4X Command").
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
					if err := handleCommand(commandText); err != nil {
						_, _ = fmt.Fprintf(enteredCommands, "[#ff0000]%s %s[white]\n", time.Now().Format("04:05"), err)
						return
					}
					application.QueueUpdateDraw(func() {
						commandInputField.SetText("")
					})
				}()
			})
		commandInputField.SetAutocompleteFunc(func(currentText string) (entries []string) {
			if len(currentText) == 0 {
				return
			}
			for _, word := range words {
				if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) {
					entries = append(entries, word)
				}
			}
			switch {
			case strings.HasPrefix(currentText, "disconnect"):
				for connectionsString, _ := range connections {
					entries = append(entries, "disconnect "+connectionsString)
				}
			case strings.HasPrefix(currentText, "subscribe"):
				for connectionsString, _ := range connections {
					entries = append(entries, "subscribe "+connectionsString)
				}
			}
			return
		})
		commandArea.AddItem(commandInputField, 2, 0, 1, 1, 0, 0, true)
	}
	return commandArea
}

func handleCommand(commandText string) error {
	switch {
	case strings.HasPrefix(commandText, "register "):
	case strings.HasPrefix(commandText, "connect "):
		host := strings.TrimPrefix(commandText, "connect ")
		if _, ok := connections[host]; ok {
			return errors.Errorf("%s already connected", host)
		}
		//TODO: we hardcode that to cbus for now
		connectionString := fmt.Sprintf("c-bus://%s?srchk=true", host)
		driverManager := plc4go.NewPlcDriverManager()
		driverManager.RegisterDriver(cbus.NewDriver())
		transports.RegisterTcpTransport(driverManager)
		connectionResult := <-driverManager.GetConnection(connectionString)
		if err := connectionResult.GetErr(); err != nil {
			return errors.Wrapf(err, "%s can't connect to", host)
		}
		log.Info().Msgf("%s connected", host)
		connections[host] = connectionResult.GetConnection()
		connectionsChanged()
	case strings.HasPrefix(commandText, "disconnect "):
		host := strings.TrimPrefix(commandText, "disconnect ")
		if connection, ok := connections[host]; !ok {
			return errors.Errorf("%s not connected", host)
		} else {
			closeResult := <-connection.Close()
			log.Info().Msgf("%s disconnected", host)
			delete(connections, host)
			connectionsChanged()
			if err := closeResult.GetErr(); err != nil {
				return errors.Wrapf(err, "%s can't close", host)
			}
		}
	case strings.HasPrefix(commandText, "subscribe "):
		host := strings.TrimPrefix(commandText, "subscribe ")
		if connection, ok := connections[host]; !ok {
			return errors.Errorf("%s not connected", host)
		} else {
			subscriptionRequest, err := connection.SubscriptionRequestBuilder().
				AddEventQuery("something", "monitor/*/*").
				AddItemHandler(func(event model.PlcSubscriptionEvent) {
					messagesReceived++
					_, _ = fmt.Fprintf(messageOutput, "[\"%d\"]\n%s[\"\"]", messagesReceived, event)
				}).
				Build()
			if err != nil {
				return errors.Wrapf(err, "%s can't subscribe", host)
			}
			subscriptionRequestResult := <-subscriptionRequest.Execute()
			if err := subscriptionRequestResult.GetErr(); err != nil {
				return errors.Wrapf(err, "%s can't subscribe", host)
			}
			log.Info().Msgf("subscription result %s", subscriptionRequestResult.GetResponse())
		}
	default:
		return errors.Errorf("%s not found", commandText)
	}
	return nil
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
				SetRegions(true).
				SetWordWrap(true).
				SetChangedFunc(func() {
					application.Draw()
				})

			log.Logger = log.
				//// Enable below if you want to see the filenames
				//With().Caller().Logger().
				Output(zerolog.ConsoleWriter{Out: tview.ANSIWriter(consoleView)}).
				Level(zerolog.InfoLevel)

			consoleView.SetBorder(false)
			outputArea.AddItem(consoleView, 2, 0, 1, 1, 0, 0, false)
		}
	}
	return outputArea
}
