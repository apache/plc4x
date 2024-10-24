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

package readwrite

import (
	"context"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/protocols/opcua/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

type OpcuaXmlParserHelper struct {
}

// Temporary imports to silent compiler warnings (TODO: migrate from static to emission based imports)
func init() {
	_ = strconv.ParseUint
	_ = strconv.ParseInt
	_ = strings.Join
	_ = utils.Dump
}

func (m OpcuaXmlParserHelper) Parse(typeName string, xmlString string, parserArguments ...string) (any, error) {
	switch typeName {
	case "LocaleId":
		return LocaleIdParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "ImageGIF":
		return ImageGIFParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "EncodedTicket":
		return EncodedTicketParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "OpenChannelMessage":
		response := parserArguments[0] == "true"
		return OpenChannelMessageParseWithBuffer[OpenChannelMessage](context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)), response)
	case "ImageJPG":
		return ImageJPGParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "PascalByteString":
		return PascalByteStringParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "DiagnosticInfo":
		return DiagnosticInfoParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "PascalString":
		return PascalStringParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "TwoByteNodeId":
		return TwoByteNodeIdParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "OpcuaAPU":
		response := parserArguments[0] == "true"
		return OpcuaAPUParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)), response)
	case "Index":
		return IndexParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "StatusCode":
		return StatusCodeParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "NormalizedString":
		return NormalizedStringParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "QualifiedName":
		return QualifiedNameParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "NumericNodeId":
		return NumericNodeIdParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "FourByteNodeId":
		return FourByteNodeIdParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "AudioDataType":
		return AudioDataTypeParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "SecurityHeader":
		return SecurityHeaderParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "ContinuationPoint":
		return ContinuationPointParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "Variant":
		return VariantParseWithBuffer[Variant](context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "Payload":
		extensible := parserArguments[0] == "true"
		parsedUint1, err := strconv.ParseUint(parserArguments[1], 10, 32)
		if err != nil {
			return nil, err
		}
		byteCount := uint32(parsedUint1)
		return PayloadParseWithBuffer[Payload](context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)), extensible, byteCount)
	case "ExtensionObjectEncodingMask":
		return ExtensionObjectEncodingMaskParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "DurationString":
		return DurationStringParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "Structure":
		return StructureParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "OpcuaConstants":
		return OpcuaConstantsParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "UtcTime":
		return UtcTimeParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "MessagePDU":
		response := parserArguments[0] == "true"
		return MessagePDUParseWithBuffer[MessagePDU](context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)), response)
	case "Counter":
		return CounterParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "SequenceHeader":
		return SequenceHeaderParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "NodeId":
		return NodeIdParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "RsaEncryptedSecret":
		return RsaEncryptedSecretParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "ExtensionObject":
		includeEncodingMask := parserArguments[0] == "true"
		return ExtensionObjectParseWithBuffer[ExtensionObject](context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)), includeEncodingMask)
	case "LocalizedText":
		return LocalizedTextParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "IntegerId":
		return IntegerIdParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "ByteStringArray":
		return ByteStringArrayParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "Handle":
		return HandleParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "ImagePNG":
		return ImagePNGParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "XmlElement":
		return XmlElementParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "SessionAuthenticationToken":
		return SessionAuthenticationTokenParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "DataValue":
		return DataValueParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "GuidNodeId":
		return GuidNodeIdParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "GuidValue":
		return GuidValueParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "TrimmedString":
		return TrimmedStringParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "ApplicationInstanceCertificate":
		return ApplicationInstanceCertificateParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "BitFieldMaskDataType":
		return BitFieldMaskDataTypeParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "ImageBMP":
		return ImageBMPParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "ExtensionObjectDefinition":
		parsedInt0, err := strconv.ParseInt(parserArguments[0], 10, 32)
		if err != nil {
			return nil, err
		}
		extensionId := int32(parsedInt0)
		return ExtensionObjectDefinitionParseWithBuffer[ExtensionObjectDefinition](context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)), extensionId)
	case "ExpandedNodeId":
		return ExpandedNodeIdParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "OpcuaProtocolLimits":
		return OpcuaProtocolLimitsParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "NumericRange":
		return NumericRangeParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "SemanticVersionString":
		return SemanticVersionStringParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "ByteStringNodeId":
		return ByteStringNodeIdParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "TimeString":
		return TimeStringParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "EccEncryptedSecret":
		return EccEncryptedSecretParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "StringNodeId":
		return StringNodeIdParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "VersionTime":
		return VersionTimeParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "UriString":
		return UriStringParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "DecimalString":
		return DecimalStringParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "NodeIdTypeDefinition":
		return NodeIdTypeDefinitionParseWithBuffer[NodeIdTypeDefinition](context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "DateString":
		return DateStringParseWithBuffer(context.Background(), utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	}
	return nil, errors.Errorf("Unsupported type %s", typeName)
}
