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

package analyzer

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/bacnetanalyzer"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/cbusanalyzer"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/common"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/pcaphandler"
	"github.com/google/gopacket/layers"
	"github.com/k0kubun/go-ansi"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"net"
)

func Analyze(pcapFile, protocolType, filter string, onlyParse, noBytesCompare bool, client string, verbosity int) {
	log.Info().Msgf("Analyzing pcap file '%s' with protocolType '%s' and filter '%s' now", pcapFile, protocolType, filter)

	handle, numberOfPackage, timestampToIndexMap := pcaphandler.GetIndexedPcapHandle(pcapFile, filter)
	log.Info().Msgf("Starting to analyze %d packages", numberOfPackage)
	defer handle.Close()
	log.Debug().Interface("handle", handle).Int("numberOfPackage", numberOfPackage).Msg("got handle")
	source := pcaphandler.GetPacketSource(handle)
	var packageParse func(common.PacketInformation, []byte) (interface{}, error)
	var serializePackage func(interface{}) ([]byte, error)
	var prettyPrint = func(item interface{}) {
		fmt.Printf("%v\n", item)
	}
	switch protocolType {
	case "bacnet":
		packageParse = bacnetanalyzer.PackageParse
		serializePackage = bacnetanalyzer.SerializePackage
	case "c-bus":
		analyzer := cbusanalyzer.Analyzer{Client: net.ParseIP(client)}
		packageParse = analyzer.PackageParse
		serializePackage = analyzer.SerializePackage
		prettyPrint = analyzer.PrettyPrint
	}
	bar := progressbar.NewOptions(numberOfPackage, progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan][1/3][reset] Analyzing packages..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	for packet := range source.Packets() {
		if packet == nil {
			log.Debug().Msg("Done reading packages. (nil returned)")
			break
		}
		if err := bar.Add(1); err != nil {
			log.Warn().Err(err).Msg("Error updating progressBar")
		}
		packetTimestamp := packet.Metadata().Timestamp
		realPacketNumber := timestampToIndexMap[packetTimestamp]
		description := fmt.Sprintf("No.[%d] timestamp: %v, %s", realPacketNumber, packetTimestamp, pcapFile)
		packetInformation := common.PacketInformation{
			PacketNumber:    realPacketNumber,
			PacketTimestamp: packetTimestamp,
			Description:     description,
		}
		if networkLayer, ok := packet.NetworkLayer().(*layers.IPv4); ok {
			packetInformation.SrcIp = networkLayer.SrcIP
			packetInformation.DstIp = networkLayer.DstIP
		}

		applicationLayer := packet.ApplicationLayer()
		if applicationLayer == nil {
			log.Info().Stringer("packetInformation", packetInformation).Msgf("No.[%d] No application layer", realPacketNumber)
			continue
		}
		payload := applicationLayer.Payload()
		if parsed, err := packageParse(packetInformation, payload); err != nil {
			// TODO: write report to xml or something
			log.Warn().Stringer("packetInformation", packetInformation).Err(err).Msgf("No.[%d] Error parsing package", realPacketNumber)
			continue
		} else {
			log.Info().Stringer("packetInformation", packetInformation).Msgf("No.[%d] Parsed", realPacketNumber)
			if verbosity > 1 {
				prettyPrint(parsed)
			}
			if onlyParse {
				log.Trace().Msg("only parsing")
				continue
			}
			serializedBytes, err := serializePackage(parsed)
			if err != nil {
				// TODO: write report to xml or something
				log.Warn().Stringer("packetInformation", packetInformation).Err(err).Msgf("No.[%d] Error serializing", realPacketNumber)
				continue
			}
			if noBytesCompare {
				log.Trace().Msg("not comparing bytes")
				continue
			}
			if compareResult := bytes.Compare(payload, serializedBytes); compareResult != 0 {
				// TODO: write report to xml or something
				log.Warn().Stringer("packetInformation", packetInformation).Msgf("No.[%d] Bytes don't match", realPacketNumber)
				if verbosity > 0 {
					println("Original bytes")
					println(hex.Dump(payload))
					println("Serialized bytes")
					println(hex.Dump(serializedBytes))
				}
			}
		}
	}

	log.Info().Msgf("Done evaluating %d packages", numberOfPackage)
}
