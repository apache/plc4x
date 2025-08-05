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
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/opcua/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type EncryptionHandler struct {
	serverCertificate *x509.Certificate
	clientCertificate *x509.Certificate
	clientPrivateKey  *rsa.PrivateKey
	clientPublicKey   crypto.PublicKey
	securityPolicy    string

	log zerolog.Logger
}

func NewEncryptionHandler(log zerolog.Logger, ckp *CertificateKeyPair, senderCertificate []byte, securityPolicy string) *EncryptionHandler {
	e := &EncryptionHandler{
		securityPolicy: securityPolicy,
		log:            log,
	}
	if ckp != nil {
		e.clientPrivateKey = ckp.getKeyPair()
		e.clientPublicKey = ckp.getKeyPair().Public()
		e.clientCertificate = ckp.getCertificate()
	}
	if senderCertificate != nil {
		var err error
		e.serverCertificate, err = e.getCertificateX509(senderCertificate)
		if err != nil {
			e.log.Error().Err(err).Msg("error getting sender certificate")
		}
	}
	return e
}

func (h *EncryptionHandler) encodeMessage(ctx context.Context, pdu readWriteModel.MessagePDU, message []byte) ([]byte, error) {
	const PREENCRYPTED_BLOCK_LENGTH = 190
	unencryptedLength := int(pdu.GetLengthInBytes(ctx))
	openRequestLength := len(message)
	positionFirstBlock := unencryptedLength - openRequestLength - 8
	paddingSize := PREENCRYPTED_BLOCK_LENGTH - ((openRequestLength + 256 + 1 + 8) % PREENCRYPTED_BLOCK_LENGTH)
	preEncryptedLength := openRequestLength + 256 + 1 + 8 + paddingSize
	if preEncryptedLength%PREENCRYPTED_BLOCK_LENGTH != 0 {
		return nil, errors.Errorf("Pre encrypted block length %d isn't a multiple of the block size", preEncryptedLength)
	}
	numberOfBlocks := preEncryptedLength / PREENCRYPTED_BLOCK_LENGTH
	encryptedLength := numberOfBlocks*256 + positionFirstBlock
	buf := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := readWriteModel.NewOpcuaAPU(pdu, false, true).SerializeWithWriteBuffer(ctx, buf); err != nil {
		return nil, errors.Wrap(err, "error serializing")
	}
	paddingByte := byte(paddingSize)
	if err := buf.WriteByte("", paddingByte); err != nil {
		return nil, errors.Wrap(err, "error writing byte")
	}
	for i := 0; i < paddingSize; i++ {
		if err := buf.WriteByte("", paddingByte); err != nil {
			return nil, errors.Wrap(err, "error writing byte")
		}
	}
	//Writing Message Length
	{
		if err := buf.WriteInt32("", 32, int32(encryptedLength)); err != nil {
			return nil, errors.Wrap(err, "error writing int32")
		}
		allBytes := buf.GetBytes()
		encryptedLengthBytes := allBytes[len(allBytes)-4:]
		allBytes = allBytes[:len(allBytes)-4]
		allBytes = append(allBytes[:4], append(encryptedLengthBytes, allBytes[8:]...)...)
		buf = utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
		if err := buf.WriteByteArray("", allBytes); err != nil {
			return nil, errors.Wrap(err, "error writing back bytes")
		}
	}

	signature, err := h.sign(buf.GetBytes()[:unencryptedLength+paddingSize+1])
	if err != nil {
		return nil, errors.Wrap(err, "error signing")
	}
	//Write the signature to the end of the buffer
	for _, b := range signature {
		if err := buf.WriteByte("", b); err != nil {
			return nil, errors.Wrap(err, "error writing byte")
		}
	}
	//buf.SetPos(uint16(positionFirstBlock))// TODO: check if we need to move the position at all
	if err := h.encryptBlock(buf, buf.GetBytes()[positionFirstBlock:positionFirstBlock+preEncryptedLength]); err != nil {
		return nil, errors.Wrap(err, "error encrypting block")
	}
	return buf.GetBytes(), nil
}

