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

package tests

import (
	"os"
	"testing"

	"github.com/apache/plc4x/plc4go/internal/eip"
	"github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/apache/plc4x/plc4go/protocols/eip/readwrite"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	_ "github.com/apache/plc4x/plc4go/tests/initializetest"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestEIPDriver(t *testing.T) {
	log.Logger = log.
		With().Caller().Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		Level(zerolog.DebugLevel)
	config.TraceTransactionManagerWorkers = true
	config.TraceTransactionManagerTransactions = true
	config.TraceDefaultMessageCodecWorker = true
	testutils.RunDriverTestsuite(t, eip.NewDriver(), "assets/testing/protocols/eip/DriverTestsuite.xml", readwrite.EipXmlParserHelper{})
}
