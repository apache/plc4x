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
package org.apache.plc4x.java.opcua.context;

import org.bouncycastle.asn1.DERSequence;
import org.bouncycastle.asn1.x500.X500NameBuilder;
import org.bouncycastle.asn1.x500.style.BCStyle;
import org.bouncycastle.asn1.x509.*;
import org.bouncycastle.cert.X509CertificateHolder;
import org.bouncycastle.cert.X509v3CertificateBuilder;
import org.bouncycastle.cert.jcajce.JcaX509CertificateConverter;
import org.bouncycastle.cert.jcajce.JcaX509ExtensionUtils;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.operator.ContentSigner;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.net.InetAddress;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.util.Calendar;
import java.util.Date;
import java.util.Locale;

public class CertificateGenerator<PKCS10CertificateRequest> {

    private static final Logger LOGGER = LoggerFactory.getLogger(CertificateGenerator.class);
    private static final String APPURI = "urn:eclipse:milo:plc4x:server";

    /** Key size used when the connection string doesn't ask for a specific one. */
    public static final int DEFAULT_KEY_SIZE = 2048;

    /**
     * Generates a self-signed application instance certificate with a 2048 bit key.
     */
    public static CertificateKeyPair generateCertificate() {
        return generateCertificate(DEFAULT_KEY_SIZE);
    }

    /**
     * Generates a self-signed application instance certificate.
     *
     * @param keySize size of the RSA key in bits; servers may demand a minimum, commonly 4096
     */
    public static CertificateKeyPair generateCertificate(int keySize) {
        X500NameBuilder nameBuilder = new X500NameBuilder();

        nameBuilder.addRDN(BCStyle.CN, "Apache PLC4X Driver Client");
        nameBuilder.addRDN(BCStyle.O, "Apache Software Foundation");
        nameBuilder.addRDN(BCStyle.OU, "dev");
        nameBuilder.addRDN(BCStyle.L, "");
        nameBuilder.addRDN(BCStyle.ST, "DE");
        nameBuilder.addRDN(BCStyle.C, "US");

        // X.509 requires a positive serial number of at most 20 octets. Building it from 40 random
        // bytes produced an over-long and, half the time, negative serial, which strict servers
        // reject.
        BigInteger serial = new BigInteger(127, new SecureRandom()).setBit(126);

        final Calendar calender = Calendar.getInstance();
        calender.add(Calendar.DATE, -1);
        Date startDate = calender.getTime();
        calender.add(Calendar.DATE, 365 * 25);
        Date expiryDate = calender.getTime();

        KeyPairGenerator generator = null;
        try {
            generator = KeyPairGenerator.getInstance("RSA");
            generator.initialize(keySize, new SecureRandom());
            KeyPair keyPair = generator.generateKeyPair();

            SubjectPublicKeyInfo subjectPublicKeyInfo = SubjectPublicKeyInfo.getInstance(
                keyPair.getPublic().getEncoded()
            );

            X509v3CertificateBuilder certificateBuilder = new X509v3CertificateBuilder(
                nameBuilder.build(),
                serial,
                startDate,
                expiryDate,
                Locale.ENGLISH,
                nameBuilder.build(),
                subjectPublicKeyInfo
            );

            GeneralName[] gnArray = new GeneralName[]{new GeneralName(GeneralName.dNSName, InetAddress.getLocalHost().getHostName()), new GeneralName(GeneralName.uniformResourceIdentifier, APPURI)};

            JcaX509ExtensionUtils extensionUtils = new JcaX509ExtensionUtils();
            certificateBuilder.addExtension(Extension.authorityKeyIdentifier, false, extensionUtils.createAuthorityKeyIdentifier(keyPair.getPublic()));
            // Servers look the certificate up by its subject key identifier; without it some
            // reject the certificate outright.
            certificateBuilder.addExtension(Extension.subjectKeyIdentifier, false, extensionUtils.createSubjectKeyIdentifier(keyPair.getPublic()));
            certificateBuilder.addExtension(Extension.extendedKeyUsage, false, new ExtendedKeyUsage(new KeyPurposeId[]{KeyPurposeId.id_kp_clientAuth, KeyPurposeId.id_kp_serverAuth}));
            // The key usages an application instance certificate is supposed to carry
            // (OPC UA Part 6). keyCertSign has no business here - this certificate signs no other
            // certificates.
            certificateBuilder.addExtension(Extension.keyUsage, false, new KeyUsage(KeyUsage.dataEncipherment | KeyUsage.digitalSignature | KeyUsage.keyEncipherment | KeyUsage.nonRepudiation));
            // This is an end entity certificate, not a CA. Claiming to be a CA made servers with
            // stricter certificate checks refuse it.
            certificateBuilder.addExtension(Extension.basicConstraints, false, new BasicConstraints(false));

            GeneralNames subjectAltNames = GeneralNames.getInstance(new DERSequence(gnArray));
            certificateBuilder.addExtension(Extension.subjectAlternativeName, false, subjectAltNames);

            // SHA-1 is deprecated and rejected by the SHA-256 based security policies as well as by
            // any server holding a modern security baseline.
            ContentSigner sigGen = new JcaContentSignerBuilder("SHA256withRSA")
                // Pass the provider instance rather than the name "BC": this class would
                // otherwise depend on some other class having registered it first.
                .setProvider(new BouncyCastleProvider()).build(keyPair.getPrivate());

            X509CertificateHolder certificateHolder = certificateBuilder.build(sigGen);

            JcaX509CertificateConverter certificateConvertor = new JcaX509CertificateConverter();
            certificateConvertor.setProvider(new BouncyCastleProvider());

            CertificateKeyPair ckp = new CertificateKeyPair(keyPair, certificateConvertor.getCertificate(certificateHolder));

            return ckp;

        } catch (Exception e) {
            LOGGER.error("Security Algorithm is unsupported for certificate");
            return null;
        }
    }
}