func (h *EncryptionHandler) decodeMessage(ctx context.Context, pdu readWriteModel.OpcuaAPU) (readWriteModel.OpcuaAPU, error) {
	h.log.Info().Str("securityPolicy", h.securityPolicy).Msg("Decoding Message with Security policy")
	switch h.securityPolicy {
	case "None":
		return pdu, nil
	case "Basic256Sha256":
		var message []byte
		switch pduMessage := pdu.GetMessage().(type) {
		case readWriteModel.OpcuaOpenResponse:
			message = pduMessage.(readWriteModel.BinaryPayload).GetPayload()
		case readWriteModel.OpcuaMessageResponse:
			message = pduMessage.(readWriteModel.BinaryPayload).GetPayload()
		default:
			h.log.Trace().Type("pdu", pdu).Msg("unhandled type")
			return pdu, nil
		}
		encryptedLength := int(pdu.GetLengthInBytes(ctx))
		encryptedMessageLength := len(message) + 8
		headerLength := encryptedLength - encryptedMessageLength
		buf := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
		if err := pdu.SerializeWithWriteBuffer(ctx, buf); err != nil {
			return nil, errors.Wrap(err, "error serializing")
		}
		allBytes := buf.GetBytes()
		data := allBytes[headerLength:encryptedLength]
		buf = utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
		if err := buf.WriteByteArray("", allBytes[:headerLength-1]); err != nil {
			return nil, errors.Wrap(err, "error serializing")
		}
		if err := h.decryptBlock(buf, data); err != nil {
			return nil, errors.Wrap(err, "error decrypting")
		}
		{
			if err := buf.WriteInt32("", 32, int32(encryptedLength)); err != nil {
				return nil, errors.Wrap(err, "error writing int32")
			}
			encryptedLengthBytes := allBytes[len(allBytes)-4:]
			allBytes = allBytes[:len(allBytes)-4]
			allBytes = append(allBytes[:4], append(encryptedLengthBytes, allBytes[8:]...)...)
			buf = utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
			if err := buf.WriteByteArray("", allBytes); err != nil {
				return nil, errors.Wrap(err, "error writing back bytes")
			}
		}

		readBuffer := utils.NewReadBufferByteBased(buf.GetBytes(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
		return readWriteModel.OpcuaAPUParseWithBuffer(ctx, readBuffer, true, true)
	default:
		h.log.Trace().Msg("unmapped security policy")
		return pdu, nil
	}
}

func (h *EncryptionHandler) decryptBlock(buf utils.WriteBufferByteBased, data []byte) error {
	oaep, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, h.clientPrivateKey, data, nil)
	if err != nil {
		return errors.Wrap(err, "error DecryptOAEP")
	}
	return buf.WriteByteArray("", oaep)
}

func (h *EncryptionHandler) getCertificateX509(senderCertificate []byte) (*x509.Certificate, error) {
	return x509.ParseCertificate(senderCertificate)
}

func (h *EncryptionHandler) setServerCertificate(serverCertificate *x509.Certificate) {
	h.serverCertificate = serverCertificate
}

func (h *EncryptionHandler) encryptPassword(password []byte) ([]byte, error) {
	publicKey := h.serverCertificate.PublicKey.(rsa.PublicKey)
	encryptOAEP, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &publicKey, password, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error EncryptOAEP")
	}
	return encryptOAEP, nil
}

func (h *EncryptionHandler) encryptBlock(buf utils.WriteBufferByteBased, data []byte) error {
	publicKey := h.serverCertificate.PublicKey.(rsa.PublicKey)
	encryptOAEP, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &publicKey, data, nil)
	if err != nil {
		return errors.Wrap(err, "error EncryptOAEP")
	}
	return buf.WriteByteArray("", encryptOAEP)
}

func (h *EncryptionHandler) sign(data []byte) ([]byte, error) {
	return h.clientPrivateKey.Sign(rand.Reader, data, crypto.SHA256)
}

func (h *EncryptionHandler) String() string {
	return "EncryptionHandler{" + h.securityPolicy + "}"
}
