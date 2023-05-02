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
	"context"
	"testing"

	"github.com/apache/plc4x/plc4go/internal/knxnetip"
	knxIO "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite"
	knxModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/utils"
	_ "github.com/apache/plc4x/plc4go/tests/initializetest"
)

func TestKNXNetIPDriver(t *testing.T) {
	t.Skip("No test yet")
	parser := func(readBufferByteBased utils.ReadBufferByteBased) (any, error) {
		return knxModel.KnxNetIpMessageParseWithBuffer(context.Background(), readBufferByteBased)
	}
	testutils.RunDriverTestsuite(t, knxnetip.NewDriver(), "assets/testing/protocols/knxnetip/DriverTestsuite.xml", knxIO.KnxnetipXmlParserHelper{}, testutils.WithRootTypeParser(parser))
}
