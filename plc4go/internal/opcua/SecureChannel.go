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
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/pkg/api"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/opcua/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

const (
	VERSION                       = uint32(0)
	DEFAULT_MAX_CHUNK_COUNT       = 64
	DEFAULT_MAX_MESSAGE_SIZE      = uint32(2097152)
	DEFAULT_RECEIVE_BUFFER_SIZE   = uint32(65535)
	DEFAULT_SEND_BUFFER_SIZE      = uint32(65535)
	REQUEST_TIMEOUT               = 10 * time.Second
	REQUEST_TIMEOUT_LONG          = 10000
	PASSWORD_ENCRYPTION_ALGORITHM = "http://www.w3.org/2001/04/xmlenc#rsa-oaep"
	EPOCH_OFFSET                  = 116444736000000000 //Offset between OPC UA epoch time and linux epoch time.
)

var (
	SECURITY_POLICY_NONE = readWriteModel.NewPascalString(utils.ToPtr("http://opcfoundation.org/UA/SecurityPolicy#None"))
	NULL_STRING          = readWriteModel.NewPascalString(nil)
	NULL_BYTE_STRING     = readWriteModel.NewPascalByteString(-1, nil)
	NULL_EXPANDED_NODEID = readWriteModel.NewExpandedNodeId(false,
		false,
		readWriteModel.NewNodeIdTwoByte(0),
		nil,
		nil,
	)
	BINARY_ENCODING_MASK  = readWriteModel.NewExtensionObjectEncodingMask(false, false, true)
	NULL_EXTENSION_OBJECT = readWriteModel.NewNullExtensionObjectWithMask(NULL_EXPANDED_NODEID,
		readWriteModel.NewExtensionObjectEncodingMask(false, false, false),
		0,
		false) // Body

	INET_ADDRESS_PATTERN = regexp.MustCompile(`(.(?P<transportCode>tcp))?://(?P<transportHost>[\w.-]+)(:(?P<transportPort>\d*))?`)

	URI_PATTERN                 = regexp.MustCompile(`^(?P<protocolCode>opc)` + INET_ADDRESS_PATTERN.String() + `(?P<transportEndpoint>[\w/=]*)[?]?`)
	APPLICATION_URI             = readWriteModel.NewPascalString(utils.ToPtr("urn:apache:plc4x:client"))
	PRODUCT_URI                 = readWriteModel.NewPascalString(utils.ToPtr("urn:apache:plc4x:client"))
	APPLICATION_TEXT            = readWriteModel.NewPascalString(utils.ToPtr("OPCUA client for the Apache PLC4X:PLC4J project"))
	DEFAULT_CONNECTION_LIFETIME = uint32(36000000)
)

//go:generate go tool plc4xGenerator -type=SecureChannel
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
	channelId                 atomic.Uint32
	tokenId                   atomic.Uint32
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

	log zerolog.Logger
}

func NewSecureChannel(log zerolog.Logger, ctx DriverContext, configuration Configuration) *SecureChannel {
	s := &SecureChannel{
		configuration:             configuration,
		endpoint:                  readWriteModel.NewPascalString(&configuration.Endpoint),
		username:                  configuration.Username,
		password:                  configuration.Password,
		securityPolicy:            "http://opcfoundation.org/UA/SecurityPolicy#" + configuration.SecurityPolicy,
		sessionName:               "UaSession:" + *APPLICATION_TEXT.GetStringValue() + ":" + utils.RandomString(20),
		authenticationToken:       readWriteModel.NewNodeIdTwoByte(0),
		clientNonce:               []byte(utils.RandomString(40)),
		keyStoreFile:              configuration.KeyStoreFile,
		channelTransactionManager: NewSecureChannelTransactionManager(log),
		lifetime:                  DEFAULT_CONNECTION_LIFETIME,
		log:                       log,
	}
	s.requestHandleGenerator.Store(1)
	s.channelId.Store(1)
	s.tokenId.Store(1)
	ckp := configuration.Ckp
	if configuration.SecurityPolicy == "Basic256Sha256" {
		//Sender Certificate gets populated during the 'discover' phase when encryption is enabled.
		s.senderCertificate = configuration.SenderCertificate
		s.encryptionHandler = NewEncryptionHandler(s.log, ckp, s.senderCertificate, configuration.SecurityPolicy)
		certificate := ckp.getCertificate()
		s.publicCertificate = readWriteModel.NewPascalByteString(int32(len(certificate.Raw)), certificate.Raw)
		s.isEncrypted = true

		s.thumbprint = configuration.Thumbprint
	} else {
		s.encryptionHandler = NewEncryptionHandler(s.log, ckp, s.senderCertificate, configuration.SecurityPolicy)
		s.publicCertificate = NULL_BYTE_STRING
		s.thumbprint = NULL_BYTE_STRING
		s.isEncrypted = false
	}

	// Generate a list of endpoints we can use.
	{
		var err error
		address, err := url.Parse("none://" + configuration.Host)
		if err == nil {
			if names, lookupErr := net.LookupHost(address.Host); lookupErr == nil {
				s.endpoints = append(s.endpoints, names[rand.Intn(len(names))])
				s.endpoints = append(s.endpoints, address.Host)
				//s.endpoints = append(s.endpoints, address.Host)//TODO: not sure if golang can do
			} else {
				err = lookupErr
			}
		}
		if err != nil {
			s.log.Warn().Err(err).Msg("Unable to resolve host name. Using original host from connection string which may cause issues connecting to server")
			s.endpoints = append(s.endpoints, address.Host)
		}
	}

	s.channelId.Store(1)
	return s
}

