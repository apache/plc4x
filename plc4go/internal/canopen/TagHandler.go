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

package canopen

import (
	"regexp"
	"strconv"
	"strings"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/canopen/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
)

// maxNodeID is the highest CANopen node id: the node is the low 7 bits of an 11 bit COB-ID.
const maxNodeID = 127

var (
	sdoPattern       = regexp.MustCompile(`^SDO:(?P<node>\d+)(?:/(?P<answerNode>\d+))?:(?P<index>0[xX][0-9a-fA-F]+|\d+)/(?P<subIndex>0[xX][0-9a-fA-F]+|\d+):(?P<dataType>\w+)$`)
	pdoPattern       = regexp.MustCompile(`^(?P<pdo>(?:RECEIVE|TRANSMIT)_PDO_[1-4]):(?P<node>\d+):(?P<dataType>\w+)$`)
	nmtPattern       = regexp.MustCompile(`^NMT(?::(?P<node>\d+))?$`)
	heartbeatPattern = regexp.MustCompile(`^HEARTBEAT(?::(?P<node>\d+))?$`)
)

type TagHandler struct{}

func NewTagHandler() TagHandler {
	return TagHandler{}
}

func (m TagHandler) ParseTag(tagAddress string) (apiModel.PlcTag, error) {
	address := strings.TrimSpace(tagAddress)
	switch {
	case sdoPattern.MatchString(address):
		return parseSDO(address)
	case pdoPattern.MatchString(address):
		return parsePDO(address)
	case nmtPattern.MatchString(address):
		node, err := optionalNode(nmtPattern, address)
		if err != nil {
			return nil, err
		}
		return Tag{Kind: TagKindNMT, Node: node}, nil
	case heartbeatPattern.MatchString(address):
		node, err := optionalNode(heartbeatPattern, address)
		if err != nil {
			return nil, err
		}
		return Tag{Kind: TagKindHeartbeat, Node: node}, nil
	}
	return nil, errors.Errorf("unable to parse CANopen address %q: expected SDO:<node>:<index>/<subindex>:<type>, "+
		"<RECEIVE|TRANSMIT>_PDO_<n>:<node>:<type>, NMT[:<node>] or HEARTBEAT[:<node>]", tagAddress)
}

// ParseQuery is not supported: there is no standard way to list a node's object dictionary over the
// bus (that is what EDS files are for).
func (m TagHandler) ParseQuery(_ string) (apiModel.PlcQuery, error) {
	return nil, errors.New("the canopen driver doesn't support browsing")
}

func parseSDO(address string) (Tag, error) {
	groups := submatches(sdoPattern, address)
	node, err := parseNode(groups["node"], false)
	if err != nil {
		return Tag{}, err
	}
	var answerNode uint8
	if groups["answerNode"] != "" {
		if answerNode, err = parseNode(groups["answerNode"], false); err != nil {
			return Tag{}, err
		}
	}
	index, err := parseNumber(groups["index"], 16)
	if err != nil {
		return Tag{}, errors.Wrapf(err, "bad SDO index in %q", address)
	}
	subIndex, err := parseNumber(groups["subIndex"], 8)
	if err != nil {
		return Tag{}, errors.Wrapf(err, "bad SDO subindex in %q", address)
	}
	dataType, err := parseDataType(groups["dataType"])
	if err != nil {
		return Tag{}, err
	}
	return Tag{
		Kind:       TagKindSDO,
		Node:       node,
		AnswerNode: answerNode,
		Index:      uint16(index),
		SubIndex:   uint8(subIndex),
		DataType:   dataType,
	}, nil
}

func parsePDO(address string) (Tag, error) {
	groups := submatches(pdoPattern, address)
	node, err := parseNode(groups["node"], false)
	if err != nil {
		return Tag{}, err
	}
	service, ok := readWriteModel.CANOpenServiceByName(groups["pdo"])
	if !ok {
		return Tag{}, errors.Errorf("unknown PDO %s", groups["pdo"])
	}
	dataType, err := parseDataType(groups["dataType"])
	if err != nil {
		return Tag{}, err
	}
	if dataType.NumBits() > 64 {
		return Tag{}, errors.Errorf("%s does not fit the 8 bytes of a PDO", dataType)
	}
	return Tag{Kind: TagKindPDO, Node: node, Service: service, DataType: dataType}, nil
}

func optionalNode(pattern *regexp.Regexp, address string) (uint8, error) {
	node := submatches(pattern, address)["node"]
	if node == "" {
		return 0, nil
	}
	return parseNode(node, true)
}

// parseNode reads a node id. Zero is "every node", which only NMT and HEARTBEAT have a use for.
func parseNode(text string, allowZero bool) (uint8, error) {
	node, err := strconv.ParseUint(text, 10, 8)
	if err != nil || node > maxNodeID || (node == 0 && !allowZero) {
		if allowZero {
			return 0, errors.Errorf("CANopen node id %q is not in 0..%d", text, maxNodeID)
		}
		return 0, errors.Errorf("CANopen node id %q is not in 1..%d", text, maxNodeID)
	}
	return uint8(node), nil
}

func parseNumber(text string, bits int) (uint64, error) {
	if strings.HasPrefix(text, "0x") || strings.HasPrefix(text, "0X") {
		return strconv.ParseUint(text[2:], 16, bits)
	}
	return strconv.ParseUint(text, 10, bits)
}

func parseDataType(name string) (readWriteModel.CANOpenDataType, error) {
	dataType, ok := readWriteModel.CANOpenDataTypeByName(name)
	if !ok {
		return 0, errors.Errorf("unknown CANopen data type %s", name)
	}
	switch dataType {
	case readWriteModel.CANOpenDataType_TIME_OF_DAY, readWriteModel.CANOpenDataType_TIME_DIFFERENCE:
		// The data item spec has no case for these (plc4j has the same gap), so a value could be
		// neither decoded nor encoded.
		return 0, errors.Errorf("CANopen data type %s is not supported yet", name)
	}
	return dataType, nil
}

func submatches(pattern *regexp.Regexp, text string) map[string]string {
	match := pattern.FindStringSubmatch(text)
	groups := map[string]string{}
	for i, name := range pattern.SubexpNames() {
		if name != "" && i < len(match) {
			groups[name] = match[i]
		}
	}
	return groups
}
