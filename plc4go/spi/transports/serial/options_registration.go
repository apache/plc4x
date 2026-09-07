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

package serial

import "github.com/apache/plc4x/plc4go/spi/options"

// The options this transport reads, declared beside the code that reads them so a driver's report
// of unread connection-string options does not name them.
//
// A user addresses a transport's option under the transport's own code, as PLC4J does and as the
// documentation says: "tcp.connect-timeout-ms", not "connect-timeout-ms". An option a driver
// injects into the map itself is not addressed by anyone and carries no prefix.
func init() {
	options.RegisterTransportOptions(
		"serial.baud-rate", "serial.data-bits", "serial.stop-bits", "serial.parity",
		"serial.flow-control", "serial.dtr", "serial.rts", "serial.read-timeout-ms",
		"serial.write-timeout-ms", "serial.connect-timeout-ms", "serial.reuse-port",
		"serial.interframe-delay",
	)
}
