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

package opcua

import (
	"strconv"
	"time"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/opcua/readwrite/model"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func generateNodeId(tag Tag) (readWriteModel.NodeId, error) {
	var nodeId readWriteModel.NodeId
	if tag.GetIdentifierType() == readWriteModel.OpcuaIdentifierType_BINARY_IDENTIFIER {
		parsedIdentifier, err := strconv.ParseUint(tag.GetIdentifier(), 10, 8)
		if err != nil {
			return nil, errors.New("Error parsing identifier")
		}
		nodeId = readWriteModel.NewNodeId(readWriteModel.NewNodeIdTwoByte(uint8(parsedIdentifier)))
	} else if tag.GetIdentifierType() == readWriteModel.OpcuaIdentifierType_NUMBER_IDENTIFIER {
		parsedIdentifier, err := strconv.ParseUint(tag.GetIdentifier(), 10, 32)
		if err != nil {
			return nil, errors.New("Error parsing identifier")
		}
		nodeId = readWriteModel.NewNodeId(readWriteModel.NewNodeIdNumeric( /*TODO: do we want to check for overflow?*/ uint16(tag.GetNamespace()), uint32(parsedIdentifier)))
	} else if tag.GetIdentifierType() == readWriteModel.OpcuaIdentifierType_GUID_IDENTIFIER {
		guid, err := uuid.Parse(tag.GetIdentifier())
		if err != nil {
			return nil, errors.Wrap(err, "error parsing guid")
		}
		guidBytes, err := guid.MarshalBinary() // TODO: do we need to do flip it here?
		if err != nil {
			return nil, errors.Wrap(err, "error marshaling guid")
		}
		nodeId = readWriteModel.NewNodeId(readWriteModel.NewNodeIdGuid( /*TODO: do we want to check for overflow?*/ uint16(tag.GetNamespace()), guidBytes))
	} else if tag.GetIdentifierType() == readWriteModel.OpcuaIdentifierType_STRING_IDENTIFIER {
		nodeId = readWriteModel.NewNodeId(readWriteModel.NewNodeIdString( /*TODO: do we want to check for overflow?*/ uint16(tag.GetNamespace()), readWriteModel.NewPascalString(tag.GetIdentifier())))
	}
	return nodeId, nil
}

func getDateTime(dateTime int64) time.Time {
	return time.UnixMilli((dateTime - EPOCH_OFFSET) / 10000)
}
