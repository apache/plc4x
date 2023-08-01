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
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/binary"
	"math"
	"math/rand"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/opcua/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/dchest/uniuri"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/exp/slices"
)

const (
	FINAL_CHUNK                   = "F"
	CONTINUATION_CHUNK            = "C"
	ABORT_CHUNK                   = "A"
	VERSION                       = 0
	DEFAULT_MAX_CHUNK_COUNT       = 64
	DEFAULT_MAX_MESSAGE_SIZE      = 2097152
	DEFAULT_RECEIVE_BUFFER_SIZE   = 65535
	DEFAULT_SEND_BUFFER_SIZE      = 65535
	REQUEST_TIMEOUT               = 10000 * time.Millisecond
	REQUEST_TIMEOUT_LONG          = 10000
	PASSWORD_ENCRYPTION_ALGORITHM = "http://www.w3.org/2001/04/xmlenc#rsa-oaep"
	EPOCH_OFFSET                  = 116444736000000000 //Offset between OPC UA epoch time and linux epoch time.
)

var (
	SECURITY_POLICY_NONE = readWriteModel.NewPascalString("http://opcfoundation.org/UA/SecurityPolicy#None")
	NULL_STRING          = readWriteModel.NewPascalString("")
	NULL_BYTE_STRING     = readWriteModel.NewPascalByteString(-1, nil)
	NULL_EXPANDED_NODEID = readWriteModel.NewExpandedNodeId(false,
		false,
		readWriteModel.NewNodeIdTwoByte(0),
		nil,
		nil,
	)
	NULL_EXTENSION_OBJECT = readWriteModel.NewExtensionObject(NULL_EXPANDED_NODEID,
		readWriteModel.NewExtensionObjectEncodingMask(false, false, false),
		readWriteModel.NewNullExtension(),
		false) // Body

	INET_ADDRESS_PATTERN = regexp.MustCompile(`(.(?P<transportCode>tcp))?://` +
		`(?P<transportHost>[\\w.-]+)(:` +
		`(?P<transportPort>\\d*))?`)

	URI_PATTERN = regexp.MustCompile(`^(?P<protocolCode>opc)` +
		INET_ADDRESS_PATTERN.String() +
		`(?P<transportEndpoint>[\\w/=]*)[\\?]?`,
	)
	APPLICATION_URI             = readWriteModel.NewPascalString("urn:apache:plc4x:client")
	PRODUCT_URI                 = readWriteModel.NewPascalString("urn:apache:plc4x:client")
	APPLICATION_TEXT            = readWriteModel.NewPascalString("OPCUA client for the Apache PLC4X:PLC4J project")
	DEFAULT_CONNECTION_LIFETIME = uint32(36000000)
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=SecureChannel
type SecureChannel struct {
	sessionName               string
	clientNonce               []byte
	requestHandleGenerator    atomic.Uint32
	policyId                  readWriteModel.PascalString
	tokenType                 readWriteModel.UserTokenType `stringer:"true"`
	discovery                 bool
	certFile                  string
	keyStoreFile              string
	ckp                       CertificateKeyPair
	endpoint                  readWriteModel.PascalString
	username                  string
	password                  string
	securityPolicy            string
	publicCertificate         readWriteModel.PascalByteString
	thumbprint                readWriteModel.PascalByteString
	isEncrypted               bool
	senderCertificate         []byte
	senderNonce               []byte
	certificateThumbprint     readWriteModel.PascalByteString
	checkedEndpoints          bool
	encryptionHandler         *EncryptionHandler
	configuration             Configuration `stringer:"true"`
	channelId                 atomic.Int32
	tokenId                   atomic.Int32
	authenticationToken       readWriteModel.NodeIdTypeDefinition
	codec                     *MessageCodec
	channelTransactionManager *SecureChannelTransactionManager
	lifetime                  uint32
	keepAliveStateChange      sync.Mutex
	keepAliveIndicator        atomic.Bool
	keepAliveWg               sync.WaitGroup
	sendBufferSize            int
	maxMessageSize            int
	endpoints                 []string
	senderSequenceNumber      atomic.Int32

	log zerolog.Logger `ignore:"true"`
}

func NewSecureChannel(log zerolog.Logger, ctx DriverContext, configuration Configuration) *SecureChannel {
	s := &SecureChannel{
		configuration:             configuration,
		endpoint:                  readWriteModel.NewPascalString(configuration.endpoint),
		username:                  configuration.username,
		password:                  configuration.password,
		securityPolicy:            "http://opcfoundation.org/UA/SecurityPolicy#" + configuration.securityPolicy,
		sessionName:               "UaSession:" + APPLICATION_TEXT.GetStringValue() + ":" + uniuri.NewLen(20),
		clientNonce:               []byte(uniuri.NewLen(40)),
		keyStoreFile:              configuration.keyStoreFile,
		channelTransactionManager: NewSecureChannelTransactionManager(log),
		lifetime:                  DEFAULT_CONNECTION_LIFETIME,
		log:                       log,
	}
	s.requestHandleGenerator.Store(1)
	s.channelId.Store(1)
	s.tokenId.Store(1)
	ckp := configuration.ckp
	if configuration.securityPolicy == "Basic256Sha256" {
		//Sender Certificate gets populated during the 'discover' phase when encryption is enabled.
		s.senderCertificate = configuration.senderCertificate
		s.encryptionHandler = NewEncryptionHandler(s.log, ckp, s.senderCertificate, configuration.securityPolicy)
		certificate := ckp.getCertificate()
		s.publicCertificate = readWriteModel.NewPascalByteString(int32(len(certificate.Raw)), certificate.Raw)
		s.isEncrypted = true

		s.thumbprint = configuration.thumbprint
	} else {
		s.encryptionHandler = NewEncryptionHandler(s.log, ckp, s.senderCertificate, configuration.securityPolicy)
		s.publicCertificate = NULL_BYTE_STRING
		s.thumbprint = NULL_BYTE_STRING
		s.isEncrypted = false
	}

	// Generate a list of endpoints we can use.
	if address, err := url.Parse("none:" + configuration.host); err != nil {
		if names, err := net.LookupAddr(address.Host); err != nil {
			s.endpoints = append(s.endpoints, names[rand.Intn(len(names))])
		}
		s.endpoints = append(s.endpoints, address.Host)
		//s.endpoints = append(s.endpoints, address.Host)//TODO: not sure if golang can do
	} else {
		s.log.Warn().Msg("Unable to resolve host name. Using original host from connection string which may cause issues connecting to server")
		s.endpoints = append(s.endpoints, address.Host)
	}

	s.channelId.Store(1)
	return s
}

func (s *SecureChannel) submit(ctx context.Context, codec *MessageCodec, errorDispatcher func(err error), consumer func(opcuaResponse []byte), buffer utils.WriteBufferByteBased) {
	transactionId := s.channelTransactionManager.getTransactionIdentifier()

	//TODO: We need to split large messages up into chunks if it is larger than the sendBufferSize
	//      This value is negotiated when opening a channel

	messageRequest := readWriteModel.NewOpcuaMessageRequest(FINAL_CHUNK,
		s.channelId.Load(),
		s.tokenId.Load(),
		transactionId,
		transactionId,
		buffer.GetBytes())

	var apu readWriteModel.OpcuaAPU
	if s.isEncrypted {
		message, err := s.encryptionHandler.encodeMessage(ctx, messageRequest, buffer.GetBytes())
		if err != nil {
			errorDispatcher(err)
			return
		}
		apu, err = readWriteModel.OpcuaAPUParse(ctx, message, false)
		if err != nil {
			errorDispatcher(err)
			return
		}
	} else {
		apu = readWriteModel.NewOpcuaAPU(messageRequest, false)
	}

	requestConsumer := func(transactionId int32) {
		var messageBuffer []byte
		if err := codec.SendRequest(ctx, apu,
			func(message spi.Message) bool {
				opcuaAPU, ok := message.(readWriteModel.OpcuaAPUExactly)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				if decodedOpcuaAPU, err := s.encryptionHandler.decodeMessage(ctx, opcuaAPU); err != nil {
					s.log.Debug().Err(err).Msg("error decoding")
					return false
				} else {
					opcuaAPU = decodedOpcuaAPU.(readWriteModel.OpcuaAPUExactly)
				}
				messagePDU := opcuaAPU.GetMessage()
				opcuaResponse, ok := messagePDU.(readWriteModel.OpcuaMessageResponseExactly)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				if requestId := opcuaResponse.GetRequestId(); requestId != transactionId {
					s.log.Debug().Int32("requestId", requestId).Int32("transactionId", transactionId).Msg("Not relevant")
					return false
				} else {
					messageBuffer = opcuaResponse.GetMessage()
					if !(s.senderSequenceNumber.Add(1) == (opcuaResponse.GetSequenceNumber())) {
						s.log.Error().Msgf("Sequence number isn't as expected, we might have missed a packet. - %d != %d", s.senderSequenceNumber.Add(1), opcuaResponse.GetSequenceNumber())
						codec.fireDisconnected()
					}
				}
				return true
			},
			func(message spi.Message) error {
				opcuaAPU := message.(readWriteModel.OpcuaAPU)
				opcuaAPU, _ = s.encryptionHandler.decodeMessage(ctx, opcuaAPU)
				messagePDU := opcuaAPU.GetMessage()
				opcuaResponse := messagePDU.(readWriteModel.OpcuaMessageResponse)
				if opcuaResponse.GetChunk() == (FINAL_CHUNK) {
					s.tokenId.Store(opcuaResponse.GetSecureTokenId())
					s.channelId.Store(opcuaResponse.GetSecureChannelId())

					consumer(messageBuffer)
				}
				return nil
			},
			func(err error) error {
				errorDispatcher(err)
				return nil
			},
			REQUEST_TIMEOUT); err != nil {
			errorDispatcher(err)
		}
	}

	s.log.Debug().Msgf("Submitting Transaction to TransactionManager %v", transactionId)
	if err := s.channelTransactionManager.submit(requestConsumer, transactionId); err != nil {
		s.log.Debug().Err(err).Msg("error submitting")
	}
}

func (s *SecureChannel) onConnect(ctx context.Context, codec *MessageCodec) {
	s.log.Trace().Msg("on connect")
	// Only the TCP transport supports login.
	s.log.Debug().Msg("Opcua Driver running in ACTIVE mode.")
	s.codec = codec

	hello := readWriteModel.NewOpcuaHelloRequest(FINAL_CHUNK,
		VERSION,
		DEFAULT_RECEIVE_BUFFER_SIZE,
		DEFAULT_SEND_BUFFER_SIZE,
		DEFAULT_MAX_MESSAGE_SIZE,
		DEFAULT_MAX_CHUNK_COUNT,
		s.endpoint)

	requestConsumer := func(transactionId int32) {
		s.log.Trace().Int32("transactionId", transactionId).Msg("request consumer called")
		if err := codec.SendRequest(
			ctx,
			hello,
			func(message spi.Message) bool {
				opcuaAPU, ok := message.(readWriteModel.OpcuaAPUExactly)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				messagePDU := opcuaAPU.GetMessage()
				_, ok = messagePDU.(readWriteModel.OpcuaAcknowledgeResponseExactly)
				if !ok {
					s.log.Debug().Type("type", messagePDU).Msg("Not relevant")
					return false
				}
				return true
			},
			func(message spi.Message) error {
				opcuaAPU := message.(readWriteModel.OpcuaAPU)
				messagePDU := opcuaAPU.GetMessage()
				opcuaAcknowledgeResponse := messagePDU.(readWriteModel.OpcuaAcknowledgeResponse)
				s.onConnectOpenSecureChannel(ctx, codec, opcuaAcknowledgeResponse)
				return nil
			},
			func(err error) error {
				s.log.Debug().Err(err).Msg("error submitting")
				return nil
			},
			REQUEST_TIMEOUT); err != nil {
			s.log.Debug().Err(err).Msg("error sending")
		}
	}
	if err := s.channelTransactionManager.submit(requestConsumer, s.channelTransactionManager.getTransactionIdentifier()); err != nil {
		s.log.Debug().Err(err).Msg("error submitting")
	}
}

func (s *SecureChannel) onConnectOpenSecureChannel(ctx context.Context, codec *MessageCodec, response readWriteModel.OpcuaAcknowledgeResponse) {
	transactionId := s.channelTransactionManager.getTransactionIdentifier()

	requestHeader := readWriteModel.NewRequestHeader(readWriteModel.NewNodeId(s.authenticationToken),
		s.getCurrentDateTime(),
		0, //RequestHandle
		0,
		NULL_STRING,
		REQUEST_TIMEOUT_LONG,
		NULL_EXTENSION_OBJECT)

	var openSecureChannelRequest readWriteModel.OpenSecureChannelRequest
	if s.isEncrypted {
		openSecureChannelRequest = readWriteModel.NewOpenSecureChannelRequest(
			requestHeader,
			VERSION,
			readWriteModel.SecurityTokenRequestType_securityTokenRequestTypeIssue,
			readWriteModel.MessageSecurityMode_messageSecurityModeSignAndEncrypt,
			readWriteModel.NewPascalByteString(int32(len(s.clientNonce)), s.clientNonce),
			uint32(s.lifetime))
	} else {
		openSecureChannelRequest = readWriteModel.NewOpenSecureChannelRequest(
			requestHeader,
			VERSION,
			readWriteModel.SecurityTokenRequestType_securityTokenRequestTypeIssue,
			readWriteModel.MessageSecurityMode_messageSecurityModeNone,
			NULL_BYTE_STRING,
			uint32(s.lifetime))
	}

	identifier, err := strconv.ParseUint(openSecureChannelRequest.GetIdentifier(), 10, 16)
	if err != nil {
		s.log.Debug().Err(err).Msg("error parsing identifier")
		return
	}
	expandedNodeId := readWriteModel.NewExpandedNodeId(
		false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil,
	)

	extObject := readWriteModel.NewExtensionObject(
		expandedNodeId,
		nil,
		openSecureChannelRequest,
		false,
	)

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		s.log.Debug().Err(err).Msg("error serializing")
		return
	}

	openRequest := readWriteModel.NewOpcuaOpenRequest(
		FINAL_CHUNK,
		0,
		readWriteModel.NewPascalString(s.securityPolicy),
		s.publicCertificate,
		s.thumbprint,
		transactionId,
		transactionId,
		buffer.GetBytes(),
	)

	var apu readWriteModel.OpcuaAPU

	if s.isEncrypted {
		message, err := s.encryptionHandler.encodeMessage(ctx, openRequest, buffer.GetBytes())
		if err != nil {
			s.log.Debug().Err(err).Msg("error encoding")
			return
		}
		apu, err = readWriteModel.OpcuaAPUParse(ctx, message, false)
		if err != nil {
			s.log.Debug().Err(err).Msg("error parsing")
			return
		}
	} else {
		apu = readWriteModel.NewOpcuaAPU(openRequest, false)
	}

	requestConsumer := func(transactionId int32) {
		if err := codec.SendRequest(
			ctx,
			apu,
			func(message spi.Message) bool {
				opcuaAPU, ok := message.(readWriteModel.OpcuaAPUExactly)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				messagePDU := opcuaAPU.GetMessage()
				openResponse, ok := messagePDU.(readWriteModel.OpcuaOpenResponseExactly)
				if !ok {
					s.log.Debug().Type("type", messagePDU).Msg("Not relevant")
					return false
				}
				return openResponse.GetRequestId() == transactionId
			},
			func(message spi.Message) error {
				opcuaAPU := message.(readWriteModel.OpcuaAPU)
				messagePDU := opcuaAPU.GetMessage()
				opcuaOpenResponse := messagePDU.(readWriteModel.OpcuaOpenResponse)
				readBuffer := utils.NewReadBufferByteBased(opcuaOpenResponse.GetMessage(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
				extensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer(ctx, readBuffer, false)
				if err != nil {
					return errors.Wrap(err, "error parsing")
				}
				//Store the initial sequence number from the server. there's no requirement for the server and client to use the same starting number.
				s.senderSequenceNumber.Store(opcuaOpenResponse.GetSequenceNumber())

				if fault, ok := extensionObject.GetBody().(readWriteModel.ServiceFaultExactly); ok {
					statusCode := fault.GetResponseHeader().(readWriteModel.ResponseHeader).GetServiceResult().GetStatusCode()
					statusCodeByValue, _ := readWriteModel.OpcuaStatusCodeByValue(statusCode)
					s.log.Error().Msgf("Failed to connect to opc ua server for the following reason:- %v, %v",
						statusCode,
						statusCodeByValue)
				} else {
					s.log.Debug().Msg("Got Secure Response Connection Response")
					openSecureChannelResponse := extensionObject.GetBody().(readWriteModel.OpenSecureChannelResponse)
					s.tokenId.Store(int32(openSecureChannelResponse.GetSecurityToken().(readWriteModel.ChannelSecurityToken).GetTokenId())) // TODO: strange that int32 and uint32 missmatch
					s.channelId.Store(int32(openSecureChannelResponse.GetSecurityToken().(readWriteModel.ChannelSecurityToken).GetChannelId()))
					s.onConnectCreateSessionRequest(ctx, codec)
				}
				return nil
			},
			func(err error) error {
				s.log.Debug().Err(err).Msg("error submitting")
				return nil
			},
			REQUEST_TIMEOUT,
		); err != nil {
			s.log.Debug().Err(err).Msg("a error")
		}
	}
	s.log.Debug().Msgf("Submitting OpenSecureChannel with id of %d", transactionId)
	if err := s.channelTransactionManager.submit(requestConsumer, transactionId); err != nil {
		s.log.Debug().Err(err).Msg("error submitting")
	}
}

func (s *SecureChannel) onConnectCreateSessionRequest(ctx context.Context, codec *MessageCodec) {
	requestHeader := readWriteModel.NewRequestHeader(
		readWriteModel.NewNodeId(s.authenticationToken),
		s.getCurrentDateTime(),
		0,
		0,
		NULL_STRING,
		REQUEST_TIMEOUT_LONG,
		NULL_EXTENSION_OBJECT)

	applicationName := readWriteModel.NewLocalizedText(
		true,
		true,
		readWriteModel.NewPascalString("en"),
		APPLICATION_TEXT)

	noOfDiscoveryUrls := int32(-1)
	var discoveryUrls []readWriteModel.PascalString

	clientDescription := readWriteModel.NewApplicationDescription(APPLICATION_URI,
		PRODUCT_URI,
		applicationName,
		readWriteModel.ApplicationType_applicationTypeClient,
		NULL_STRING,
		NULL_STRING,
		noOfDiscoveryUrls,
		discoveryUrls)

	createSessionRequest := readWriteModel.NewCreateSessionRequest(
		requestHeader,
		clientDescription,
		NULL_STRING,
		s.endpoint,
		readWriteModel.NewPascalString(s.sessionName),
		readWriteModel.NewPascalByteString(int32(len(s.clientNonce)), s.clientNonce),
		NULL_BYTE_STRING,
		120000,
		0,
	)

	identifier, err := strconv.ParseUint(createSessionRequest.GetIdentifier(), 10, 16)
	if err != nil {
		s.log.Debug().Err(err).Msg("error parsing identifier")
		return
	}
	expandedNodeId := readWriteModel.NewExpandedNodeId(
		false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil)

	extObject := readWriteModel.NewExtensionObject(
		expandedNodeId,
		nil,
		createSessionRequest,
		false,
	)

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		s.log.Debug().Err(err).Msg("error serializing")
		return
	}

	consumer := func(opcuaResponse []byte) {
		message, err := readWriteModel.ExtensionObjectParseWithBuffer(ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
		if err != nil {
			s.log.Error().Err(err).Msg("error parsing")
			return
		}
		if fault, ok := message.GetBody().(readWriteModel.ServiceFaultExactly); ok {
			statusCode := fault.GetResponseHeader().(readWriteModel.ResponseHeader).GetServiceResult().GetStatusCode()
			statusCodeByValue, _ := readWriteModel.OpcuaStatusCodeByValue(statusCode)
			s.log.Error().Msgf("Failed to connect to opc ua server for the following reason:- %v, %v",
				statusCode,
				statusCodeByValue)
		} else {
			s.log.Debug().Msg("Got Create Session Response Connection Response")

			extensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer(ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
			if err != nil {
				s.log.Error().Err(err).Msg("error parsing")
				return
			}
			unknownExtensionObject := extensionObject.GetBody()
			if responseMessage, ok := unknownExtensionObject.(readWriteModel.CreateSessionResponseExactly); ok {
				s.authenticationToken = responseMessage.GetAuthenticationToken().GetNodeId()

				s.onConnectActivateSessionRequest(ctx, codec, responseMessage, message.GetBody().(readWriteModel.CreateSessionResponse))
			} else {
				serviceFault := unknownExtensionObject.(readWriteModel.ServiceFault)
				header := serviceFault.GetResponseHeader().(readWriteModel.ResponseHeader)
				s.log.Error().Msgf("Subscription ServiceFault returned from server with error code, '%s'", header.GetServiceResult())
			}
		}
	}

	errorDispatcher := func(err error) {
		s.log.Error().Err(err).Msg("Error while waiting for subscription response")
	}

	s.submit(ctx, codec, errorDispatcher, consumer, buffer)
}

func (s *SecureChannel) onConnectActivateSessionRequest(ctx context.Context, codec *MessageCodec, opcuaMessageResponse readWriteModel.CreateSessionResponse, sessionResponse readWriteModel.CreateSessionResponse) {
	s.senderCertificate = sessionResponse.GetServerCertificate().GetStringValue()
	certificate, err := s.encryptionHandler.getCertificateX509(s.senderCertificate)
	if err != nil {
		s.log.Error().Err(err).Msg("error getting certificate")
		return
	}
	s.encryptionHandler.setServerCertificate(certificate)
	s.senderNonce = sessionResponse.GetServerNonce().GetStringValue()
	endpoints := make([]string, 3)
	if address, err := url.Parse(s.configuration.host); err != nil {
		if names, err := net.LookupAddr(address.Host); err != nil {
			endpoints[0] = "opc.tcp://" + names[rand.Intn(len(names))] + ":" + s.configuration.port + s.configuration.transportEndpoint
		}
		endpoints[1] = "opc.tcp://" + address.Hostname() + ":" + s.configuration.port + s.configuration.transportEndpoint
		//endpoints[2] = "opc.tcp://" + address.getCanonicalHostName() + ":" + s.configuration.getPort() + s.configuration.transportEndpoint// TODO: not sure how to get that in golang
	}

	s.selectEndpoint(sessionResponse)

	if s.policyId == nil {
		s.log.Error().Msg("Unable to find endpoint - " + endpoints[1])
		return
	}

	userIdentityToken := s.getIdentityToken(s.tokenType, s.policyId.GetStringValue())

	requestHandle := s.getRequestHandle()

	requestHeader := readWriteModel.NewRequestHeader(
		readWriteModel.NewNodeId(s.authenticationToken),
		s.getCurrentDateTime(),
		requestHandle,
		0,
		NULL_STRING,
		REQUEST_TIMEOUT_LONG,
		NULL_EXTENSION_OBJECT)

	clientSignature := readWriteModel.NewSignatureData(NULL_STRING, NULL_BYTE_STRING)

	activateSessionRequest := readWriteModel.NewActivateSessionRequest(
		requestHeader,
		clientSignature,
		0,
		nil,
		0,
		nil,
		userIdentityToken,
		clientSignature)

	identifier, err := strconv.ParseUint(activateSessionRequest.GetIdentifier(), 10, 16)
	if err != nil {
		s.log.Debug().Err(err).Msg("error parsing identifier")
		return
	}

	expandedNodeId := readWriteModel.NewExpandedNodeId(false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil)

	extObject := readWriteModel.NewExtensionObject(
		expandedNodeId,
		nil,
		activateSessionRequest,
		false,
	)

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		s.log.Debug().Err(err).Msg("error serializing")
		return
	}

	consumer := func(opcuaResponse []byte) {
		message, err := readWriteModel.ExtensionObjectParseWithBuffer(ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
		if err != nil {
			s.log.Error().Err(err).Msg("error parsing")
			return
		}
		if fault, ok := message.GetBody().(readWriteModel.ServiceFaultExactly); ok {
			statusCode := fault.GetResponseHeader().(readWriteModel.ResponseHeader).GetServiceResult().GetStatusCode()
			statusCodeByValue, _ := readWriteModel.OpcuaStatusCodeByValue(statusCode)
			s.log.Error().Msgf("Failed to connect to opc ua server for the following reason:- %v, %v",
				statusCode,
				statusCodeByValue)
		} else {
			s.log.Debug().Msg("Got Activate Session Response Connection Response")

			extensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer(ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
			if err != nil {
				s.log.Error().Err(err).Msg("error parsing")
				return
			}
			unknownExtensionObject := extensionObject.GetBody()
			if responseMessage, ok := unknownExtensionObject.(readWriteModel.ActivateSessionResponseExactly); ok {
				returnedRequestHandle := responseMessage.GetResponseHeader().(readWriteModel.ResponseHeader).GetRequestHandle()
				if !(requestHandle == returnedRequestHandle) {
					s.log.Error().Msgf("Request handle isn't as expected, we might have missed a packet. %d != %d", requestHandle, returnedRequestHandle)
				}

				// Send an event that connection setup is complete.
				s.keepAlive()
				codec.fireConnected()
			} else {
				serviceFault := unknownExtensionObject.(readWriteModel.ServiceFault)
				header := serviceFault.GetResponseHeader().(readWriteModel.ResponseHeader)
				s.log.Error().Msgf("Subscription ServiceFault returned from server with error code, '%s'", header.GetServiceResult())
			}
		}
	}

	errorDispatcher := func(err error) {
		s.log.Error().Err(err).Msg("Error while waiting for subscription response")
	}

	s.submit(ctx, codec, errorDispatcher, consumer, buffer)
}

func (s *SecureChannel) onDisconnect(ctx context.Context, codec *MessageCodec) {
	s.log.Info().Msg("disconnecting")
	requestHandle := s.getRequestHandle()

	s.keepAliveIndicator.Store(false)

	expandedNodeId := readWriteModel.NewExpandedNodeId(false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, 473),
		nil,
		nil) //Identifier for OpenSecureChannel

	requestHeader := readWriteModel.NewRequestHeader(
		readWriteModel.NewNodeId(s.authenticationToken),
		s.getCurrentDateTime(),
		requestHandle, //RequestHandle
		0,
		NULL_STRING,
		5000,
		NULL_EXTENSION_OBJECT)

	closeSessionRequest := readWriteModel.NewCloseSessionRequest(
		requestHeader,
		true)

	extObject := readWriteModel.NewExtensionObject(
		expandedNodeId,
		nil,
		closeSessionRequest,
		false,
	)

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		s.log.Debug().Err(err).Msg("error serializing")
		return
	}

	consumer := func(opcuaResponse []byte) {
		message, err := readWriteModel.ExtensionObjectParseWithBuffer(ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
		if err != nil {
			s.log.Error().Err(err).Msg("error parsing")
			return
		}
		if fault, ok := message.GetBody().(readWriteModel.ServiceFaultExactly); ok {
			statusCode := fault.GetResponseHeader().(readWriteModel.ResponseHeader).GetServiceResult().GetStatusCode()
			statusCodeByValue, _ := readWriteModel.OpcuaStatusCodeByValue(statusCode)
			s.log.Error().Msgf("Failed to connect to opc ua server for the following reason:- %v, %v",
				statusCode,
				statusCodeByValue)
		} else {
			s.log.Debug().Msg("Got Close Session Response Connection Response")

			extensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer(ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
			if err != nil {
				s.log.Error().Err(err).Msg("error parsing")
				return
			}
			unknownExtensionObject := extensionObject.GetBody()
			if responseMessage, ok := unknownExtensionObject.(readWriteModel.CloseSessionResponseExactly); ok {
				s.onDisconnectCloseSecureChannel(ctx, codec, responseMessage, message.GetBody().(readWriteModel.CloseSessionResponse))
			} else {
				serviceFault := unknownExtensionObject.(readWriteModel.ServiceFault)
				header := serviceFault.GetResponseHeader().(readWriteModel.ResponseHeader)
				s.log.Error().Msgf("Subscription ServiceFault returned from server with error code, '%s'", header.GetServiceResult())
			}
		}
	}

	errorDispatcher := func(err error) {
		s.log.Error().Err(err).Msg("Error while waiting for close session response")
	}

	s.submit(ctx, codec, errorDispatcher, consumer, buffer)
}

func (s *SecureChannel) onDisconnectCloseSecureChannel(ctx context.Context, codec *MessageCodec, message readWriteModel.CloseSessionResponseExactly, response readWriteModel.CloseSessionResponse) {
	transactionId := s.channelTransactionManager.getTransactionIdentifier()

	requestHeader := readWriteModel.NewRequestHeader(readWriteModel.NewNodeId(s.authenticationToken),
		s.getCurrentDateTime(),
		0, //RequestHandle
		0,
		NULL_STRING,
		REQUEST_TIMEOUT_LONG,
		NULL_EXTENSION_OBJECT)

	closeSecureChannelRequest := readWriteModel.NewCloseSecureChannelRequest(requestHeader)

	identifier, err := strconv.ParseUint(closeSecureChannelRequest.GetIdentifier(), 10, 16)
	if err != nil {
		s.log.Debug().Err(err).Msg("error parsing identifier")
		return
	}
	expandedNodeId := readWriteModel.NewExpandedNodeId(
		false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil,
	)

	closeRequest := readWriteModel.NewOpcuaCloseRequest(FINAL_CHUNK,
		s.channelId.Load(),
		s.tokenId.Load(),
		transactionId,
		transactionId,
		readWriteModel.NewExtensionObject(
			expandedNodeId,
			nil,
			closeSecureChannelRequest,
			false,
		),
	)

	apu := readWriteModel.NewOpcuaAPU(closeRequest, false)

	requestConsumer := func(transactionId int32) {
		if err := codec.SendRequest(
			ctx,
			apu,
			func(message spi.Message) bool {
				opcuaAPU, ok := message.(readWriteModel.OpcuaAPUExactly)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				messagePDU := opcuaAPU.GetMessage()
				openResponse, ok := messagePDU.(readWriteModel.OpcuaMessageResponseExactly)
				if !ok {
					s.log.Debug().Type("type", messagePDU).Msg("Not relevant")
					return false
				}
				return openResponse.GetRequestId() == transactionId
			},
			func(message spi.Message) error {
				opcuaAPU := message.(readWriteModel.OpcuaAPU)
				messagePDU := opcuaAPU.GetMessage()
				opcuaMessageResponse := messagePDU.(readWriteModel.OpcuaMessageResponse)
				s.log.Trace().Msgf("Got close secure channel response:\n%s", opcuaMessageResponse)
				return nil
			},
			func(err error) error {
				s.log.Debug().Err(err).Msg("error submitting")
				return nil
			},
			REQUEST_TIMEOUT,
		); err != nil {
			s.log.Debug().Err(err).Msg("a error")
		}
	}
	s.log.Debug().Msgf("Submitting CloseSecureChannel with id of %d", transactionId)
	if err := s.channelTransactionManager.submit(requestConsumer, transactionId); err != nil {
		s.log.Debug().Err(err).Msg("error submitting")
	}
}

func (s *SecureChannel) onDiscover(ctx context.Context, codec *MessageCodec) {
	// Only the TCP transport supports login.
	s.log.Debug().Msg("Opcua Driver running in ACTIVE mode, discovering endpoints")

	hello := readWriteModel.NewOpcuaHelloRequest(FINAL_CHUNK,
		VERSION,
		DEFAULT_RECEIVE_BUFFER_SIZE,
		DEFAULT_SEND_BUFFER_SIZE,
		DEFAULT_MAX_MESSAGE_SIZE,
		DEFAULT_MAX_CHUNK_COUNT,
		s.endpoint,
	)

	apu := readWriteModel.NewOpcuaAPU(hello, false)

	requestConsumer := func(transactionId int32) {
		if err := codec.SendRequest(
			ctx,
			apu,
			func(message spi.Message) bool {
				opcuaAPU, ok := message.(readWriteModel.OpcuaAPUExactly)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				messagePDU := opcuaAPU.GetMessage()
				_, ok = messagePDU.(readWriteModel.OpcuaAcknowledgeResponseExactly)
				if !ok {
					s.log.Debug().Type("type", messagePDU).Msg("Not relevant")
					return false
				}
				return true
			},
			func(message spi.Message) error {
				opcuaAPU := message.(readWriteModel.OpcuaAPU)
				messagePDU := opcuaAPU.GetMessage()
				opcuaAcknowledgeResponse := messagePDU.(readWriteModel.OpcuaAcknowledgeResponse)
				s.log.Trace().Msgf("Got Hello Response Connection Response:\n%s", opcuaAcknowledgeResponse)
				s.onDiscoverOpenSecureChannel(ctx, codec, opcuaAcknowledgeResponse)
				return nil
			},
			func(err error) error {
				s.log.Debug().Err(err).Msg("error submitting")
				return nil
			},
			REQUEST_TIMEOUT,
		); err != nil {
			s.log.Debug().Err(err).Msg("a error")
		}
	}
	if err := s.channelTransactionManager.submit(requestConsumer, s.channelTransactionManager.getTransactionIdentifier()); err != nil {
		s.log.Debug().Err(err).Msg("error submitting")
	}
}

func (s *SecureChannel) onDiscoverOpenSecureChannel(ctx context.Context, codec *MessageCodec, opcuaAcknowledgeResponse readWriteModel.OpcuaAcknowledgeResponse) {
	transactionId := s.channelTransactionManager.getTransactionIdentifier()

	requestHeader := readWriteModel.NewRequestHeader(readWriteModel.NewNodeId(s.authenticationToken),
		s.getCurrentDateTime(),
		0, //RequestHandle
		0,
		NULL_STRING,
		REQUEST_TIMEOUT_LONG,
		NULL_EXTENSION_OBJECT)

	openSecureChannelRequest := readWriteModel.NewOpenSecureChannelRequest(
		requestHeader,
		VERSION,
		readWriteModel.SecurityTokenRequestType_securityTokenRequestTypeIssue,
		readWriteModel.MessageSecurityMode_messageSecurityModeNone,
		NULL_BYTE_STRING,
		s.lifetime,
	)

	identifier, err := strconv.ParseUint(openSecureChannelRequest.GetIdentifier(), 10, 16)
	if err != nil {
		s.log.Debug().Err(err).Msg("error parsing identifier")
		return
	}
	expandedNodeId := readWriteModel.NewExpandedNodeId(
		false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil,
	)

	extObject := readWriteModel.NewExtensionObject(
		expandedNodeId,
		nil,
		openSecureChannelRequest,
		false,
	)

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		s.log.Debug().Err(err).Msg("error serializing")
		return
	}

	openRequest := readWriteModel.NewOpcuaOpenRequest(
		FINAL_CHUNK,
		0,
		SECURITY_POLICY_NONE,
		NULL_BYTE_STRING,
		NULL_BYTE_STRING,
		transactionId,
		transactionId,
		buffer.GetBytes(),
	)

	apu := readWriteModel.NewOpcuaAPU(openRequest, false)

	requestConsumer := func(transactionId int32) {
		if err := codec.SendRequest(
			ctx,
			apu,
			func(message spi.Message) bool {
				opcuaAPU, ok := message.(readWriteModel.OpcuaAPUExactly)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				messagePDU := opcuaAPU.GetMessage()
				openResponse, ok := messagePDU.(readWriteModel.OpcuaOpenResponseExactly)
				if !ok {
					s.log.Debug().Type("type", messagePDU).Msg("Not relevant")
					return false
				}
				return openResponse.GetRequestId() == transactionId
			},
			func(message spi.Message) error {
				opcuaAPU := message.(readWriteModel.OpcuaAPU)
				messagePDU := opcuaAPU.GetMessage()
				opcuaOpenResponse := messagePDU.(readWriteModel.OpcuaOpenResponse)
				readBuffer := utils.NewReadBufferByteBased(opcuaOpenResponse.GetMessage(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
				extensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer(ctx, readBuffer, false)
				if err != nil {
					return errors.Wrap(err, "error parsing")
				}

				if fault, ok := extensionObject.GetBody().(readWriteModel.ServiceFaultExactly); ok {
					statusCode := fault.GetResponseHeader().(readWriteModel.ResponseHeader).GetServiceResult().GetStatusCode()
					statusCodeByValue, _ := readWriteModel.OpcuaStatusCodeByValue(statusCode)
					s.log.Error().Msgf("Failed to connect to opc ua server for the following reason:- %v, %v",
						statusCode,
						statusCodeByValue)
				} else {
					s.log.Debug().Msg("Got Secure Response Connection Response")
					openSecureChannelResponse := extensionObject.GetBody().(readWriteModel.OpenSecureChannelResponse)
					s.onDiscoverGetEndpointsRequest(ctx, codec, opcuaOpenResponse, openSecureChannelResponse)
				}
				return nil
			},
			func(err error) error {
				s.log.Debug().Err(err).Msg("error submitting")
				return nil
			},
			REQUEST_TIMEOUT,
		); err != nil {
			s.log.Debug().Err(err).Msg("a error")
		}
	}
	if err := s.channelTransactionManager.submit(requestConsumer, transactionId); err != nil {
		s.log.Debug().Err(err).Msg("error submitting")
	}
}

func (s *SecureChannel) onDiscoverGetEndpointsRequest(ctx context.Context, codec *MessageCodec, opcuaOpenResponse readWriteModel.OpcuaOpenResponse, openSecureChannelResponse readWriteModel.OpenSecureChannelResponse) {
	s.tokenId.Store(int32(openSecureChannelResponse.GetSecurityToken().(readWriteModel.ChannelSecurityToken).GetTokenId()))
	s.channelId.Store(int32(openSecureChannelResponse.GetSecurityToken().(readWriteModel.ChannelSecurityToken).GetChannelId()))

	transactionId := s.channelTransactionManager.getTransactionIdentifier()

	nextSequenceNumber := opcuaOpenResponse.GetSequenceNumber() + 1
	nextRequestId := opcuaOpenResponse.GetRequestId() + 1

	if !(transactionId == nextSequenceNumber) {
		s.log.Error().Msgf("Sequence number isn't as expected, we might have missed a packet. - %d != %d", transactionId, nextSequenceNumber)
		return
	}

	requestHeader := readWriteModel.NewRequestHeader(readWriteModel.NewNodeId(s.authenticationToken),
		s.getCurrentDateTime(),
		0, //RequestHandle
		0,
		NULL_STRING,
		REQUEST_TIMEOUT_LONG,
		NULL_EXTENSION_OBJECT)

	endpointsRequest := readWriteModel.NewGetEndpointsRequest(
		requestHeader,
		s.endpoint,
		0,
		nil,
		0,
		nil)

	identifier, err := strconv.ParseUint(endpointsRequest.GetIdentifier(), 10, 16)
	if err != nil {
		s.log.Debug().Err(err).Msg("error parsing identifier")
		return
	}
	expandedNodeId := readWriteModel.NewExpandedNodeId(
		false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil,
	)

	extObject := readWriteModel.NewExtensionObject(
		expandedNodeId,
		nil,
		endpointsRequest,
		false,
	)

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		s.log.Debug().Err(err).Msg("error serializing")
		return
	}

	messageRequest := readWriteModel.NewOpcuaMessageRequest(
		FINAL_CHUNK,
		s.channelId.Load(),
		s.tokenId.Load(),
		nextSequenceNumber,
		nextRequestId,
		buffer.GetBytes(),
	)

	apu := readWriteModel.NewOpcuaAPU(messageRequest, false)

	requestConsumer := func(transactionId int32) {
		if err := codec.SendRequest(
			ctx,
			apu,
			func(message spi.Message) bool {
				opcuaAPU, ok := message.(readWriteModel.OpcuaAPUExactly)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				messagePDU := opcuaAPU.GetMessage()
				messageResponse, ok := messagePDU.(readWriteModel.OpcuaMessageResponseExactly)
				if !ok {
					s.log.Debug().Type("type", messagePDU).Msg("Not relevant")
					return false
				}
				return messageResponse.GetRequestId() == transactionId
			},
			func(message spi.Message) error {
				opcuaAPU := message.(readWriteModel.OpcuaAPU)
				messagePDU := opcuaAPU.GetMessage()
				messageResponse := messagePDU.(readWriteModel.OpcuaMessageResponse)
				readBuffer := utils.NewReadBufferByteBased(messageResponse.GetMessage(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
				extensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer(ctx, readBuffer, false)
				if err != nil {
					return errors.Wrap(err, "error parsing")
				}

				if fault, ok := extensionObject.GetBody().(readWriteModel.ServiceFaultExactly); ok {
					statusCode := fault.GetResponseHeader().(readWriteModel.ResponseHeader).GetServiceResult().GetStatusCode()
					statusCodeByValue, _ := readWriteModel.OpcuaStatusCodeByValue(statusCode)
					s.log.Error().Msgf("Failed to connect to opc ua server for the following reason:- %v, %v",
						statusCode,
						statusCodeByValue)
				} else {
					s.log.Debug().Msg("Got Secure Response Connection Response")
					response := extensionObject.GetBody().(readWriteModel.GetEndpointsResponse)

					endpoints := response.GetEndpoints()
					for _, endpoint := range endpoints {
						endpointDescription := endpoint.(readWriteModel.EndpointDescription)
						if endpointDescription.GetEndpointUrl().GetStringValue() == (s.endpoint.GetStringValue()) && endpointDescription.GetSecurityPolicyUri().GetStringValue() == (s.securityPolicy) {
							s.log.Info().Msgf("Found OPC UA endpoint %s", s.endpoint.GetStringValue())
							s.configuration.senderCertificate = endpointDescription.GetServerCertificate().GetStringValue()
						}
					}

					digest := sha1.Sum(s.configuration.senderCertificate)
					s.thumbprint = readWriteModel.NewPascalByteString(int32(len(digest)), digest[:])

					s.onDiscoverCloseSecureChannel(ctx, codec, response)
				}
				return nil
			},
			func(err error) error {
				s.log.Debug().Err(err).Msg("error submitting")
				return nil
			},
			REQUEST_TIMEOUT,
		); err != nil {
			s.log.Debug().Err(err).Msg("a error")
		}
	}
	if err := s.channelTransactionManager.submit(requestConsumer, transactionId); err != nil {
		s.log.Debug().Err(err).Msg("error submitting")
	}
}

func (s *SecureChannel) onDiscoverCloseSecureChannel(ctx context.Context, codec *MessageCodec, response readWriteModel.GetEndpointsResponse) {
	transactionId := s.channelTransactionManager.getTransactionIdentifier()

	requestHeader := readWriteModel.NewRequestHeader(readWriteModel.NewNodeId(s.authenticationToken),
		s.getCurrentDateTime(),
		0, //RequestHandle
		0,
		NULL_STRING,
		REQUEST_TIMEOUT_LONG,
		NULL_EXTENSION_OBJECT)

	closeSecureChannelRequest := readWriteModel.NewCloseSecureChannelRequest(requestHeader)

	identifier, err := strconv.ParseUint(closeSecureChannelRequest.GetIdentifier(), 10, 16)
	if err != nil {
		s.log.Debug().Err(err).Msg("error parsing identifier")
		return
	}
	expandedNodeId := readWriteModel.NewExpandedNodeId(
		false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil,
	)

	closeRequest := readWriteModel.NewOpcuaCloseRequest(FINAL_CHUNK,
		s.channelId.Load(),
		s.tokenId.Load(),
		transactionId,
		transactionId,
		readWriteModel.NewExtensionObject(
			expandedNodeId,
			nil,
			closeSecureChannelRequest,
			false,
		),
	)

	apu := readWriteModel.NewOpcuaAPU(closeRequest, false)

	requestConsumer := func(transactionId int32) {
		if err := codec.SendRequest(
			ctx,
			apu,
			func(message spi.Message) bool {
				opcuaAPU, ok := message.(readWriteModel.OpcuaAPUExactly)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				messagePDU := opcuaAPU.GetMessage()
				openResponse, ok := messagePDU.(readWriteModel.OpcuaMessageResponseExactly)
				if !ok {
					s.log.Debug().Type("type", messagePDU).Msg("Not relevant")
					return false
				}
				return openResponse.GetRequestId() == transactionId
			},
			func(message spi.Message) error {
				opcuaAPU := message.(readWriteModel.OpcuaAPU)
				messagePDU := opcuaAPU.GetMessage()
				opcuaMessageResponse := messagePDU.(readWriteModel.OpcuaMessageResponse)
				s.log.Trace().Msgf("Got close secure channel response:\n%s", opcuaMessageResponse)
				return nil
			},
			func(err error) error {
				s.log.Debug().Err(err).Msg("error submitting")
				return nil
			},
			REQUEST_TIMEOUT,
		); err != nil {
			s.log.Debug().Err(err).Msg("a error")
		}
	}
	s.log.Debug().Msgf("Submitting CloseSecureChannel with id of %d", transactionId)
	if err := s.channelTransactionManager.submit(requestConsumer, transactionId); err != nil {
		s.log.Debug().Err(err).Msg("error submitting")
	}
}

func (s *SecureChannel) keepAlive() {
	s.keepAliveStateChange.Lock()
	defer s.keepAliveStateChange.Unlock()
	if s.keepAliveIndicator.Load() {
		s.log.Warn().Msg("keepalive already running")
		return
	}
	s.keepAliveWg.Add(1)
	go func() {
		defer s.keepAliveWg.Done()
		s.keepAliveIndicator.Store(true)
		defer s.keepAliveIndicator.Store(false)
		defer s.log.Info().Msg("ending keepalive")
		ctx := context.Background()
		for (s.codec == nil || s.codec.IsRunning()) && s.keepAliveIndicator.Load() {
			sleepTime := time.Duration(math.Ceil(float64(s.lifetime)*0.75)) * time.Millisecond
			s.log.Trace().Dur("sleepTime", sleepTime).Msg("Sleeping")
			time.Sleep(sleepTime)

			transactionId := s.channelTransactionManager.getTransactionIdentifier()

			requestHeader := readWriteModel.NewRequestHeader(readWriteModel.NewNodeId(s.authenticationToken),
				s.getCurrentDateTime(),
				0, //RequestHandle
				0,
				NULL_STRING,
				REQUEST_TIMEOUT_LONG,
				NULL_EXTENSION_OBJECT)

			var openSecureChannelRequest readWriteModel.OpenSecureChannelRequest
			if s.isEncrypted {
				openSecureChannelRequest = readWriteModel.NewOpenSecureChannelRequest(
					requestHeader,
					VERSION,
					readWriteModel.SecurityTokenRequestType_securityTokenRequestTypeIssue,
					readWriteModel.MessageSecurityMode_messageSecurityModeSignAndEncrypt,
					readWriteModel.NewPascalByteString(int32(len(s.clientNonce)), s.clientNonce),
					uint32(s.lifetime))
			} else {
				openSecureChannelRequest = readWriteModel.NewOpenSecureChannelRequest(
					requestHeader,
					VERSION,
					readWriteModel.SecurityTokenRequestType_securityTokenRequestTypeIssue,
					readWriteModel.MessageSecurityMode_messageSecurityModeNone,
					NULL_BYTE_STRING,
					uint32(s.lifetime))
			}
			identifier, err := strconv.ParseUint(openSecureChannelRequest.GetIdentifier(), 10, 16)
			if err != nil {
				s.log.Error().Err(err).Msg("error parsing identifier")
				return
			}

			expandedNodeId := readWriteModel.NewExpandedNodeId(false, //Namespace Uri Specified
				false, //Server Index Specified
				readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
				nil,
				nil)

			extObject := readWriteModel.NewExtensionObject(
				expandedNodeId,
				nil,
				openSecureChannelRequest,
				false,
			)

			buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
			if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
				s.log.Error().Err(err).Msg("error serializing")
				return
			}

			openRequest := readWriteModel.NewOpcuaOpenRequest(
				FINAL_CHUNK,
				0,
				readWriteModel.NewPascalString(s.securityPolicy),
				s.publicCertificate,
				s.thumbprint,
				transactionId,
				transactionId,
				buffer.GetBytes(),
			)

			var apu readWriteModel.OpcuaAPU

			if s.isEncrypted {
				message, err := s.encryptionHandler.encodeMessage(ctx, openRequest, buffer.GetBytes())
				if err != nil {
					s.log.Error().Err(err).Msg("error encoding")
					return
				}
				apu, err = readWriteModel.OpcuaAPUParse(ctx, message, false)
				if err != nil {
					s.log.Error().Err(err).Msg("error parsing")
					return
				}
			} else {
				apu = readWriteModel.NewOpcuaAPU(openRequest, false)
			}

			requestConsumer := func(transactionId int32) {
				if err := s.codec.SendRequest(
					ctx,
					apu,
					func(message spi.Message) bool {
						opcuaAPU, ok := message.(readWriteModel.OpcuaAPUExactly)
						if !ok {
							s.log.Debug().Type("type", message).Msg("Not relevant")
							return false
						}
						messagePDU := opcuaAPU.GetMessage()
						openResponse, ok := messagePDU.(readWriteModel.OpcuaOpenResponseExactly)
						if !ok {
							s.log.Debug().Type("type", messagePDU).Msg("Not relevant")
							return false
						}
						return openResponse.GetRequestId() == transactionId
					},
					func(message spi.Message) error {
						opcuaAPU := message.(readWriteModel.OpcuaAPU)
						messagePDU := opcuaAPU.GetMessage()
						opcuaOpenResponse := messagePDU.(readWriteModel.OpcuaOpenResponse)
						readBuffer := utils.NewReadBufferByteBased(opcuaOpenResponse.GetMessage(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
						extensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer(ctx, readBuffer, false)
						if err != nil {
							return errors.Wrap(err, "error parsing")
						}

						if fault, ok := extensionObject.GetBody().(readWriteModel.ServiceFaultExactly); ok {
							statusCode := fault.GetResponseHeader().(readWriteModel.ResponseHeader).GetServiceResult().GetStatusCode()
							statusCodeByValue, _ := readWriteModel.OpcuaStatusCodeByValue(statusCode)
							s.log.Error().Msgf("Failed to connect to opc ua server for the following reason:- %v, %v",
								statusCode,
								statusCodeByValue)
						} else {
							s.log.Debug().Msg("Got Secure Response Connection Response")
							openSecureChannelResponse := extensionObject.GetBody().(readWriteModel.OpenSecureChannelResponse)
							token := openSecureChannelResponse.GetSecurityToken().(readWriteModel.ChannelSecurityToken)
							s.tokenId.Store(int32(token.GetTokenId())) // TODO: strange that int32 and uint32 missmatch
							s.channelId.Store(int32(token.GetChannelId()))
							s.lifetime = token.GetRevisedLifetime()
						}
						return nil
					},
					func(err error) error {
						s.log.Debug().Err(err).Msg("error submitting")
						return nil
					},
					REQUEST_TIMEOUT,
				); err != nil {
					s.log.Debug().Err(err).Msg("a error")
				}
			}
			s.log.Debug().Msgf("Submitting OpenSecureChannel with id of %d", transactionId)
			if err := s.channelTransactionManager.submit(requestConsumer, transactionId); err != nil {
				s.log.Debug().Err(err).Msg("error submitting")
			}
		}
	}()
	return
}

// getRequestHandle returns the next request handle
func (s *SecureChannel) getRequestHandle() uint32 {
	return s.requestHandleGenerator.Add(1) - 1
}

// getAuthenticationToken returns the authentication token for the current connection
func (s *SecureChannel) getAuthenticationToken() readWriteModel.NodeId {
	return readWriteModel.NewNodeId(s.authenticationToken)
}

// getChannelId gets the Channel identifier for the current channel
func (s *SecureChannel) getChannelId() int32 {
	return s.channelId.Load()
}

// getTokenId gets the Token Identifier
func (s *SecureChannel) getTokenId() int32 {
	return s.tokenId.Load()
}

// selectEndpoint Selects the endpoint to use based on the connection string provided.
//   - If Discovery is disabled it will use the host address return from the server
//   - @param sessionResponse - The CreateSessionResponse message returned by the server
//   - @throws PlcRuntimeException - If no endpoint with a compatible policy is found raise and error.
func (s *SecureChannel) selectEndpoint(sessionResponse readWriteModel.CreateSessionResponse) {
	// Get a list of the endpoints which match ours.
	var filteredEndpoints []readWriteModel.EndpointDescription
	for _, endpoint := range sessionResponse.GetServerEndpoints() {
		endpointDescription := endpoint.(readWriteModel.EndpointDescription)
		if s.isEndpoint(endpointDescription) {
			filteredEndpoints = append(filteredEndpoints, endpointDescription)
		}
	}

	//Determine if the requested security policy is included in the endpoint
	for _, endpoint := range filteredEndpoints {
		userIdentityTokens := make([]readWriteModel.UserTokenPolicy, len(endpoint.GetUserIdentityTokens()))
		for i, definition := range endpoint.GetUserIdentityTokens() {
			userIdentityTokens[i] = definition.(readWriteModel.UserTokenPolicy)
		}
		s.hasIdentity(userIdentityTokens)
	}

	if s.policyId == nil {
		s.log.Error().Msgf("Unable to find endpoint - %s", s.endpoints[0])
		return
	}

	if s.tokenType == 0xffffffff { // TODO: what did we use as undefined
		s.log.Error().Msgf("Unable to find Security Policy for endpoint - %s", s.endpoints[0])
		return
	}
}

// isEndpoint checks each component of the return endpoint description against the connection string.
//   - If all are correct then return true.
//   - @param endpoint - EndpointDescription returned from server
//   - @return true if this endpoint matches our configuration
//   - @return error - If the returned endpoint string doesn't match the format expected
func (s *SecureChannel) isEndpoint(endpoint readWriteModel.EndpointDescription) bool {
	// Split up the connection string into its individual segments.
	matches := utils.GetSubgroupMatches(URI_PATTERN, endpoint.GetEndpointUrl().GetStringValue())
	if len(matches) == 0 {
		s.log.Error().Msgf("Endpoint returned from the server doesn't match the format '{protocol-code}:({transport-code})?//{transport-host}(:{transport-port})(/{transport-endpoint})'")
		return false
	}
	s.log.Trace().Msgf("Using Endpoint %s %s %s", matches["transportHost"], matches["transportPort"], matches["transportEndpoint"])

	if s.configuration.discovery && !slices.Contains(s.endpoints, matches["transportHost"]) {
		return false
	}

	if s.configuration.port != matches["transportPort"] {
		return false
	}

	if s.configuration.transportEndpoint != matches["transportEndpoint"] {
		return false
	}

	if !s.configuration.discovery {
		s.configuration.host = matches["transportHost"]
	}

	return true
}

// hasIdentity confirms that a policy that matches the connection string is available from
//   - the returned endpoints. It sets the policyId and tokenType for the policy to use.
//   - @param policies - A list of policies returned with the endpoint description.
func (s *SecureChannel) hasIdentity(policies []readWriteModel.UserTokenPolicy) {
	for _, identityToken := range policies {
		if (identityToken.GetTokenType() == readWriteModel.UserTokenType_userTokenTypeAnonymous) && (s.username == "") {
			s.policyId = identityToken.GetPolicyId()
			s.tokenType = identityToken.GetTokenType()
		} else if (identityToken.GetTokenType() == readWriteModel.UserTokenType_userTokenTypeUserName) && (s.username != "") {
			s.policyId = identityToken.GetPolicyId()
			s.tokenType = identityToken.GetTokenType()
		}
	}
}

// getIdentityToken creates an IdentityToken to authenticate with a server.
//   - @param tokenType      the token type
//   - @param securityPolicy the security policy
//   - @return returns an ExtensionObject with an IdentityToken.
func (s *SecureChannel) getIdentityToken(tokenType readWriteModel.UserTokenType, value string) readWriteModel.ExtensionObject {
	switch tokenType {
	case readWriteModel.UserTokenType_userTokenTypeAnonymous:
		//If we aren't using authentication tell the server we would like to log in anonymously
		anonymousIdentityToken := readWriteModel.NewAnonymousIdentityToken()
		extExpandedNodeId := readWriteModel.NewExpandedNodeId(
			false, //Namespace Uri Specified
			false, //Server Index Specified
			readWriteModel.NewNodeIdFourByte(
				0, uint16(readWriteModel.OpcuaNodeIdServices_AnonymousIdentityToken_Encoding_DefaultBinary)),
			nil,
			nil,
		)
		return readWriteModel.NewExtensionObject(
			extExpandedNodeId,
			readWriteModel.NewExtensionObjectEncodingMask(false, false, true),
			readWriteModel.NewUserIdentityToken(readWriteModel.NewPascalString(s.securityPolicy), anonymousIdentityToken),
			false,
		)
	case readWriteModel.UserTokenType_userTokenTypeUserName:
		//Encrypt the password using the server nonce and server public key
		passwordBytes := []byte(s.password)
		encodeableBuffer := new(bytes.Buffer)
		var err error
		err = binary.Write(encodeableBuffer, binary.LittleEndian, len(passwordBytes)+len(s.senderNonce))
		s.log.Debug().Err(err).Msg("write")
		err = binary.Write(encodeableBuffer, binary.LittleEndian, passwordBytes)
		s.log.Debug().Err(err).Msg("write")
		err = binary.Write(encodeableBuffer, binary.LittleEndian, s.senderNonce)
		s.log.Debug().Err(err).Msg("write")
		encodeablePassword := make([]byte, 4+len(passwordBytes)+len(s.senderNonce))
		n, err := encodeableBuffer.Read(encodeablePassword)
		s.log.Debug().Err(err).Int("n", n).Msg("read")
		encryptedPassword, err := s.encryptionHandler.encryptPassword(encodeablePassword)
		if err != nil {
			s.log.Error().Err(err).Msg("error encrypting password")
			return nil
		}
		userNameIdentityToken := readWriteModel.NewUserNameIdentityToken(
			readWriteModel.NewPascalString(s.username),
			readWriteModel.NewPascalByteString(int32(len(encryptedPassword)), encryptedPassword),
			readWriteModel.NewPascalString(PASSWORD_ENCRYPTION_ALGORITHM),
		)
		extExpandedNodeId := readWriteModel.NewExpandedNodeId(
			false, //Namespace Uri Specified
			false, //Server Index Specified
			readWriteModel.NewNodeIdFourByte(0, uint16(readWriteModel.OpcuaNodeIdServices_UserNameIdentityToken_Encoding_DefaultBinary)),
			nil,
			nil)
		return readWriteModel.NewExtensionObject(
			extExpandedNodeId,
			readWriteModel.NewExtensionObjectEncodingMask(false, false, true),
			readWriteModel.NewUserIdentityToken(readWriteModel.NewPascalString(s.securityPolicy), userNameIdentityToken),
			false,
		)
	}
	return nil
}

func (s *SecureChannel) getCurrentDateTime() int64 {
	return (time.Now().UnixMilli() * 10000) + EPOCH_OFFSET
}
