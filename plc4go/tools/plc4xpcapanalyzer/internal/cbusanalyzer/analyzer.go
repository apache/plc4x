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

package cbusanalyzer

import (
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net"
)

type Analyzer struct {
	Client net.IP
}

func (a Analyzer) PackageParse(packetInformation common.PacketInformation, payload []byte) (interface{}, error) {
	log.Debug().Msgf("Parsing %s", packetInformation)
	// TODO: srcchk we need to pull that out of the config
	isResponse := packetInformation.DstIp.Equal(a.Client)
	log.Debug().Stringer("packetInformation", packetInformation).Msgf("isResponse: %t", isResponse)
	parse, err := model.CBusMessageParse(utils.NewReadBufferByteBased(payload), isResponse, true)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing CBusCommand")
	}
	log.Debug().Msgf("Parsed c-bus command %s", parse)
	return parse, nil
}

func (a Analyzer) SerializePackage(command interface{}) ([]byte, error) {
	if command, ok := command.(model.CBusMessage); !ok {
		log.Fatal().Msgf("Unsupported type %T supplied", command)
		panic("unreachable statement")
	} else {
		based := utils.NewWriteBufferByteBased()
		if err := command.Serialize(based); err != nil {
			return nil, errors.Wrap(err, "Error serializing")
		}
		return based.GetBytes(), nil
	}
}