func (s *SecureChannel) submit(ctx context.Context, codec *MessageCodec, errorDispatcher func(err error), consumer func(opcuaResponse []byte), buffer utils.WriteBufferByteBased) {
	transactionId := s.channelTransactionManager.getTransactionIdentifier()

	//TODO: We need to split large messages up into chunks if it is larger than the sendBufferSize
	//      This value is negotiated when opening a channel

	messageRequest := readWriteModel.NewOpcuaMessageRequest(
		readWriteModel.ChunkType_FINAL,
		readWriteModel.NewSecurityHeader(
			s.channelId.Load(),
			s.tokenId.Load(),
		),
		readWriteModel.NewBinaryPayload(
			readWriteModel.NewSequenceHeader(transactionId, transactionId),
			buffer.GetBytes(),
			uint32(len(buffer.GetBytes())),
		),
		uint32(len(buffer.GetBytes())),
		true,
	)

	var apu readWriteModel.OpcuaAPU
	if s.isEncrypted {
		message, err := s.encryptionHandler.encodeMessage(ctx, messageRequest, buffer.GetBytes())
		if err != nil {
			errorDispatcher(err)
			return
		}
		apu, err = readWriteModel.OpcuaAPUParse(ctx, message, false, true)
		if err != nil {
			errorDispatcher(err)
			return
		}
	} else {
		apu = readWriteModel.NewOpcuaAPU(messageRequest, false, true)
	}

	requestConsumer := func(transactionId int32) {
		var messageBuffer []byte
		if err := codec.SendRequest(ctx, apu,
			func(message spi.Message) bool {
				opcuaAPU, ok := message.(readWriteModel.OpcuaAPU)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				if decodedOpcuaAPU, err := s.encryptionHandler.decodeMessage(ctx, opcuaAPU); err != nil {
					s.log.Debug().Err(err).Msg("error decoding")
					return false
				} else {
					opcuaAPU = decodedOpcuaAPU.(readWriteModel.OpcuaAPU)
				}
				messagePDU := opcuaAPU.GetMessage()
				s.log.Trace().Stringer("messagePDU", messagePDU).Msg("looking at messagePDU")
				opcuaResponse, ok := messagePDU.(readWriteModel.OpcuaMessageResponse)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				if requestId := opcuaResponse.GetMessage().GetSequenceHeader().GetRequestId(); requestId != transactionId {
					s.log.Debug().Int32("requestId", requestId).Int32("transactionId", transactionId).Msg("Not relevant")
					return false
				} else {
					messageBuffer = opcuaResponse.(readWriteModel.BinaryPayload).GetPayload()
					if !(s.senderSequenceNumber.Add(1) == (opcuaResponse.GetMessage().GetSequenceHeader().GetSequenceNumber())) {
						s.log.Error().
							Int32("senderSequenceNumber", s.senderSequenceNumber.Load()).
							Int32("responseSequenceNumber", opcuaResponse.GetMessage().GetSequenceHeader().GetSequenceNumber()).
							Msg("Sequence number isn't as expected, we might have missed a packet. - senderSequenceNumber != responseSequenceNumber")
						errorDispatcher(errors.New("unexpected sequence number"))
					}
				}
				return true
			},
			func(message spi.Message) error {
				opcuaAPU := message.(readWriteModel.OpcuaAPU)
				opcuaAPU, _ = s.encryptionHandler.decodeMessage(ctx, opcuaAPU)
				messagePDU := opcuaAPU.GetMessage()
				s.log.Trace().Stringer("messagePDU", messagePDU).Msg("looking at messagePDU")
				opcuaResponse := messagePDU.(readWriteModel.OpcuaMessageResponse)
				if opcuaResponse.GetChunk() == (readWriteModel.ChunkType_FINAL) {
					s.tokenId.Store(opcuaResponse.GetSecurityHeader().GetSecureTokenId())
					s.channelId.Store(opcuaResponse.GetSecurityHeader().GetSecureChannelId())

					consumer(messageBuffer)
				} else {
					s.log.Warn().Stringer("chunk", opcuaResponse.GetChunk()).Msg("Message discarded")
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

	s.log.Debug().Int32("transactionId", transactionId).Msg("Submitting Transaction to TransactionManager")
	if err := s.channelTransactionManager.submit(requestConsumer, transactionId); err != nil {
		s.log.Debug().Err(err).Msg("error submitting")
	}
}

func (s *SecureChannel) onConnect(ctx context.Context, connection *Connection, ch chan plc4go.PlcConnectionConnectResult) {
	s.log.Trace().Msg("on connect")
	// Only the TCP transport supports login.
	s.log.Debug().Msg("Opcua Driver running in ACTIVE mode.")
	s.codec = connection.messageCodec // TODO: why would we need to set that?

	hello := readWriteModel.NewOpcuaHelloRequest(
		readWriteModel.ChunkType_FINAL,
		VERSION,
		readWriteModel.NewOpcuaProtocolLimits(
			DEFAULT_RECEIVE_BUFFER_SIZE,
			DEFAULT_SEND_BUFFER_SIZE,
			DEFAULT_MAX_MESSAGE_SIZE,
			DEFAULT_MAX_CHUNK_COUNT,
		),
		s.endpoint,
		true,
	)

	requestConsumer := func(transactionId int32) {
		s.log.Trace().Int32("transactionId", transactionId).Msg("request consumer called")
		if err := s.codec.SendRequest(
			ctx,
			hello,
			func(message spi.Message) bool {
				opcuaAPU, ok := message.(readWriteModel.OpcuaAPU)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				messagePDU := opcuaAPU.GetMessage()
				_, ok = messagePDU.(readWriteModel.OpcuaAcknowledgeResponse)
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
				go s.onConnectOpenSecureChannel(ctx, connection, ch, opcuaAcknowledgeResponse)
				return nil
			},
			func(err error) error {
				s.log.Debug().Err(err).Msg("error submitting")
				connection.fireConnectionError(err, ch)
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

func (s *SecureChannel) onConnectOpenSecureChannel(ctx context.Context, connection *Connection, ch chan plc4go.PlcConnectionConnectResult, response readWriteModel.OpcuaAcknowledgeResponse) {
	transactionId := s.channelTransactionManager.getTransactionIdentifier()

	requestHeader := readWriteModel.NewRequestHeader(
		s.getAuthenticationToken(),
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
			s.lifetime)
	} else {
		openSecureChannelRequest = readWriteModel.NewOpenSecureChannelRequest(
			requestHeader,
			VERSION,
			readWriteModel.SecurityTokenRequestType_securityTokenRequestTypeIssue,
			readWriteModel.MessageSecurityMode_messageSecurityModeNone,
			NULL_BYTE_STRING,
			s.lifetime)
	}

	identifier := openSecureChannelRequest.GetExtensionId()
	expandedNodeId := readWriteModel.NewExpandedNodeId(
		false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil,
	)

	extObject := readWriteModel.NewRootExtensionObject(
		expandedNodeId,
		openSecureChannelRequest,
		identifier,
	)

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		s.log.Debug().Err(err).Msg("error serializing")
		connection.fireConnectionError(err, ch)
		return
	}

	openRequest := readWriteModel.NewOpcuaOpenRequest(
		readWriteModel.ChunkType_FINAL,
		readWriteModel.NewOpenChannelMessageRequest(
			0,
			readWriteModel.NewPascalString(&s.securityPolicy),
			s.publicCertificate,
			s.thumbprint),
		readWriteModel.NewBinaryPayload(
			readWriteModel.NewSequenceHeader(transactionId, transactionId),
			buffer.GetBytes(),
			uint32(len(buffer.GetBytes())),
		),
		uint32(len(buffer.GetBytes())),
		true,
	)

	var apu readWriteModel.OpcuaAPU

	if s.isEncrypted {
		message, err := s.encryptionHandler.encodeMessage(ctx, openRequest, buffer.GetBytes())
		if err != nil {
			s.log.Debug().Err(err).Msg("error encoding")
			connection.fireConnectionError(err, ch)
			return
		}
		apu, err = readWriteModel.OpcuaAPUParse(ctx, message, false, true)
		if err != nil {
			s.log.Debug().Err(err).Msg("error parsing")
			connection.fireConnectionError(err, ch)
			return
		}
	} else {
		apu = readWriteModel.NewOpcuaAPU(openRequest, false, true)
	}

	requestConsumer := func(transactionId int32) {
		if err := s.codec.SendRequest(
			ctx,
			apu,
			func(message spi.Message) bool {
				opcuaAPU, ok := message.(readWriteModel.OpcuaAPU)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				messagePDU := opcuaAPU.GetMessage()
				openResponse, ok := messagePDU.(readWriteModel.OpcuaOpenResponse)
				if !ok {
					s.log.Debug().Type("type", messagePDU).Msg("Not relevant")
					return false
				}
				return openResponse.GetMessage().GetSequenceHeader().GetRequestId() == transactionId
			},
			func(message spi.Message) error {
				opcuaAPU := message.(readWriteModel.OpcuaAPU)
				messagePDU := opcuaAPU.GetMessage()
				opcuaOpenResponse := messagePDU.(readWriteModel.OpcuaOpenResponse)
				readBuffer := utils.NewReadBufferByteBased(opcuaOpenResponse.(readWriteModel.BinaryPayload).GetPayload(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
				extensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer[readWriteModel.ExtensionObject](ctx, readBuffer, false)
				if err != nil {
					return errors.Wrap(err, "error parsing")
				}
				//Store the initial sequence number from the server. there's no requirement for the server and client to use the same starting number.
				s.senderSequenceNumber.Store(opcuaOpenResponse.GetMessage().GetSequenceHeader().GetSequenceNumber())

				if fault, ok := extensionObject.GetBody().(readWriteModel.ServiceFault); ok {
					statusCode := fault.GetResponseHeader().(readWriteModel.ResponseHeader).GetServiceResult().GetStatusCode()
					statusCodeByValue, _ := readWriteModel.OpcuaStatusCodeByValue(statusCode)
					s.log.Error().
						Uint32("statusCode", statusCode).
						Stringer("statusCodeByValue", statusCodeByValue).
						Msg("Failed to connect to opc ua server for the following reason")
					connection.fireConnectionError(errors.New("service fault received"), ch)
					return nil
				}
				s.log.Debug().Msg("Got Secure Response Connection Response")
				openSecureChannelResponse := extensionObject.GetBody().(readWriteModel.OpenSecureChannelResponse)
				s.tokenId.Store(openSecureChannelResponse.GetSecurityToken().(readWriteModel.ChannelSecurityToken).GetTokenId())
				s.channelId.Store(openSecureChannelResponse.GetSecurityToken().(readWriteModel.ChannelSecurityToken).GetChannelId())
				go s.onConnectCreateSessionRequest(ctx, connection, ch)
				return nil
			},
			func(err error) error {
				s.log.Debug().Err(err).Msg("error submitting")
				connection.fireConnectionError(err, ch)
				return nil
			},
			REQUEST_TIMEOUT,
		); err != nil {
			s.log.Debug().Err(err).Msg("a error")
			connection.fireConnectionError(err, ch)
		}
	}
	s.log.Debug().Int32("transactionId", transactionId).Msg("Submitting OpenSecureChannel with id")
	if err := s.channelTransactionManager.submit(requestConsumer, transactionId); err != nil {
		s.log.Debug().Err(err).Msg("error submitting")
		connection.fireConnectionError(err, ch)
	}
}

func (s *SecureChannel) onConnectCreateSessionRequest(ctx context.Context, connection *Connection, ch chan plc4go.PlcConnectionConnectResult) {
	requestHeader := readWriteModel.NewRequestHeader(
		s.getAuthenticationToken(),
		s.getCurrentDateTime(),
		0,
		0,
		NULL_STRING,
		REQUEST_TIMEOUT_LONG,
		NULL_EXTENSION_OBJECT)

	applicationName := readWriteModel.NewLocalizedText(
		true,
		true,
		readWriteModel.NewPascalString(utils.ToPtr("en")),
		APPLICATION_TEXT)

	var discoveryUrls []readWriteModel.PascalString

	clientDescription := readWriteModel.NewApplicationDescription(APPLICATION_URI,
		PRODUCT_URI,
		applicationName,
		readWriteModel.ApplicationType_applicationTypeClient,
		NULL_STRING,
		NULL_STRING,
		discoveryUrls)

	createSessionRequest := readWriteModel.NewCreateSessionRequest(
		requestHeader,
		clientDescription,
		NULL_STRING,
		s.endpoint,
		readWriteModel.NewPascalString(&s.sessionName),
		readWriteModel.NewPascalByteString(int32(len(s.clientNonce)), s.clientNonce),
		NULL_BYTE_STRING,
		120000,
		0,
	)

	identifier := createSessionRequest.GetExtensionId()
	expandedNodeId := readWriteModel.NewExpandedNodeId(
		false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil)

	extObject := readWriteModel.NewRootExtensionObject(
		expandedNodeId,
		createSessionRequest,
		identifier,
	)

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		s.log.Debug().Err(err).Msg("error serializing")
		connection.fireConnectionError(err, ch)
		return
	}

	consumer := func(opcuaResponse []byte) {
		extensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer[readWriteModel.ExtensionObject](ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
		if err != nil {
			s.log.Error().Err(err).Msg("error parsing")
			connection.fireConnectionError(err, ch)
			return
		}
		s.log.Trace().Stringer("extensionObject", extensionObject).Msg("looking at message")
		if fault, ok := extensionObject.GetBody().(readWriteModel.ServiceFault); ok {
			statusCode := fault.GetResponseHeader().(readWriteModel.ResponseHeader).GetServiceResult().GetStatusCode()
			statusCodeByValue, _ := readWriteModel.OpcuaStatusCodeByValue(statusCode)
			s.log.Error().
				Uint32("statusCode", statusCode).
				Stringer("statusCodeByValue", statusCodeByValue).
				Msg("Failed to connect to opc ua server for the following reason")
			connection.fireConnectionError(errors.New("service fault received"), ch)
			return
		}
		s.log.Debug().Msg("Got Create Session Response Connection Response")

		unknownExtensionObject := extensionObject.GetBody()
		if responseMessage, ok := unknownExtensionObject.(readWriteModel.CreateSessionResponse); ok {
			s.authenticationToken = responseMessage.GetAuthenticationToken().GetNodeId()

			go s.onConnectActivateSessionRequest(ctx, connection, ch, responseMessage, responseMessage)
		} else {
			serviceFault := unknownExtensionObject.(readWriteModel.ServiceFault)
			header := serviceFault.GetResponseHeader().(readWriteModel.ResponseHeader)
			s.log.Error().
				Stringer("serviceResult", header.GetServiceResult()).
				Msg("Subscription ServiceFault returned from server with error code, '%s'")
		}
	}

	errorDispatcher := func(err error) {
		s.log.Error().Err(err).Msg("Error while waiting for subscription response")
		connection.fireConnectionError(err, ch)
	}

	s.submit(ctx, connection.messageCodec, errorDispatcher, consumer, buffer)
}

func (s *SecureChannel) onConnectActivateSessionRequest(ctx context.Context, connection *Connection, ch chan plc4go.PlcConnectionConnectResult, opcuaMessageResponse readWriteModel.CreateSessionResponse, sessionResponse readWriteModel.CreateSessionResponse) {
	s.senderCertificate = sessionResponse.GetServerCertificate().GetStringValue()
	certificate, err := s.encryptionHandler.getCertificateX509(s.senderCertificate)
	if err != nil {
		s.log.Error().Err(err).Msg("error getting certificate")
		connection.fireConnectionError(err, ch)
		return
	}
	s.log.Debug().Interface("senderCertificate", certificate).Msg("working with senderCertificate")
	s.encryptionHandler.setServerCertificate(certificate)
	s.senderNonce = sessionResponse.GetServerNonce().GetStringValue()
	endpoints := make([]string, 3)
	if address, err := url.Parse(s.configuration.Host); err != nil {
		if names, err := net.LookupAddr(address.Host); err != nil {
			endpoints[0] = "opc.tcp://" + names[rand.Intn(len(names))] + ":" + s.configuration.Port + s.configuration.TransportEndpoint
		}
		endpoints[1] = "opc.tcp://" + address.Hostname() + ":" + s.configuration.Port + s.configuration.TransportEndpoint
		//endpoints[2] = "opc.tcp://" + address.getCanonicalHostName() + ":" + s.configuration.getPort() + s.configuration.transportEndpoint// TODO: not sure how to get that in golang
	}

	s.selectEndpoint(sessionResponse)

	if s.policyId == nil {
		s.log.Error().Msg("Unable to find endpoint - " + endpoints[1])
		connection.fireConnectionError(err, ch)
		return
	}

	userIdentityToken := s.getIdentityToken(s.tokenType, s.policyId.GetStringValue())

	requestHandle := s.getRequestHandle()

	requestHeader := readWriteModel.NewRequestHeader(
		s.getAuthenticationToken(),
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
		nil,
		nil,
		userIdentityToken,
		clientSignature,
	)

	identifier := activateSessionRequest.GetExtensionId()
	expandedNodeId := readWriteModel.NewExpandedNodeId(false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil)

	extObject := readWriteModel.NewRootExtensionObject(
		expandedNodeId,
		activateSessionRequest,
		identifier,
	)

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		s.log.Debug().Err(err).Msg("error serializing")
		connection.fireConnectionError(err, ch)
		return
	}

	consumer := func(opcuaResponse []byte) {
		message, err := readWriteModel.ExtensionObjectParseWithBuffer[readWriteModel.ExtensionObject](ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
		if err != nil {
			s.log.Error().Err(err).Msg("error parsing")
			return
		}
		s.log.Trace().Stringer("message", message).Msg("looking at message")
		if fault, ok := message.GetBody().(readWriteModel.ServiceFault); ok {
			statusCode := fault.GetResponseHeader().(readWriteModel.ResponseHeader).GetServiceResult().GetStatusCode()
			statusCodeByValue, _ := readWriteModel.OpcuaStatusCodeByValue(statusCode)
			s.log.Error().
				Uint32("statusCode", statusCode).
				Stringer("statusCodeByValue", statusCodeByValue).
				Msg("Failed to connect to opc ua server for the following reason")
			connection.fireConnectionError(errors.New("service fault received"), ch)
			return
		}
		s.log.Debug().Msg("Got Activate Session Response Connection Response")

		extensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer[readWriteModel.ExtensionObject](ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
		if err != nil {
			s.log.Error().Err(err).Msg("error parsing")
			return
		}
		unknownExtensionObject := extensionObject.GetBody()
		if responseMessage, ok := unknownExtensionObject.(readWriteModel.ActivateSessionResponse); ok {
			returnedRequestHandle := responseMessage.GetResponseHeader().(readWriteModel.ResponseHeader).GetRequestHandle()
			if !(requestHandle == returnedRequestHandle) {
				s.log.Error().
					Uint32("requestHandle", requestHandle).
					Uint32("returnedRequestHandle", returnedRequestHandle).
					Msg("Request handle isn't as expected, we might have missed a packet. requestHandle != returnedRequestHandle")
			}

			// Send an event that connection setup is complete.
			s.keepAlive()
			connection.fireConnected(ch)
		} else {
			serviceFault := unknownExtensionObject.(readWriteModel.ServiceFault)
			header := serviceFault.GetResponseHeader().(readWriteModel.ResponseHeader)
			s.log.Error().
				Stringer("serviceResult", header.GetServiceResult()).
				Msg("Subscription ServiceFault returned from server with error code")
		}
	}

	errorDispatcher := func(err error) {
		s.log.Error().Err(err).Msg("Error while waiting for subscription response")
		connection.fireConnectionError(err, ch)
	}

	s.submit(ctx, connection.messageCodec, errorDispatcher, consumer, buffer)
}

func (s *SecureChannel) onDisconnect(ctx context.Context, connection *Connection) {
	s.log.Info().Msg("disconnecting")
	requestHandle := s.getRequestHandle()

	s.keepAliveIndicator.Store(false)

	expandedNodeId := readWriteModel.NewExpandedNodeId(false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, 473),
		nil,
		nil) //Identifier for OpenSecureChannel

	if s.authenticationToken == nil {
		// TODO: this or nil?? What do we do when we don't have one?
		s.log.Error().Msg("no authentication token, so we can't disconnect")
		return
	}
	requestHeader := readWriteModel.NewRequestHeader(
		s.getAuthenticationToken(),
		s.getCurrentDateTime(),
		requestHandle, //RequestHandle
		0,
		NULL_STRING,
		5000,
		NULL_EXTENSION_OBJECT)

	closeSessionRequest := readWriteModel.NewCloseSessionRequest(
		requestHeader,
		true)

	extObject := readWriteModel.NewRootExtensionObject(
		expandedNodeId,
		closeSessionRequest,
		closeSessionRequest.GetExtensionId(),
	)

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		s.log.Debug().Err(err).Msg("error serializing")
		return
	}

	consumer := func(opcuaResponse []byte) {
		message, err := readWriteModel.ExtensionObjectParseWithBuffer[readWriteModel.ExtensionObject](ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
		if err != nil {
			s.log.Error().Err(err).Msg("error parsing")
			return
		}
		s.log.Trace().Stringer("message", message).Msg("looking at message")
		if fault, ok := message.GetBody().(readWriteModel.ServiceFault); ok {
			statusCode := fault.GetResponseHeader().(readWriteModel.ResponseHeader).GetServiceResult().GetStatusCode()
			statusCodeByValue, _ := readWriteModel.OpcuaStatusCodeByValue(statusCode)
			s.log.Error().
				Uint32("statusCode", statusCode).
				Stringer("statusCodeByValue", statusCodeByValue).
				Msg("Failed to connect to opc ua server for the following reason")
			return
		}
		s.log.Debug().Msg("Got Close Session Response Connection Response")

		extensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer[readWriteModel.ExtensionObject](ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
		if err != nil {
			s.log.Error().Err(err).Msg("error parsing")
			return
		}
		unknownExtensionObject := extensionObject.GetBody()
		if responseMessage, ok := unknownExtensionObject.(readWriteModel.CloseSessionResponse); ok {
			go s.onDisconnectCloseSecureChannel(ctx, connection, responseMessage, message.GetBody().(readWriteModel.CloseSessionResponse))
		} else {
			serviceFault := unknownExtensionObject.(readWriteModel.ServiceFault)
			header := serviceFault.GetResponseHeader().(readWriteModel.ResponseHeader)
			s.log.Error().
				Stringer("serviceResult", header.GetServiceResult()).
				Msg("Subscription ServiceFault returned from server with error code")
		}
	}

	errorDispatcher := func(err error) {
		s.log.Error().Err(err).Msg("Error while waiting for close session response")
	}

	s.submit(ctx, connection.messageCodec, errorDispatcher, consumer, buffer)
}

func (s *SecureChannel) onDisconnectCloseSecureChannel(ctx context.Context, connection *Connection, message readWriteModel.CloseSessionResponse, response readWriteModel.CloseSessionResponse) {
	transactionId := s.channelTransactionManager.getTransactionIdentifier()

	requestHeader := readWriteModel.NewRequestHeader(
		s.getAuthenticationToken(),
		s.getCurrentDateTime(),
		0, //RequestHandle
		0,
		NULL_STRING,
		REQUEST_TIMEOUT_LONG,
		NULL_EXTENSION_OBJECT)

	closeSecureChannelRequest := readWriteModel.NewCloseSecureChannelRequest(requestHeader)

	identifier := closeSecureChannelRequest.GetExtensionId()
	expandedNodeId := readWriteModel.NewExpandedNodeId(
		false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil,
	)

	closeRequest := readWriteModel.NewOpcuaCloseRequest(
		readWriteModel.ChunkType_FINAL,
		readWriteModel.NewSecurityHeader(s.channelId.Load(), s.tokenId.Load()),
		readWriteModel.NewExtensiblePayload(
			readWriteModel.NewSequenceHeader(transactionId, transactionId),
			readWriteModel.NewRootExtensionObject(
				expandedNodeId,
				closeSecureChannelRequest,
				identifier,
			),
			0,
		),
		true,
	)

	apu := readWriteModel.NewOpcuaAPU(closeRequest, false, true)

	requestConsumer := func(transactionId int32) {
		if err := connection.messageCodec.SendRequest(
			ctx,
			apu,
			func(message spi.Message) bool {
				opcuaAPU, ok := message.(readWriteModel.OpcuaAPU)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				messagePDU := opcuaAPU.GetMessage()
				openResponse, ok := messagePDU.(readWriteModel.OpcuaMessageResponse)
				if !ok {
					s.log.Debug().Type("type", messagePDU).Msg("Not relevant")
					return false
				}
				return openResponse.GetMessage().GetSequenceHeader().GetRequestId() == transactionId
			},
			func(message spi.Message) error {
				opcuaAPU := message.(readWriteModel.OpcuaAPU)
				messagePDU := opcuaAPU.GetMessage()
				opcuaMessageResponse := messagePDU.(readWriteModel.OpcuaMessageResponse)
				s.log.Trace().Stringer("opcuaMessageResponse", opcuaMessageResponse).Msg("Got close secure channel response")
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
	s.log.Debug().Int32("transactionId", transactionId).Msg("Submitting CloseSecureChannel with id")
	if err := s.channelTransactionManager.submit(requestConsumer, transactionId); err != nil {
		s.log.Debug().Err(err).Msg("error submitting")
	}
}

func (s *SecureChannel) onDiscover(ctx context.Context, codec *MessageCodec) {
	s.log.Trace().Msg("onDiscover")
	// Only the TCP transport supports login.
	s.log.Debug().Msg("Opcua Driver running in ACTIVE mode, discovering endpoints")

	hello := readWriteModel.NewOpcuaHelloRequest(
		readWriteModel.ChunkType_FINAL,
		VERSION,
		readWriteModel.NewOpcuaProtocolLimits(
			DEFAULT_RECEIVE_BUFFER_SIZE,
			DEFAULT_SEND_BUFFER_SIZE,
			DEFAULT_MAX_MESSAGE_SIZE,
			DEFAULT_MAX_CHUNK_COUNT,
		),
		s.endpoint,
		true,
	)

	apu := readWriteModel.NewOpcuaAPU(hello, false, true)

	requestConsumer := func(transactionId int32) {
		if err := codec.SendRequest(
			ctx,
			apu,
			func(message spi.Message) bool {
				opcuaAPU, ok := message.(readWriteModel.OpcuaAPU)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				messagePDU := opcuaAPU.GetMessage()
				_, ok = messagePDU.(readWriteModel.OpcuaAcknowledgeResponse)
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
				s.log.Trace().Stringer("opcuaAcknowledgeResponse", opcuaAcknowledgeResponse).Msg("Got Hello Response Connection Response")
				go s.onDiscoverOpenSecureChannel(ctx, codec, opcuaAcknowledgeResponse)
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
	s.log.Trace().Msg("onDiscoverOpenSecureChannel")
	transactionId := s.channelTransactionManager.getTransactionIdentifier()

	requestHeader := readWriteModel.NewRequestHeader(
		s.getAuthenticationToken(),
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

	identifier := openSecureChannelRequest.GetExtensionId()
	expandedNodeId := readWriteModel.NewExpandedNodeId(
		false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil,
	)

	extObject := readWriteModel.NewRootExtensionObject(
		expandedNodeId,
		openSecureChannelRequest,
		identifier,
	)

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		s.log.Debug().Err(err).Msg("error serializing")
		return
	}

	openRequest := readWriteModel.NewOpcuaOpenRequest(
		readWriteModel.ChunkType_FINAL,
		readWriteModel.NewOpenChannelMessageRequest(
			0,
			SECURITY_POLICY_NONE,
			NULL_BYTE_STRING,
			NULL_BYTE_STRING,
		),
		readWriteModel.NewBinaryPayload(
			readWriteModel.NewSequenceHeader(transactionId, transactionId),
			buffer.GetBytes(),
			uint32(len(buffer.GetBytes())),
		),
		uint32(len(buffer.GetBytes())),
		true,
	)

	apu := readWriteModel.NewOpcuaAPU(openRequest, false, true)

	requestConsumer := func(transactionId int32) {
		if err := codec.SendRequest(
			ctx,
			apu,
			func(message spi.Message) bool {
				opcuaAPU, ok := message.(readWriteModel.OpcuaAPU)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				messagePDU := opcuaAPU.GetMessage()
				openResponse, ok := messagePDU.(readWriteModel.OpcuaOpenResponse)
				if !ok {
					s.log.Debug().Type("type", messagePDU).Msg("Not relevant")
					return false
				}
				return openResponse.GetMessage().GetSequenceHeader().GetRequestId() == transactionId
			},
			func(message spi.Message) error {
				opcuaAPU := message.(readWriteModel.OpcuaAPU)
				messagePDU := opcuaAPU.GetMessage()
				opcuaOpenResponse := messagePDU.(readWriteModel.OpcuaOpenResponse)
				readBuffer := utils.NewReadBufferByteBased(opcuaOpenResponse.(readWriteModel.BinaryPayload).GetPayload(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
				extensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer[readWriteModel.ExtensionObject](ctx, readBuffer, false)
				if err != nil {
					return errors.Wrap(err, "error parsing")
				}

				if fault, ok := extensionObject.GetBody().(readWriteModel.ServiceFault); ok {
					statusCode := fault.GetResponseHeader().(readWriteModel.ResponseHeader).GetServiceResult().GetStatusCode()
					statusCodeByValue, _ := readWriteModel.OpcuaStatusCodeByValue(statusCode)
					s.log.Error().
						Uint32("statusCode", statusCode).
						Stringer("statusCodeByValue", statusCodeByValue).
						Msg("Failed to connect to opc ua server for the following reason")
					return nil
				}
				s.log.Debug().Msg("Got Secure Response Connection Response")
				openSecureChannelResponse := extensionObject.GetBody().(readWriteModel.OpenSecureChannelResponse)
				go s.onDiscoverGetEndpointsRequest(ctx, codec, opcuaOpenResponse, openSecureChannelResponse)
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
	s.log.Trace().Msg("onDiscoverGetEndpointsRequest")
	s.tokenId.Store(openSecureChannelResponse.GetSecurityToken().(readWriteModel.ChannelSecurityToken).GetTokenId())
	s.channelId.Store(openSecureChannelResponse.GetSecurityToken().(readWriteModel.ChannelSecurityToken).GetChannelId())

	transactionId := s.channelTransactionManager.getTransactionIdentifier()

	nextSequenceNumber := opcuaOpenResponse.GetMessage().GetSequenceHeader().GetSequenceNumber() + 1
	nextRequestId := opcuaOpenResponse.GetMessage().GetSequenceHeader().GetRequestId() + 1

	if !(transactionId == nextSequenceNumber) {
		s.log.Error().
			Int32("transactionId", transactionId).
			Int32("nextSequenceNumber", nextSequenceNumber).
			Msg("Sequence number isn't as expected, we might have missed a packet. - transactionId != nextSequenceNumber")
		return
	}

	requestHeader := readWriteModel.NewRequestHeader(
		s.getAuthenticationToken(),
		s.getCurrentDateTime(),
		0, //RequestHandle
		0,
		NULL_STRING,
		REQUEST_TIMEOUT_LONG,
		NULL_EXTENSION_OBJECT)

	endpointsRequest := readWriteModel.NewGetEndpointsRequest(
		requestHeader,
		s.endpoint,
		nil,
		nil)

	identifier := endpointsRequest.GetExtensionId()
	expandedNodeId := readWriteModel.NewExpandedNodeId(
		false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil,
	)

	extObject := readWriteModel.NewRootExtensionObject(
		expandedNodeId,
		endpointsRequest,
		identifier,
	)

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		s.log.Debug().Err(err).Msg("error serializing")
		return
	}

	messageRequest := readWriteModel.NewOpcuaMessageRequest(
		readWriteModel.ChunkType_FINAL,
		readWriteModel.NewSecurityHeader(
			s.channelId.Load(),
			s.tokenId.Load(),
		),
		readWriteModel.NewBinaryPayload(
			readWriteModel.NewSequenceHeader(nextSequenceNumber, nextRequestId),
			buffer.GetBytes(),
			uint32(len(buffer.GetBytes())),
		),
		uint32(len(buffer.GetBytes())),
		true,
	)

	apu := readWriteModel.NewOpcuaAPU(messageRequest, false, true)

	requestConsumer := func(transactionId int32) {
		if err := codec.SendRequest(
			ctx,
			apu,
			func(message spi.Message) bool {
				opcuaAPU, ok := message.(readWriteModel.OpcuaAPU)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				messagePDU := opcuaAPU.GetMessage()
				messageResponse, ok := messagePDU.(readWriteModel.OpcuaMessageResponse)
				if !ok {
					s.log.Debug().Type("type", messagePDU).Msg("Not relevant")
					return false
				}
				return messageResponse.GetMessage().GetSequenceHeader().GetRequestId() == transactionId
			},
			func(message spi.Message) error {
				opcuaAPU := message.(readWriteModel.OpcuaAPU)
				messagePDU := opcuaAPU.GetMessage()
				messageResponse := messagePDU.(readWriteModel.OpcuaMessageResponse)
				readBuffer := utils.NewReadBufferByteBased(messageResponse.(readWriteModel.BinaryPayload).GetPayload(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
				extensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer[readWriteModel.ExtensionObject](ctx, readBuffer, false)
				if err != nil {
					return errors.Wrap(err, "error parsing")
				}

				if fault, ok := extensionObject.GetBody().(readWriteModel.ServiceFault); ok {
					statusCode := fault.GetResponseHeader().(readWriteModel.ResponseHeader).GetServiceResult().GetStatusCode()
					statusCodeByValue, _ := readWriteModel.OpcuaStatusCodeByValue(statusCode)
					s.log.Error().
						Uint32("statusCode", statusCode).
						Stringer("statusCodeByValue", statusCodeByValue).
						Msg("Failed to connect to opc ua server for the following reason")
				} else {
					s.log.Debug().Msg("Got Secure Response Connection Response")
					response := extensionObject.GetBody().(readWriteModel.GetEndpointsResponse)

					endpoints := response.GetEndpoints()
					for _, endpoint := range endpoints {
						endpointDescription := endpoint.(readWriteModel.EndpointDescription)
						if endpointDescription.GetEndpointUrl().GetStringValue() == (s.endpoint.GetStringValue()) && *endpointDescription.GetSecurityPolicyUri().GetStringValue() == (s.securityPolicy) {
							s.log.Info().Str("stringValue", *s.endpoint.GetStringValue()).Msg("Found OPC UA endpoint")
							s.configuration.SenderCertificate = endpointDescription.GetServerCertificate().GetStringValue()
						}
					}

					digest := sha1.Sum(s.configuration.SenderCertificate)
					s.thumbprint = readWriteModel.NewPascalByteString(int32(len(digest)), digest[:])

					go s.onDiscoverCloseSecureChannel(ctx, codec, response)
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
	s.log.Trace().Msg("onDiscoverCloseSecureChannel")
	transactionId := s.channelTransactionManager.getTransactionIdentifier()

	requestHeader := readWriteModel.NewRequestHeader(
		s.getAuthenticationToken(),
		s.getCurrentDateTime(),
		0, //RequestHandle
		0,
		NULL_STRING,
		REQUEST_TIMEOUT_LONG,
		NULL_EXTENSION_OBJECT)

	closeSecureChannelRequest := readWriteModel.NewCloseSecureChannelRequest(requestHeader)

	identifier := closeSecureChannelRequest.GetExtensionId()
	expandedNodeId := readWriteModel.NewExpandedNodeId(
		false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil,
	)

	closeRequest := readWriteModel.NewOpcuaCloseRequest(
		readWriteModel.ChunkType_FINAL,
		readWriteModel.NewSecurityHeader(
			s.channelId.Load(),
			s.tokenId.Load(),
		),
		readWriteModel.NewExtensiblePayload(
			readWriteModel.NewSequenceHeader(transactionId, transactionId),
			readWriteModel.NewRootExtensionObject(
				expandedNodeId,
				closeSecureChannelRequest,
				identifier,
			),
			uint32(0),
		),
		true,
	)

	apu := readWriteModel.NewOpcuaAPU(closeRequest, false, true)

	requestConsumer := func(transactionId int32) {
		if err := codec.SendRequest(
			ctx,
			apu,
			func(message spi.Message) bool {
				opcuaAPU, ok := message.(readWriteModel.OpcuaAPU)
				if !ok {
					s.log.Debug().Type("type", message).Msg("Not relevant")
					return false
				}
				messagePDU := opcuaAPU.GetMessage()
				openResponse, ok := messagePDU.(readWriteModel.OpcuaMessageResponse)
				if !ok {
					s.log.Debug().Type("type", messagePDU).Msg("Not relevant")
					return false
				}
				return openResponse.GetMessage().GetSequenceHeader().GetRequestId() == transactionId
			},
			func(message spi.Message) error {
				opcuaAPU := message.(readWriteModel.OpcuaAPU)
				messagePDU := opcuaAPU.GetMessage()
				opcuaMessageResponse := messagePDU.(readWriteModel.OpcuaMessageResponse)
				s.log.Trace().Stringer("opcuaMessageResponse", opcuaMessageResponse).Msg("Got close secure channel response")
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
	s.log.Debug().Int32("transactionId", transactionId).Msg("Submitting CloseSecureChannel with id")
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

			requestHeader := readWriteModel.NewRequestHeader(
				s.getAuthenticationToken(),
				s.getCurrentDateTime(),
				0, //RequestHandle
				0,
				NULL_STRING,
				REQUEST_TIMEOUT_LONG,
				NULL_EXTENSION_OBJECT,
			)

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
			identifier := openSecureChannelRequest.GetExtensionId()
			expandedNodeId := readWriteModel.NewExpandedNodeId(false, //Namespace Uri Specified
				false, //Server Index Specified
				readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
				nil,
				nil)

			extObject := readWriteModel.NewRootExtensionObject(
				expandedNodeId,
				openSecureChannelRequest,
				identifier,
			)

			buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
			if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
				s.log.Error().Err(err).Msg("error serializing")
				return
			}

			openRequest := readWriteModel.NewOpcuaOpenRequest(
				readWriteModel.ChunkType_FINAL,
				readWriteModel.NewOpenChannelMessageRequest(0,
					readWriteModel.NewPascalString(&s.securityPolicy),
					s.publicCertificate,
					s.thumbprint,
				),
				readWriteModel.NewBinaryPayload(
					readWriteModel.NewSequenceHeader(transactionId, transactionId),
					buffer.GetBytes(),
					uint32(len(buffer.GetBytes())),
				),
				uint32(len(buffer.GetBytes())),
				true,
			)

			var apu readWriteModel.OpcuaAPU

			if s.isEncrypted {
				message, err := s.encryptionHandler.encodeMessage(ctx, openRequest, buffer.GetBytes())
				if err != nil {
					s.log.Error().Err(err).Msg("error encoding")
					return
				}
				apu, err = readWriteModel.OpcuaAPUParse(ctx, message, false, true)
				if err != nil {
					s.log.Error().Err(err).Msg("error parsing")
					return
				}
			} else {
				apu = readWriteModel.NewOpcuaAPU(openRequest, false, true)
			}

			requestConsumer := func(transactionId int32) {
				if err := s.codec.SendRequest(
					ctx,
					apu,
					func(message spi.Message) bool {
						opcuaAPU, ok := message.(readWriteModel.OpcuaAPU)
						if !ok {
							s.log.Debug().Type("type", message).Msg("Not relevant")
							return false
						}
						messagePDU := opcuaAPU.GetMessage()
						openResponse, ok := messagePDU.(readWriteModel.OpcuaOpenResponse)
						if !ok {
							s.log.Debug().Type("type", messagePDU).Msg("Not relevant")
							return false
						}
						return openResponse.GetMessage().GetSequenceHeader().GetRequestId() == transactionId
					},
					func(message spi.Message) error {
						opcuaAPU := message.(readWriteModel.OpcuaAPU)
						messagePDU := opcuaAPU.GetMessage()
						opcuaOpenResponse := messagePDU.(readWriteModel.OpcuaOpenResponse)
						readBuffer := utils.NewReadBufferByteBased(opcuaOpenResponse.(readWriteModel.BinaryPayload).GetPayload(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
						extensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer[readWriteModel.ExtensionObject](ctx, readBuffer, false)
						if err != nil {
							return errors.Wrap(err, "error parsing")
						}

						if fault, ok := extensionObject.GetBody().(readWriteModel.ServiceFault); ok {
							statusCode := fault.GetResponseHeader().(readWriteModel.ResponseHeader).GetServiceResult().GetStatusCode()
							statusCodeByValue, _ := readWriteModel.OpcuaStatusCodeByValue(statusCode)
							s.log.Error().
								Uint32("statusCode", statusCode).
								Stringer("statusCodeByValue", statusCodeByValue).
								Msg("Failed to connect to opc ua server for the following reason")
						} else {
							s.log.Debug().Msg("Got Secure Response Connection Response")
							openSecureChannelResponse := extensionObject.GetBody().(readWriteModel.OpenSecureChannelResponse)
							token := openSecureChannelResponse.GetSecurityToken().(readWriteModel.ChannelSecurityToken)
							s.tokenId.Store(token.GetTokenId())
							s.channelId.Store(token.GetChannelId())
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
			s.log.Debug().Int32("transactionId", transactionId).Msg("Submitting OpenSecureChannel with id")
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
	if s.authenticationToken == nil {
		panic("authenticationToken should be set at this point")
	}
	return readWriteModel.NewNodeId(s.authenticationToken)
}

// getChannelId gets the Channel identifier for the current channel
func (s *SecureChannel) getChannelId() uint32 {
	return s.channelId.Load()
}

// getTokenId gets the Token Identifier
func (s *SecureChannel) getTokenId() uint32 {
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
		s.log.Error().Str("endpoint", s.endpoints[0]).Msg("Unable to find endpoint")
		return
	}

	if s.tokenType == 0xffffffff { // TODO: what did we use as undefined
		s.log.Error().Str("endpoint", s.endpoints[0]).Msg("Unable to find Security Policy for endpoint")
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
	matches := utils.GetSubgroupMatches(URI_PATTERN, *endpoint.GetEndpointUrl().GetStringValue())
	if len(matches) == 0 {
		s.log.Error().Stringer("endpoint", endpoint).Msg("Endpoint returned from the server doesn't match the format '{protocol-code}:({transport-code})?//{transport-host}(:{transport-port})(/{transport-endpoint})'")
		return false
	}
	s.log.Trace().
		Str("transportHost", matches["transportHost"]).
		Str("transportPort", matches["transportPort"]).
		Str("transportEndpoint", matches["transportEndpoint"]).
		Msg("Using Endpoint")

	if s.configuration.Discovery && !slices.Contains(s.endpoints, matches["transportHost"]) {
		return false
	}

	if s.configuration.Port != matches["transportPort"] {
		return false
	}

	if s.configuration.TransportEndpoint != matches["transportEndpoint"] {
		return false
	}

	if !s.configuration.Discovery {
		s.configuration.Host = matches["transportHost"]
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
//   - @param policyId 	 	 the policy id
//   - @return returns an ExtensionObject with an IdentityToken.
func (s *SecureChannel) getIdentityToken(tokenType readWriteModel.UserTokenType, policyId *string) readWriteModel.ExtensionObject {
	switch tokenType {
	case readWriteModel.UserTokenType_userTokenTypeAnonymous:
		//If we aren't using authentication tell the server we would like to log in anonymously
		anonymousIdentityToken := readWriteModel.NewAnonymousIdentityToken(readWriteModel.NewPascalString(policyId))
		extExpandedNodeId := readWriteModel.NewExpandedNodeId(
			false, //Namespace Uri Specified
			false, //Server Index Specified
			readWriteModel.NewNodeIdFourByte(0, uint16(anonymousIdentityToken.GetExtensionId())),
			nil,
			nil,
		)
		return readWriteModel.NewBinaryExtensionObjectWithMask(
			extExpandedNodeId,
			BINARY_ENCODING_MASK,
			anonymousIdentityToken,
			anonymousIdentityToken.GetExtensionId(),
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
			readWriteModel.NewPascalString(policyId),
			readWriteModel.NewPascalString(&s.username),
			readWriteModel.NewPascalByteString(int32(len(encryptedPassword)), encryptedPassword),
			readWriteModel.NewPascalString(utils.ToPtr(PASSWORD_ENCRYPTION_ALGORITHM)),
		)
		extExpandedNodeId := readWriteModel.NewExpandedNodeId(
			false, //Namespace Uri Specified
			false, //Server Index Specified
			readWriteModel.NewNodeIdFourByte(0, uint16(userNameIdentityToken.GetExtensionId())),
			nil,
			nil)
		return readWriteModel.NewBinaryExtensionObjectWithMask(
			extExpandedNodeId,
			BINARY_ENCODING_MASK,
			userNameIdentityToken,
			userNameIdentityToken.GetExtensionId(),
			false,
		)
	}
	return nil
}

func (s *SecureChannel) getCurrentDateTime() int64 {
	return (time.Now().UnixMilli() * 10000) + EPOCH_OFFSET
}
