/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package pcaphandler

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"time"
)

// GetPacketSource gets a packet source from a handle
func GetPacketSource(handle *pcap.Handle) *gopacket.PacketSource {
	return gopacket.NewPacketSource(handle, handle.LinkType())
}

// GetIndexedPcapHandle returns a *pcap.Handle, the number of packages found and an index which maps timestamp to
// absolute package number and panics if an error occurs
func GetIndexedPcapHandle(file, filterExpression string) (*pcap.Handle, int, map[time.Time]int) {
	timestampToIndexMap := make(map[time.Time]int)
	// Count absolute packages and set timestamp map
	temporaryHandle := GetPcapHandle(file, "")
	defer temporaryHandle.Close()
	packetSource := GetPacketSource(temporaryHandle)
	packages := 0
	for packet := range packetSource.Packets() {
		if packet == nil {
			break
		}
		packages++
		timestampToIndexMap[packet.Metadata().Timestamp] = packages
	}
	// Just count filtered packages
	temporaryFilteredHandle := GetPcapHandle(file, filterExpression)
	defer temporaryFilteredHandle.Close()
	filteredPacketSource := GetPacketSource(temporaryFilteredHandle)
	packages = 0
	for packet := range filteredPacketSource.Packets() {
		if packet == nil {
			break
		}
		packages++
	}
	return GetPcapHandle(file, filterExpression), packages, timestampToIndexMap
}

// GetPcapHandle returns a *pcap.Handle and panics if an error occurs
func GetPcapHandle(file, filterExpression string) *pcap.Handle {
	handle := getPcapHandleOrPanic(file)
	if filterExpression != "" {
		if err := handle.SetBPFFilter(filterExpression); err != nil {
			panic(err)
		}
	}
	return handle
}

func getPcapHandleOrPanic(file string) *pcap.Handle {
	handle, err := pcap.OpenOffline(file)
	if err != nil {
		panic(err)
	}
	return handle
}
