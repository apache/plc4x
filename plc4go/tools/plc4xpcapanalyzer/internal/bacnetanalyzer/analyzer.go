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

package bacnetanalyzer

import (
	"context"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func PackageParse(packetInformation common.PacketInformation, payload []byte) (spi.Message, error) {
	log.Debug().Stringer("packetInformation", packetInformation).Msg("Parsing")
	parse, err := model.BVLCParse(context.TODO(), payload)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing bvlc")
	}
	log.Debug().Stringer("parse", parse).Msg("Parsed bvlc")
	return parse, nil
}

func SerializePackage(bvlc spi.Message) ([]byte, error) {
	if bvlc, ok := bvlc.(model.BVLC); !ok {
		log.Fatal().Type("bvlc", bvlc).Msg("Unsupported type supplied")
		panic("unreachable statement")
	} else {
		theBytes, err := bvlc.Serialize()
		if err != nil {
			return nil, errors.Wrap(err, "Error serializing")
		}
		return theBytes, nil
	}
}
