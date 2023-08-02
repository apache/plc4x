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
	"context"
	"encoding/hex"
	"fmt"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/config"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/bacnetanalyzer"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/cbusanalyzer"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/common"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/pcaphandler"
	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/k0kubun/go-ansi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"io"
	"net"
	"os"
	"time"
)

func Analyze(pcapFile, protocolType string) error {
	return AnalyzeWithOutput(pcapFile, protocolType, os.Stdout, os.Stderr)
}

func AnalyzeWithOutput(pcapFile, protocolType string, stdout, stderr io.Writer) error {
	return AnalyzeWithOutputAndCallback(context.TODO(), pcapFile, protocolType, stdout, stderr, nil)
}

func AnalyzeWithOutputAndCallback(ctx context.Context, pcapFile, protocolType string, stdout, stderr io.Writer, messageCallback func(parsed spi.Message)) error {
	var filterExpression = config.AnalyzeConfigInstance.Filter
	if filterExpression != "" {
		log.Info().Str("filterExpression", filterExpression).Msg("Using global filter")
	}
	var mapPackets = func(in chan gopacket.Packet, packetInformationCreator func(packet gopacket.Packet) common.PacketInformation) chan gopacket.Packet {
		return in
	}
	var packageParse func(common.PacketInformation, []byte) (spi.Message, error)
	var serializePackage func(spi.Message) ([]byte, error)
	var prettyPrint = func(item spi.Message) {
		_, _ = fmt.Fprintf(stdout, "%v\n", item)
	}
	var byteOutput = hex.Dump
	switch protocolType {
	case "bacnetip":
		if !config.AnalyzeConfigInstance.NoFilter {
			if config.AnalyzeConfigInstance.Filter == "" && config.BacnetConfigInstance.BacnetFilter != "" {
				log.Debug().Str("filter", config.BacnetConfigInstance.Filter).Msg("Setting bacnet filter")
				filterExpression = config.BacnetConfigInstance.BacnetFilter
			}
		} else {
			log.Info().Msg("All filtering disabled")
		}
		packageParse = bacnetanalyzer.PackageParse
		serializePackage = bacnetanalyzer.SerializePackage
	case "c-bus":
		if !config.AnalyzeConfigInstance.NoFilter {
			if config.AnalyzeConfigInstance.Filter == "" && config.CBusConfigInstance.CBusFilter != "" {
				log.Debug().Str("filter", config.CBusConfigInstance.Filter).Msg("Setting cbus filter")
				filterExpression = config.CBusConfigInstance.CBusFilter
			}
		} else {
			log.Info().Msg("All filtering disabled")
		}
		analyzer := cbusanalyzer.Analyzer{Client: net.ParseIP(config.AnalyzeConfigInstance.Client)}
		analyzer.Init()
		packageParse = analyzer.PackageParse
		serializePackage = analyzer.SerializePackage
		mapPackets = analyzer.MapPackets
		if !config.AnalyzeConfigInstance.NoCustomMapping {
			byteOutput = analyzer.ByteOutput
		} else {
			log.Info().Msg("Custom mapping disabled")
		}
	default:
		return errors.Errorf("Unsupported protocol type %s", protocolType)
	}

	log.Info().
		Str("pcapFile", pcapFile).
		Str("protocolType", protocolType).
		Str("filterExpression", filterExpression).
		Msg("Analyzing pcap file pcapFile with protocolType and filter filterExpression now")
	handle, numberOfPackage, timestampToIndexMap, err := pcaphandler.GetIndexedPcapHandle(pcapFile, filterExpression)
	if err != nil {
		return errors.Wrap(err, "Error getting handle")
	}
	log.Info().Int("numberOfPackage", numberOfPackage).Msg("Starting to analyze numberOfPackage packages")
	defer handle.Close()
	log.Debug().Interface("handle", handle).Int("numberOfPackage", numberOfPackage).Msg("got handle")
	source := pcaphandler.GetPacketSource(handle)
	bar := progressbar.NewOptions(numberOfPackage, progressbar.OptionSetWriter(ansi.NewAnsiStderr()),
		progressbar.OptionSetVisibility(!config.RootConfigInstance.HideProgressBar),
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
	currentPackageNum := uint(0)
	parseFails := 0
	serializeFails := 0
	compareFails := 0
	for packet := range mapPackets(source.Packets(), func(packet gopacket.Packet) common.PacketInformation {
		return createPacketInformation(pcapFile, packet, timestampToIndexMap)
	}) {
		if err := ctx.Err(); err != nil {
			log.Info().Err(err).Uint("currentPackageNum", currentPackageNum).Msg("Aborted after currentPackageNum packages")
			break
		}
		currentPackageNum++
		if currentPackageNum < config.AnalyzeConfigInstance.StartPackageNumber {
			log.Debug().
				Uint("currentPackageNum", currentPackageNum).
				Uint("startPackageNum", config.AnalyzeConfigInstance.StartPackageNumber).
				Msg("Skipping package number currentPackageNum (till no. startPackageNum)")
			continue
		}
		if currentPackageNum > config.AnalyzeConfigInstance.PackageNumberLimit {
			log.Warn().
				Uint("PackageNumberLimit", config.AnalyzeConfigInstance.PackageNumberLimit).
				Msg("Aborting reading packages because we hit the limit of packageNumberLimit")
			break
		}
		if packet == nil {
			log.Debug().Msg("Done reading packages. (nil returned)")
			break
		}
		if err := bar.Add(1); err != nil {
			log.Warn().Err(err).Msg("Error updating progressBar")
		}
		packetInformation := createPacketInformation(pcapFile, packet, timestampToIndexMap)
		realPacketNumber := packetInformation.PacketNumber
		if filteredPackage, ok := packet.(common.FilteredPackage); ok {
			log.Info().Err(filteredPackage.FilterReason()).
				Int("realPacketNumber", realPacketNumber).Msg("No.[realPacketNumber] was filtered")
			continue
		}

		applicationLayer := packet.ApplicationLayer()
		if applicationLayer == nil {
			log.Info().Stringer("packetInformation", packetInformation).
				Int("realPacketNumber", realPacketNumber).
				Msg("No.[realPacketNumber] No application layer")
			continue
		}
		payload := applicationLayer.Payload()
		if parsed, err := packageParse(packetInformation, payload); err != nil {
			switch err {
			case common.ErrUnterminatedPackage:
				log.Info().Stringer("packetInformation", packetInformation).
					Int("realPacketNumber", realPacketNumber).
					Msg("No.[realPacketNumber] is unterminated")
			case common.ErrEmptyPackage:
				log.Info().Stringer("packetInformation", packetInformation).
					Int("realPacketNumber", realPacketNumber).
					Msg("No.[realPacketNumber] is empty")
			case common.ErrEcho:
				log.Info().Stringer("packetInformation", packetInformation).
					Int("realPacketNumber", realPacketNumber).
					Msg("No.[realPacketNumber] is echo")
			default:
				parseFails++
				// TODO: write report to xml or something
				log.Error().
					Err(err).
					Stringer("packetInformation", packetInformation).
					Int("realPacketNumber", realPacketNumber).
					Str("byteOutput", byteOutput(payload)).
					Msg("No.[realPacketNumber] Error parsing package.")
			}
			continue
		} else {
			if messageCallback != nil {
				messageCallback(parsed)
			}
			log.Info().
				Stringer("packetInformation", packetInformation).
				Int("realPacketNumber", realPacketNumber).
				Msg("No.[realPacketNumber] Parsed")
			if config.AnalyzeConfigInstance.Verbosity > 1 {
				prettyPrint(parsed)
			}
			if config.AnalyzeConfigInstance.OnlyParse {
				log.Trace().Msg("only parsing")
				continue
			}
			serializedBytes, err := serializePackage(parsed)
			if err != nil {
				serializeFails++
				// TODO: write report to xml or something
				log.Warn().
					Err(err).
					Stringer("packetInformation", packetInformation).
					Int("realPacketNumber", realPacketNumber).
					Msg("No.[realPacketNumber] Error serializing")
				continue
			}
			if config.AnalyzeConfigInstance.NoBytesCompare {
				log.Trace().Msg("not comparing bytes")
				continue
			}
			if compareResult := bytes.Compare(payload, serializedBytes); compareResult != 0 {
				compareFails++
				// TODO: write report to xml or something
				log.Warn().
					Stringer("packetInformation", packetInformation).
					Int("realPacketNumber", realPacketNumber).
					Str("byteOutputPayload", byteOutput(payload)).
					Str("byteSerializedBytes", byteOutput(serializedBytes)).
					Msg("No.[realPacketNumber] Bytes don't match.")
				if config.AnalyzeConfigInstance.Verbosity > 0 {
					_, _ = fmt.Fprintf(stdout, "Original bytes\n%s\n%s\n", hex.Dump(payload), hex.Dump(serializedBytes))
				}
			}
		}
	}

	log.Info().
		Uint("currentPackageNum", currentPackageNum).
		Int("numberOfPackage", numberOfPackage).
		Int("parseFails", parseFails).
		Int("serializeFails", serializeFails).
		Int("compareFails", compareFails).
		Msg("Done evaluating currentPackageNum of numberOfPackage packages (parseFails failed to parse, serializeFails failed to serialize and compareFails failed in byte comparison)")
	return nil
}

func createPacketInformation(pcapFile string, packet gopacket.Packet, timestampToIndexMap map[time.Time]int) common.PacketInformation {
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
	return packetInformation
}
