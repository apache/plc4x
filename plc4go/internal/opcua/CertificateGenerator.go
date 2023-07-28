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
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"

	"github.com/pkg/errors"
)

func generateCertificate() (*CertificateKeyPair, error) {
	// From: https://gist.github.com/shaneutt/5e1995295cff6721c89a71d13a71c251
	// set up our CA certificate
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(40),
		Subject: pkix.Name{
			CommonName:         "Apache PLC4X Driver Client",
			Organization:       []string{"Apache Software Foundation"},
			OrganizationalUnit: []string{"dev"},
			Locality:           []string{""},
			Country:            []string{"US"},
			Province:           []string{""},
		},
		NotBefore:             time.Now().Add(-24 * time.Hour),
		NotAfter:              time.Now().AddDate(25, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// create our private and public key
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, errors.Wrap(err, "error generating key")
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "error creating certificate")
	}

	// pem encode
	caPEM := new(bytes.Buffer)
	if err := pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}); err != nil {
		return nil, errors.Wrap(err, "error pem encode")
	}

	caPrivateKeyPEM := new(bytes.Buffer)
	if err := pem.Encode(caPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivateKey),
	}); err != nil {
		return nil, errors.Wrap(err, "error pem encode")
	}

	// set up our server certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, errors.Wrap(err, "error generating key")
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivateKey.PublicKey, caPrivateKeyPEM)
	if err != nil {
		return nil, errors.Wrap(err, "error creating certificate")
	}

	certPEM := new(bytes.Buffer)
	if err := pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return nil, errors.Wrap(err, "error pem encode")
	}

	certPrivKeyPEM := new(bytes.Buffer)
	if err := pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivateKey),
	}); err != nil {
		return nil, errors.Wrap(err, "error pem encode")
	}

	serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "error building keypair")
	}

	serverTLSConf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}
	_ = serverTLSConf

	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(caPEM.Bytes())
	clientTLSConf := &tls.Config{
		RootCAs: certpool,
	}
	_ = clientTLSConf

	return NewCertificateKeyPair(caPrivateKey, cert), nil
}
