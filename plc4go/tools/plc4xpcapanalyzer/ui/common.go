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
	"github.com/apache/plc4x/plc4go/internal/spi"
	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	plc4goModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/rs/zerolog"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

const protocols = "ads,bacnetip,c-bus,s7"

var protocolList = strings.Split(protocols, ",")

var plc4xpcapanalyzerLog = zerolog.Nop()

var driverManager plc4go.PlcDriverManager
var driverAdded func(string)

type loadedPcapFile struct {
	name string
	path string
}

var loadedPcapFiles []loadedPcapFile
var loadedPcapFilesChanged func()

var messageReceived func(messageNumber int, receiveTime time.Time, message plc4goModel.PlcMessage)
var numberOfMessagesReceived int
var spiMessageReceived func(messageNumber int, receiveTime time.Time, message spi.Message)
var spiNumberOfMessagesReceived int
var messageOutput io.Writer
var messageOutputClear func()

var consoleOutput io.Writer
var consoleOutputClear func()

var commandsExecuted int
var commandOutput io.Writer
var commandOutputClear func()

var currentDirChanged func()
var currentDir = func() string {
	dir, _ := os.Getwd()
	return dir
}()

var shutdownMutex sync.Mutex
var hasShutdown bool

func Shutdown() {
	shutdownMutex.Lock()
	defer shutdownMutex.Unlock()
	if hasShutdown {
		return
	}
	saveConfig()
	hasShutdown = true
}
