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
package org.apache.plc4x.java.opcua.security;

import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

import javax.crypto.Cipher;
import java.security.Security;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * The password of a UserNameIdentityToken has to be encrypted with the algorithm of the selected
 * user token policy, and the token has to declare exactly that algorithm - see GH-2154. The driver
 * used to hardcode RSA-OAEP and the matching URI, so servers advertising Basic128Rsa15 (which
 * requires RSA-PKCS#1 v1.5) rejected the token with BadIdentityTokenInvalid.
 */
class SecurityPolicyPasswordEncryptionTest {

    @BeforeAll
    static void registerBouncyCastle() {
        // The driver registers BouncyCastle when it builds an EncryptionHandler; some of the
        // transformation names below (notably BC's OAEPWithSHA256AndMGF1Padding, which is *not*
        // the same as SunJCE's OAEPWithSHA-256AndMGF1Padding) only resolve with it present.
        Security.addProvider(new BouncyCastleProvider());
    }

    @Test
    void basic128Rsa15UsesPkcs1() throws Exception {
        SecurityPolicy policy = SecurityPolicy.findByUri(
            "http://opcfoundation.org/UA/SecurityPolicy#Basic128Rsa15").orElse(null);

        assertNotNull(policy);
        assertEquals("http://www.w3.org/2001/04/xmlenc#rsa-1_5",
            policy.getAsymmetricEncryptionAlgorithm().getUri());
        assertEquals("RSA/ECB/PKCS1Padding", policy.getAsymmetricEncryptionAlgorithm().getName());
    }

    @Test
    void sha1BasedPoliciesUseRsaOaep() throws Exception {
        for (String uri : new String[]{
            "http://opcfoundation.org/UA/SecurityPolicy#Basic256",
            "http://opcfoundation.org/UA/SecurityPolicy#Basic256Sha256",
            "http://opcfoundation.org/UA/SecurityPolicy#Aes128_Sha256_RsaOaep"}) {
            SecurityPolicy policy = SecurityPolicy.findByUri(uri).orElse(null);
            assertNotNull(policy, uri);
            assertEquals("http://www.w3.org/2001/04/xmlenc#rsa-oaep",
                policy.getAsymmetricEncryptionAlgorithm().getUri(), uri);
        }
    }

    @Test
    void rsaPssPolicyUsesOaepSha256() {
        SecurityPolicy policy = SecurityPolicy.findByUri(
            "http://opcfoundation.org/UA/SecurityPolicy#Aes256_Sha256_RsaPss").orElse(null);

        assertNotNull(policy);
        assertEquals("http://opcfoundation.org/UA/security#rsa-oaep-sha2-256",
            policy.getAsymmetricEncryptionAlgorithm().getUri());
    }

    /**
     * A policy of None means the password is not encrypted at all, so it declares no algorithm.
     */
    @Test
    void noneDeclaresNoAlgorithm() {
        SecurityPolicy policy = SecurityPolicy.findByUri(
            "http://opcfoundation.org/UA/SecurityPolicy#None").orElse(null);

        assertNotNull(policy);
        assertTrue(policy.getAsymmetricEncryptionAlgorithm().getUri().isEmpty());
    }

    /**
     * Every declared algorithm has to be one the JVM can actually instantiate - a typo in the
     * transformation name would only show up at connect time otherwise.
     */
    @Test
    void everyEncryptingPolicyResolvesToAUsableCipher() throws Exception {
        for (SecurityPolicy policy : SecurityPolicy.values()) {
            if (policy == SecurityPolicy.NONE) {
                continue;
            }
            Cipher cipher = policy.getAsymmetricEncryptionAlgorithm().getCipher();
            assertNotNull(cipher, policy.name());
            assertTrue(policy.getAsymmetricEncryptionAlgorithm().getUri().startsWith("http"),
                policy + " must declare an algorithm URI");
        }
    }
}
