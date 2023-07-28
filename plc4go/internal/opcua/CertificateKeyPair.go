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
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
)

type CertificateKeyPair struct {
	keyPair     *rsa.PrivateKey
	certificate *x509.Certificate
	thumbprint  []byte
}

func NewCertificateKeyPair(keyPair *rsa.PrivateKey, certificate *x509.Certificate) *CertificateKeyPair {
	thumbprint := sha1.Sum(certificate.Raw)
	return &CertificateKeyPair{
		keyPair:     keyPair,
		certificate: certificate,
		thumbprint:  thumbprint[:],
	}
}

func (p *CertificateKeyPair) getKeyPair() *rsa.PrivateKey {
	return p.keyPair
}

func (p *CertificateKeyPair) getCertificate() *x509.Certificate {
	return p.certificate
}

func (p *CertificateKeyPair) getThumbprint() []byte {
	return p.thumbprint
}
