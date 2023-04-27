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

package initializetest

import (
	"github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
)

func init() {
	SetupDefaultLoggerForTest()
}

// SetupDefaultLoggerForTest This is the default log setting for Jenkins
func SetupDefaultLoggerForTest() {
	onJenkins := os.Getenv("JENKINS_URL") != ""
	log.Logger = log.
		//// Enable below if you want to see the filenames
		//With().Caller().Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: onJenkins}).
		Level(zerolog.InfoLevel)
}

// GoFullDebug can be called to set up logger for debug (the call to this should not be committed)
func GoFullDebug() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.
		With().Caller().Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		Level(zerolog.DebugLevel)
	config.TraceTransactionManagerWorkers = false
	config.TraceTransactionManagerTransactions = false
	config.TraceDefaultMessageCodecWorker = false
}
