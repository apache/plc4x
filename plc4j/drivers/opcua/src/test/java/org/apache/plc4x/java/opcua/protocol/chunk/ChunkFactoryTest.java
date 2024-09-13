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

package org.apache.plc4x.java.opcua.protocol.chunk;

import static org.junit.jupiter.api.Assertions.*;

import java.security.PrivateKey;
import java.security.cert.X509Certificate;
import java.util.HashMap;
import java.util.Map;
import java.util.Map.Entry;
import org.apache.plc4x.java.opcua.TestCertificateGenerator;
import org.apache.plc4x.java.opcua.readwrite.MessageSecurityMode;
import org.apache.plc4x.java.opcua.readwrite.OpcuaProtocolLimits;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.test.DisableOnJenkinsFlag;
import org.junit.jupiter.api.Assumptions;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvFileSource;

@DisableOnJenkinsFlag
class ChunkFactoryTest {

    public static final Map<Integer, Entry<PrivateKey, X509Certificate>> CERTIFICATES = new HashMap<>();

    private OpcuaProtocolLimits limits = new OpcuaProtocolLimits(
        8196,
        8196,
        8196 * 10,
        10
    );

    @ParameterizedTest
    @CsvFileSource(numLinesToSkip = 1, resources = {
        "/chunk-calculation-1024.csv",
        "/chunk-calculation-2048.csv",
        "/chunk-calculation-3072.csv",
        "/chunk-calculation-1024.csv",
        "/chunk-calculation-5120.csv"
    })
    public void testChunkCalculation(
        int keySize,
        String securityPolicy,
        String messageSecurity,
        boolean asymmetric,
        boolean encrypted,
        boolean signed,
        int securityHeaderSize,
        int cipherTextBlockSize,
        int plainTextBlockSize,
        int signatureSize,
        int maxChunkSize,
        int paddingOverhead,
        int maxCipherTextSize,
        int maxCipherTextBlocks,
        int maxPlainTextSize,
        int maxBodySize
    ) throws Exception {
        verify(get(keySize),
            securityPolicy,
            messageSecurity,
            asymmetric,
            encrypted,
            signed,
            securityHeaderSize,
            cipherTextBlockSize,
            plainTextBlockSize,
            signatureSize,
            maxChunkSize,
            paddingOverhead,
            maxCipherTextSize,
            maxCipherTextBlocks,
            maxPlainTextSize,
            maxBodySize
        );
    }

    private void verify(Entry<PrivateKey, X509Certificate> certificateEntry, String securityPolicy, String messageSecurity,
        boolean asymmetric, boolean encrypted, boolean signed,
        int securityHeaderSize, int cipherTextBlockSize, int plainTextBlockSize, int signatureSize,
        int maxChunkSize, int paddingOverhead, int maxCipherTextSize, int maxCipherTextBlocks, int maxPlainTextSize, int maxBodySize) {
        SecurityPolicy channelSecurityPolicy = null;
        try {
            channelSecurityPolicy = SecurityPolicy.valueOf(securityPolicy);
        } catch (IllegalArgumentException e) {
            Assumptions.abort("Unsupported security policy " + securityPolicy);
        }
        MessageSecurityMode channelMessageSecurity = null;
        try {
            channelMessageSecurity = MessageSecurityMode.valueOf(messageSecurity);
        } catch (IllegalArgumentException e) {
            Assumptions.abort("Unsupported security policy " + securityPolicy);
        }

        ChunkFactory chunkFactory = new ChunkFactory();
        Chunk chunk = chunkFactory.create(
            asymmetric, encrypted, signed,
            channelSecurityPolicy,
            limits,
            certificateEntry.getValue(),
            certificateEntry.getValue()
        );

        assertEquals(securityHeaderSize, chunk.getSecurityHeaderSize(), "securityHeaderSize mismatch");
        assertEquals(cipherTextBlockSize, chunk.getCipherTextBlockSize(), "cipherTextBlockSize mismatch");
        assertEquals(asymmetric, chunk.isAsymmetric(), "asymmetric mismatch");
        assertEquals(encrypted, chunk.isEncrypted(), "encrypted mismatch");
        assertEquals(signed, chunk.isSigned(), "signed mismatch");
        assertEquals(plainTextBlockSize, chunk.getPlainTextBlockSize(), "plainTextBlockSize mismatch");
        assertEquals(signatureSize, chunk.getSignatureSize(), "signatureSize mismatch");
        assertEquals(maxChunkSize, chunk.getMaxChunkSize(), "maxChunkSize mismatch");
        assertEquals(paddingOverhead, chunk.getPaddingOverhead(), "paddingOverhead mismatch");
        assertEquals(maxCipherTextSize, chunk.getMaxCipherTextSize(), "maxCipherTextSize mismatch");
        assertEquals(maxCipherTextBlocks, chunk.getMaxCipherTextBlocks(), "maxCipherTextBlocks mismatch");
        assertEquals(maxPlainTextSize, chunk.getMaxPlainTextSize(), "maxPlainTextSize mismatch");
        assertEquals(maxBodySize, chunk.getMaxBodySize(), "maxBodySize mismatch");
    }

    private static Entry<PrivateKey, X509Certificate> get(int keySize) {
        return CERTIFICATES.computeIfAbsent(keySize, (ks) -> TestCertificateGenerator.generate(ks, "cn=test", 10));
    }

}