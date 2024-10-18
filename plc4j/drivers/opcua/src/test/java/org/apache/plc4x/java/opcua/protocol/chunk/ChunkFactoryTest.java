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

import java.security.PrivateKey;
import java.security.cert.X509Certificate;
import java.util.Map;
import java.util.Map.Entry;
import java.util.concurrent.ConcurrentHashMap;
import org.apache.plc4x.java.opcua.TestCertificateGenerator;
import org.apache.plc4x.java.opcua.readwrite.OpcuaProtocolLimits;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.assertj.core.api.SoftAssertions;
import org.junit.jupiter.api.Assumptions;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvFileSource;

class ChunkFactoryTest {

    public static final Map<Integer, Entry<PrivateKey, X509Certificate>> CERTIFICATES = new ConcurrentHashMap<>();

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

        ChunkFactory chunkFactory = new ChunkFactory();
        Chunk chunk = chunkFactory.create(
            asymmetric, encrypted, signed,
            channelSecurityPolicy,
            limits,
            certificateEntry.getValue(),
            certificateEntry.getValue()
        );

        SoftAssertions.assertSoftly(softly -> {
            softly.assertThat(asymmetric).isEqualTo(asymmetric);
            softly.assertThat(encrypted).isEqualTo(encrypted);
            softly.assertThat(signed).isEqualTo(chunk.isSigned());
            softly.assertThat(cipherTextBlockSize).isEqualTo(chunk.getCipherTextBlockSize());
            softly.assertThat(plainTextBlockSize).isEqualTo(chunk.getPlainTextBlockSize());
            softly.assertThat(signatureSize).isEqualTo(chunk.getSignatureSize());
            softly.assertThat(maxChunkSize).isEqualTo(chunk.getMaxChunkSize());
            softly.assertThat(paddingOverhead).isEqualTo(chunk.getPaddingOverhead());
            softly.assertThat(maxCipherTextSize).isEqualTo(chunk.getMaxCipherTextSize());
            softly.assertThat(maxCipherTextBlocks).isEqualTo(chunk.getMaxCipherTextBlocks());
            softly.assertThat(maxPlainTextSize).isEqualTo(chunk.getMaxPlainTextSize());
            softly.assertThat(maxBodySize).isEqualTo(chunk.getMaxBodySize());
            softly.assertThat(securityHeaderSize).isEqualTo(chunk.getSecurityHeaderSize());
        });
    }

    private static Entry<PrivateKey, X509Certificate> get(int keySize) {
        return CERTIFICATES.computeIfAbsent(keySize, (ks) -> TestCertificateGenerator.generate(ks, "cn=test", 10));
    }

}