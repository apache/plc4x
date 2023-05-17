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
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"

	"github.com/gdamore/tcell/v2"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

func SetupApplication() *tview.Application {
	application := tview.NewApplication()

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	fileArea := buildFileArea(newPrimitive, application)
	outputArea := buildOutputArea(newPrimitive, application)
	commandArea := buildCommandArea(newPrimitive, application)

	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive("PLC4X PCAP Analyzer"), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("https://github.com/apache/plc4x"), 2, 0, 1, 3, 0, 0, false)

	// Layout for screens narrower than 100 cells (fileArea and side bar are hidden).
	grid.AddItem(fileArea, 0, 0, 0, 0, 0, 0, false).
		AddItem(outputArea, 1, 0, 1, 3, 0, 0, false).
		AddItem(commandArea, 0, 0, 0, 0, 0, 0, true)

	// Layout for screens wider than 100 cells.
	grid.AddItem(fileArea, 1, 0, 1, 1, 0, 100, false).
		AddItem(outputArea, 1, 1, 1, 1, 0, 100, false).
		AddItem(commandArea, 1, 2, 1, 1, 0, 100, false)

	application.SetRoot(grid, true).EnableMouse(true)

	return application
}

func buildFileArea(newPrimitive func(text string) tview.Primitive, application *tview.Application) tview.Primitive {
	connectionAreaHeader := newPrimitive("Files")
	connectionArea := tview.NewGrid().
		SetRows(3, 0, 10).
		SetColumns(0).
		AddItem(connectionAreaHeader, 0, 0, 1, 1, 0, 0, false)
	{
		fileList := tview.NewList()
		loadedPcapFilesChanged = func() {
			application.QueueUpdateDraw(func() {
				fileList.Clear()
				for _, pcapFile := range loadedPcapFiles {
					fileList.AddItem(pcapFile.name, pcapFile.path, 0x0, func() {
						//TODO: disconnect popup
						_ = pcapFile
					})
				}
			})
		}
		connectionArea.AddItem(fileList, 1, 0, 1, 1, 0, 0, false)
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
		enteredCommandsView := tview.NewTextView().
			SetDynamicColors(true).
			SetRegions(true).
			SetWordWrap(true).
			SetChangedFunc(func() {
				application.Draw()
			})
		commandOutput = enteredCommandsView
		commandOutputClear = func() {
			enteredCommandsView.SetText("")
		}

		commandArea.AddItem(enteredCommandsView, 1, 0, 1, 1, 0, 0, false)

		commandInputField := tview.NewInputField().
			SetLabel("$").
			SetFieldWidth(30)
		application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyCtrlC:
				commandInputField.SetText("")
				application.SetFocus(commandInputField)
				return nil
			case tcell.KeyCtrlD:
				// TODO: maybe add a modal here
				application.Stop()
				return nil
			}
			return event
		})
		commandInputField.
			SetDoneFunc(func(key tcell.Key) {
				commandText := commandInputField.GetText()
				if commandText == "quit" {
					// TODO: maybe add a modal here
					application.Stop()
					return
				}
				commandsExecuted++
				go func() {
					commandHistoryShortcut, _ := regexp.Compile("^[0-9]$")
					if commandHistoryShortcut.MatchString(commandText) {
						atoi, _ := strconv.Atoi(commandHistoryShortcut.FindString(commandText))
						if atoi < len(config.History.Last10Commands) {
							commandText = config.History.Last10Commands[atoi]
						} else {
							_, _ = fmt.Fprintf(enteredCommandsView, "[#ff0000]%s %s[white]\n", time.Now().Format("04:05"), errors.Errorf("No such elements %d in command history", atoi))
							return
						}
					}
					_, _ = fmt.Fprintf(enteredCommandsView, "%s [\"%d\"]%s[\"\"]\n", time.Now().Format("04:05"), commandsExecuted, commandText)
					ctx, cancelFunc := context.WithCancel(rootContext)
					randomId := rand.Uint32()
					cancelFunctions[randomId] = cancelFunc
					defer delete(cancelFunctions, randomId)

					if err := Execute(ctx, commandText); err != nil {
						_, _ = fmt.Fprintf(enteredCommandsView, "[#ff0000]%s %s[white]\n", time.Now().Format("04:05"), err)
						return
					}
					application.QueueUpdateDraw(func() {
						commandInputField.SetText("")
					})
				}()
			})
		commandInputField.SetAutocompleteFunc(rootCommand.Completions)

		enteredCommandsView.SetDoneFunc(func(key tcell.Key) {
			currentSelection := enteredCommandsView.GetHighlights()
			if key == tcell.KeyEnter {
				if len(currentSelection) > 0 {
					enteredCommandsView.Highlight()
				} else {
					enteredCommandsView.Highlight("0").ScrollToHighlight()
				}
				if len(currentSelection) == 1 {
					commandInputField.SetText(enteredCommandsView.GetRegionText(currentSelection[0]))
					application.SetFocus(commandInputField)
				}
			} else if len(currentSelection) > 0 {
				index, _ := strconv.Atoi(currentSelection[0])
				if key == tcell.KeyTab {
					index = (index + 1) % commandsExecuted
				} else if key == tcell.KeyBacktab {
					index = (index - 1 + commandsExecuted) % commandsExecuted
				} else {
					return
				}
				enteredCommandsView.Highlight(strconv.Itoa(index)).ScrollToHighlight()
			}
		})

		commandArea.AddItem(commandInputField, 2, 0, 1, 1, 0, 0, true)
	}
	return commandArea
}

