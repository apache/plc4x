/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.opcua;

import java.math.BigInteger;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.NoSuchAlgorithmException;
import java.security.PrivateKey;
import java.security.SecureRandom;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.util.Date;
import java.util.Map;
import java.util.Map.Entry;
import java.util.Random;
import java.util.concurrent.TimeUnit;
import javax.security.auth.x500.X500Principal;
import org.bouncycastle.cert.X509CertificateHolder;
import org.bouncycastle.cert.X509v3CertificateBuilder;
import org.bouncycastle.cert.jcajce.JcaX509CertificateConverter;
import org.bouncycastle.cert.jcajce.JcaX509v3CertificateBuilder;
import org.bouncycastle.operator.OperatorCreationException;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;

public class TestCertificateGenerator {

    public static Entry<PrivateKey, X509Certificate> generate(int keySize, String dn, long validitySec) {
        try {
            KeyPairGenerator keyPairGenerator = KeyPairGenerator.getInstance("RSA");
            keyPairGenerator.initialize(keySize, new SecureRandom());
            KeyPair keyPair = keyPairGenerator.generateKeyPair();

            X509v3CertificateBuilder certGen = new JcaX509v3CertificateBuilder(
                new X500Principal(dn),
                BigInteger.valueOf(new Random().nextLong()),
                new Date(),
                new Date(System.currentTimeMillis() + TimeUnit.SECONDS.toMillis(validitySec)),
                new X500Principal(dn),
                keyPair.getPublic()
            );
            X509CertificateHolder cert = certGen.build(new JcaContentSignerBuilder("SHA256withRSA")
                .build(keyPair.getPrivate()));

            X509Certificate certificate = new JcaX509CertificateConverter().getCertificate(cert);
            return Map.entry(keyPair.getPrivate(), certificate);
        } catch (CertificateException | NoSuchAlgorithmException | OperatorCreationException e) {
            throw new RuntimeException("Could not initialize test - certificate generation failed");
        }
    }

}
