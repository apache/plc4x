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

package globals

import "os"

var (
	// ExtendedGeneralOutput switches extended output on for general items
	ExtendedGeneralOutput bool

	// ExtendedPDUOutput switches the PDU output to an extended format for debugging
	ExtendedPDUOutput bool

	// LogComm enables logging for communications
	LogComm bool

	// LogAppService enables logging for application services
	LogAppService bool

	// LogPDU enables logging for pdu
	LogPDU bool

	// LogVlan enables logging for vlan
	LogVlan bool
)

func init() {
	ExtendedGeneralOutput = os.Getenv("ExtendedGeneralOutput") == "true"
	ExtendedPDUOutput = os.Getenv("ExtendedPDUOutput") == "true"
	LogComm = os.Getenv("LogComm") == "true"
	LogAppService = os.Getenv("LogAppService") == "true"
	LogPDU = os.Getenv("LogPDU") == "true"
	LogVlan = os.Getenv("LogVlan") == "true"
}
