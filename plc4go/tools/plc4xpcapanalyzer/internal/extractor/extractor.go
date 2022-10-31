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

package extractor

import (
	"context"
	"fmt"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/config"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/common"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/pcaphandler"
	"github.com/fatih/color"
	"github.com/google/gopacket/layers"
	"github.com/k0kubun/go-ansi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"io"
	"net"
)

func Extract(pcapFile, protocolType string) error {
	return ExtractWithOutput(context.TODO(), pcapFile, protocolType, ansi.NewAnsiStdout(), ansi.NewAnsiStderr())
}

func ExtractWithOutput(ctx context.Context, pcapFile, protocolType string, stdout, stderr io.Writer) error {
	log.Info().Msgf("Analyzing pcap file '%s' with protocolType '%s' and filter '%s' now", pcapFile, protocolType, config.ExtractConfigInstance.Filter)

	handle, numberOfPackage, timestampToIndexMap, err := pcaphandler.GetIndexedPcapHandle(pcapFile, config.ExtractConfigInstance.Filter)
	if err != nil {
		return errors.Wrap(err, "Error getting handle")
	}
	log.Info().Msgf("Starting to analyze %d packages", numberOfPackage)
	defer handle.Close()
	log.Debug().Interface("handle", handle).Int("numberOfPackage", numberOfPackage).Msg("got handle")
	source := pcaphandler.GetPacketSource(handle)
	var printPayload = func(packetInformation common.PacketInformation, item []byte) {
		_, _ = fmt.Fprintf(stdout, "%x\n", item)
	}
	switch protocolType {
	case "bacnet":
		// nothing special as this is byte based
	case "c-bus":
		// c-bus is string based so we consume the string and print it
		clientIp := net.ParseIP(config.ExtractConfigInstance.Client)
		serverResponseWriter := color.New(color.FgRed)
		serverResponseIndicatorWriter := color.New(color.FgHiRed)
		clientRequestWriter := color.New(color.FgGreen)
		clientRequestIndicatorWriter := color.New(color.FgHiGreen)
		printPayload = func(packetInformation common.PacketInformation, payload []byte) {
			payloadString := ""
			suffix := ""
			extraInformation := ""
			if config.ExtractConfigInstance.Verbosity > 2 {
				extraInformation = fmt.Sprintf("(No.[%d])", packetInformation.PacketNumber)
			}
			if len(payload) > 0 {
				if properTerminated := payload[len(payload)-1] == 0x0D || payload[len(payload)-1] == 0x0A; properTerminated {
					suffix = "\n"
				}
				quotedPayload := fmt.Sprintf("%+q", payload)
				unquotedPayload := quotedPayload[1 : len(quotedPayload)-1]
				payloadString = unquotedPayload
			}
			if isResponse := packetInformation.DstIp.Equal(clientIp); isResponse {
				if config.ExtractConfigInstance.ShowDirectionalIndicators {
					_, _ = serverResponseIndicatorWriter.Fprintf(stderr, "%s(<--pci)", extraInformation)
				}
				_, _ = serverResponseWriter.Fprintf(stdout, "%s%s", payloadString, suffix)
			} else {
				if config.ExtractConfigInstance.ShowDirectionalIndicators {
					_, _ = clientRequestIndicatorWriter.Fprintf(stderr, "%s(-->pci)", extraInformation)
				}
				_, _ = clientRequestWriter.Fprintf(stdout, "%s%s", payloadString, suffix)
			}
		}
	}
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
	for packet := range source.Packets() {
		if ctx.Err() == context.Canceled {
			log.Info().Msgf("Aborted after %d packages", currentPackageNum)
			break
		}
		currentPackageNum++
		if currentPackageNum < config.ExtractConfigInstance.StartPackageNumber {
			log.Debug().Msgf("Skipping package number %d (till no. %d)", currentPackageNum, config.ExtractConfigInstance.StartPackageNumber)
			continue
		}
		if currentPackageNum > config.ExtractConfigInstance.PackageNumberLimit {
			log.Warn().Msgf("Aborting reading packages because we hit the limit of %d", config.ExtractConfigInstance.PackageNumberLimit)
			break
		}
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

		var payload []byte
		applicationLayer := packet.ApplicationLayer()
		if applicationLayer == nil {
			log.Info().Stringer("packetInformation", packetInformation).Msgf("No.[%d] No application layer", realPacketNumber)
		} else {
			payload = applicationLayer.Payload()
		}

		log.Debug().Msgf("Got payload %x", payload)
		if config.ExtractConfigInstance.Verbosity > 1 {
			printPayload(packetInformation, payload)
		}
	}
	_, _ = fmt.Fprintf(stdout, "\n")

	log.Info().Msgf("Done evaluating %d of %d packages (%d failed to parse, %d failed to serialize and %d failed in byte comparison)", currentPackageNum, numberOfPackage, parseFails, serializeFails, compareFails)
	return nil
}
