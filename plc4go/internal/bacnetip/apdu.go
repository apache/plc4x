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

package bacnetip

type APCI interface {
	PCI
}

type _APCI struct {
	*_PCI
	*DebugContents
}

type APDU interface {
	APCI
	PDUData

	GetApduInvokeID() *uint8 // TODO: check if we really need a pointer
}

type __APDU struct {
	*_APCI
}

type _APDU interface {
	APDU
}

type ___APDU struct {
	*__APDU
}

//TODO: implement it...
