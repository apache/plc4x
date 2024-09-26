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

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/protocols/opcua/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

type OpcuaParserHelper struct {
}

func (m OpcuaParserHelper) Parse(typeName string, arguments []string, io utils.ReadBuffer) (any, error) {
	switch typeName {
	case "LocaleId":
		return LocaleIdParseWithBuffer(context.Background(), io)
	case "ImageGIF":
		return ImageGIFParseWithBuffer(context.Background(), io)
	case "EncodedTicket":
		return EncodedTicketParseWithBuffer(context.Background(), io)
	case "OpenChannelMessage":
		response, err := utils.StrToBool(arguments[0])
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing")
		}
		return OpenChannelMessageParseWithBuffer[OpenChannelMessage](context.Background(), io, response)
	case "ImageJPG":
		return ImageJPGParseWithBuffer(context.Background(), io)
	case "PascalByteString":
		return PascalByteStringParseWithBuffer(context.Background(), io)
	case "DiagnosticInfo":
		return DiagnosticInfoParseWithBuffer(context.Background(), io)
	case "PascalString":
		return PascalStringParseWithBuffer(context.Background(), io)
	case "TwoByteNodeId":
		return TwoByteNodeIdParseWithBuffer(context.Background(), io)
	case "OpcuaAPU":
		response, err := utils.StrToBool(arguments[0])
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing")
		}
		return OpcuaAPUParseWithBuffer(context.Background(), io, response)
	case "Index":
		return IndexParseWithBuffer(context.Background(), io)
	case "StatusCode":
		return StatusCodeParseWithBuffer(context.Background(), io)
	case "NormalizedString":
		return NormalizedStringParseWithBuffer(context.Background(), io)
	case "QualifiedName":
		return QualifiedNameParseWithBuffer(context.Background(), io)
	case "NumericNodeId":
		return NumericNodeIdParseWithBuffer(context.Background(), io)
	case "FourByteNodeId":
		return FourByteNodeIdParseWithBuffer(context.Background(), io)
	case "AudioDataType":
		return AudioDataTypeParseWithBuffer(context.Background(), io)
	case "SecurityHeader":
		return SecurityHeaderParseWithBuffer(context.Background(), io)
	case "ContinuationPoint":
		return ContinuationPointParseWithBuffer(context.Background(), io)
	case "Variant":
		return VariantParseWithBuffer[Variant](context.Background(), io)
	case "Payload":
		extensible, err := utils.StrToBool(arguments[0])
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing")
		}
		byteCount, err := utils.StrToUint32(arguments[1])
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing")
		}
		return PayloadParseWithBuffer[Payload](context.Background(), io, extensible, byteCount)
	case "ExtensionObjectEncodingMask":
		return ExtensionObjectEncodingMaskParseWithBuffer(context.Background(), io)
	case "DurationString":
		return DurationStringParseWithBuffer(context.Background(), io)
	case "Structure":
		return StructureParseWithBuffer(context.Background(), io)
	case "OpcuaConstants":
		return OpcuaConstantsParseWithBuffer(context.Background(), io)
	case "UtcTime":
		return UtcTimeParseWithBuffer(context.Background(), io)
	case "MessagePDU":
		response, err := utils.StrToBool(arguments[0])
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing")
		}
		return MessagePDUParseWithBuffer[MessagePDU](context.Background(), io, response)
	case "Counter":
		return CounterParseWithBuffer(context.Background(), io)
	case "SequenceHeader":
		return SequenceHeaderParseWithBuffer(context.Background(), io)
	case "NodeId":
		return NodeIdParseWithBuffer(context.Background(), io)
	case "RsaEncryptedSecret":
		return RsaEncryptedSecretParseWithBuffer(context.Background(), io)
	case "ExtensionObject":
		includeEncodingMask, err := utils.StrToBool(arguments[0])
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing")
		}
		return ExtensionObjectParseWithBuffer[ExtensionObject](context.Background(), io, includeEncodingMask)
	case "LocalizedText":
		return LocalizedTextParseWithBuffer(context.Background(), io)
	case "IntegerId":
		return IntegerIdParseWithBuffer(context.Background(), io)
	case "ByteStringArray":
		return ByteStringArrayParseWithBuffer(context.Background(), io)
	case "Handle":
		return HandleParseWithBuffer(context.Background(), io)
	case "ImagePNG":
		return ImagePNGParseWithBuffer(context.Background(), io)
	case "XmlElement":
		return XmlElementParseWithBuffer(context.Background(), io)
	case "SessionAuthenticationToken":
		return SessionAuthenticationTokenParseWithBuffer(context.Background(), io)
	case "DataValue":
		return DataValueParseWithBuffer(context.Background(), io)
	case "GuidNodeId":
		return GuidNodeIdParseWithBuffer(context.Background(), io)
	case "GuidValue":
		return GuidValueParseWithBuffer(context.Background(), io)
	case "TrimmedString":
		return TrimmedStringParseWithBuffer(context.Background(), io)
	case "ApplicationInstanceCertificate":
		return ApplicationInstanceCertificateParseWithBuffer(context.Background(), io)
	case "BitFieldMaskDataType":
		return BitFieldMaskDataTypeParseWithBuffer(context.Background(), io)
	case "ImageBMP":
		return ImageBMPParseWithBuffer(context.Background(), io)
	case "ExtensionObjectDefinition":
		extensionId, err := utils.StrToInt32(arguments[0])
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing")
		}
		return ExtensionObjectDefinitionParseWithBuffer[ExtensionObjectDefinition](context.Background(), io, extensionId)
	case "ExpandedNodeId":
		return ExpandedNodeIdParseWithBuffer(context.Background(), io)
	case "OpcuaProtocolLimits":
		return OpcuaProtocolLimitsParseWithBuffer(context.Background(), io)
	case "NumericRange":
		return NumericRangeParseWithBuffer(context.Background(), io)
	case "SemanticVersionString":
		return SemanticVersionStringParseWithBuffer(context.Background(), io)
	case "ByteStringNodeId":
		return ByteStringNodeIdParseWithBuffer(context.Background(), io)
	case "TimeString":
		return TimeStringParseWithBuffer(context.Background(), io)
	case "EccEncryptedSecret":
		return EccEncryptedSecretParseWithBuffer(context.Background(), io)
	case "StringNodeId":
		return StringNodeIdParseWithBuffer(context.Background(), io)
	case "VersionTime":
		return VersionTimeParseWithBuffer(context.Background(), io)
	case "UriString":
		return UriStringParseWithBuffer(context.Background(), io)
	case "DecimalString":
		return DecimalStringParseWithBuffer(context.Background(), io)
	case "NodeIdTypeDefinition":
		return NodeIdTypeDefinitionParseWithBuffer[NodeIdTypeDefinition](context.Background(), io)
	case "DateString":
		return DateStringParseWithBuffer(context.Background(), io)
	}
	return nil, errors.Errorf("Unsupported type %s", typeName)
}