func buildOutputArea(newPrimitive func(text string) tview.Primitive, application *tview.Application) *tview.Grid {
	outputAreaHeader := newPrimitive("Output")
	outputArea := tview.NewGrid().
		SetRows(3, 0, 10).
		SetColumns(0, 30).
		AddItem(outputAreaHeader, 0, 0, 1, 1, 0, 0, false)
	{
		var jumpToMessageItem func(messageNumber int) bool
		{
			outputView := tview.NewTextView().
				SetDynamicColors(true).
				SetRegions(true).
				SetWordWrap(false).
				SetWrap(false).
				SetChangedFunc(func() {
					application.Draw()
				})
			jumpToMessageItem = func(messageNumber int) bool {
				regionId := strconv.Itoa(messageNumber)
				if outputView.GetRegionText(regionId) == "" {
					return false
				}
				outputView.Highlight(regionId).ScrollToHighlight()
				return true
			}
			messageOutput = outputView
			messageOutputClear = func() {
				outputView.SetText("")
			}

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
						index = (index + 1) % numberOfMessagesReceived
					} else if key == tcell.KeyBacktab {
						index = (index - 1 + numberOfMessagesReceived) % numberOfMessagesReceived
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
				SetMaxLines(config.MaxConsoleLines).
				SetChangedFunc(func() {
					application.Draw()
				})
			consoleOutput = consoleView
			consoleOutputClear = func() {
				consoleView.SetText("")
			}

			consoleView.SetBorder(false)
			outputArea.AddItem(consoleView, 2, 0, 1, 1, 0, 0, false)
		}

		{
			receivedMessagesList := tview.NewList()
			messageReceived = func(messageNumber int, receiveTime time.Time, message apiModel.PlcMessage) {
				application.QueueUpdateDraw(func() {
					receivedMessagesList.AddItem(fmt.Sprintf("No %d @%s (api)", messageNumber, receiveTime.Format("15:04:05.999999")), "", 0x0, func() {
						if ok := jumpToMessageItem(messageNumber); !ok {
							plc4xpcapanalyzerLog.Debug().Msgf("Adding new message to console output")
							_, _ = fmt.Fprintf(messageOutput, "Message nr: %d\n[\"%d\"]%s[\"\"]\n", messageNumber, messageNumber, message)
							jumpToMessageItem(messageNumber)
						}
					})
				})
			}
			spiMessageReceived = func(messageNumber int, receiveTime time.Time, message spi.Message) {
				application.QueueUpdateDraw(func() {
					receivedMessagesList.AddItem(fmt.Sprintf("No %d @%s (spi)", messageNumber, receiveTime.Format("15:04:05.999999")), "", 0x0, func() {
						if ok := jumpToMessageItem(messageNumber); !ok {
							plc4xpcapanalyzerLog.Debug().Msgf("Adding new spi message to console output")
							_, _ = fmt.Fprintf(messageOutput, "Message nr: %d\n[\"%d\"]%s[\"\"]\n", messageNumber, messageNumber, message)
							jumpToMessageItem(messageNumber)
						}
					})
				})
			}
			outputArea.AddItem(receivedMessagesList, 0, 1, 3, 1, 0, 0, false)
		}
	}
	return outputArea
}
