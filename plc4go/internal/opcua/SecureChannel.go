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
	"context"
	"github.com/rs/zerolog"
	"math/rand"
	"net"
	"net/url"
	"regexp"
	"sync/atomic"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/opcua/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/dchest/uniuri"
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
	DEFAULT_CONNECTION_LIFETIME = 36000000
)

type SecureChannel struct {
	sessionName               string
	clientNonce               []byte
	requestHandleGenerator    atomic.Int32
	policyId                  readWriteModel.PascalString
	tokenType                 readWriteModel.UserTokenType
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
	configuration             Configuration
	channelId                 atomic.Int32
	tokenId                   atomic.Int32
	authenticationToken       readWriteModel.NodeIdTypeDefinition
	context                   *MessageCodec // TODO: not sure if we need the codec here
	channelTransactionManager *SecureChannelTransactionManager
	lifetime                  int
	keepAlive                 chan struct{} // TODO: check if this is the right thing
	sendBufferSize            int
	maxMessageSize            int
	endpoints                 []string
	senderSequenceNumber      atomic.Int64
	log                       zerolog.Logger
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
	ckp := configuration.ckp
	if configuration.securityPolicy == "Basic256Sha256" {
		//Sender Certificate gets populated during the 'discover' phase when encryption is enabled.
		s.senderCertificate = configuration.senderCertificate
		s.encryptionHandler = NewEncryptionHandler(ckp, s.senderCertificate, configuration.securityPolicy)
		certificate := ckp.getCertificate()
		s.publicCertificate = readWriteModel.NewPascalByteString(int32(len(certificate.Raw)), certificate.Raw)
		s.isEncrypted = true

		s.thumbprint = configuration.thumbprint
	} else {
		s.encryptionHandler = NewEncryptionHandler(ckp, s.senderCertificate, configuration.securityPolicy)
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

func (s *SecureChannel) Something() {

}

func (s *SecureChannel) getAuthenticationToken() readWriteModel.NodeId {
	return readWriteModel.NewNodeId(s.authenticationToken)
}

func (s *SecureChannel) getCurrentDateTime() int64 {
	return (time.Now().UnixMilli() * 10000) + EPOCH_OFFSET
}

func (s *SecureChannel) submit(ctx context.Context, codec *MessageCodec, errorDispatcher func(err error), result chan apiModel.PlcReadRequestResult, consumer func(opcuaResponse []byte), buffer utils.WriteBufferByteBased) {
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
		var err error
		apu, err = readWriteModel.OpcuaAPUParse(ctx, s.encryptionHandler.encodeMessage(messageRequest, buffer.GetBytes()), false)
		if err != nil {
			errorDispatcher(err)
			return
		}
	} else {
		apu = readWriteModel.NewOpcuaAPU(messageRequest, false)
	}

	requestConsumer := func(transactionId int32) {
		// bos
		//codec.SendRequest(ctx,) // TODO: feed
		_ = apu
	}

	s.log.Debug().Msgf("Submitting Transaction to TransactionManager %v", transactionId)
	s.channelTransactionManager.submit(requestConsumer, transactionId)
}
