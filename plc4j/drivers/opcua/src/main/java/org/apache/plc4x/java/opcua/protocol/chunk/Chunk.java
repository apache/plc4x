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

import java.util.Objects;

public class Chunk {

    private final int securityHeaderSize;
    private final int cipherTextBlockSize;
    private final int plainTextBlockSize;
    private final int signatureSize;
    private final int maxChunkSize;
    private final int paddingOverhead;
    private final int maxCipherTextSize;
    private final int maxCipherTextBlocks;
    private final int maxPlainTextSize;
    private final int maxBodySize;

    private boolean asymmetric;
    private boolean encrypted;
    private boolean signed;

    public Chunk(int securityHeaderSize, int cipherTextBlockSize, int plainTextBlockSize, int signatureSize, int maxChunkSize) {
        this(securityHeaderSize, cipherTextBlockSize, plainTextBlockSize, signatureSize, maxChunkSize, false, false, false);
    }

    public Chunk(int securityHeaderSize, int cipherTextBlockSize, int plainTextBlockSize, int signatureSize, int maxChunkSize,
        boolean asymmetric, boolean encrypted, boolean signed) {
        this.securityHeaderSize = securityHeaderSize;
        this.cipherTextBlockSize = cipherTextBlockSize;
        this.plainTextBlockSize = plainTextBlockSize;
        this.signatureSize = signatureSize;
        this.maxChunkSize = maxChunkSize;
        this.asymmetric = asymmetric;
        this.encrypted = encrypted;
        this.signed = signed;
        this.maxCipherTextSize = maxChunkSize - 12 /* security header */ - securityHeaderSize;
        this.maxCipherTextBlocks = maxCipherTextSize / cipherTextBlockSize;
        this.paddingOverhead = cipherTextBlockSize > 256 ? 2 : (cipherTextBlockSize < 16 ? 0 : 1);
        this.maxPlainTextSize = maxCipherTextBlocks * plainTextBlockSize;
        this.maxBodySize = maxPlainTextSize - 8 /* sequence header */ - this.paddingOverhead - signatureSize;
    }

    public int getSecurityHeaderSize() {
        return securityHeaderSize;
    }
    public int getCipherTextBlockSize() {
        return cipherTextBlockSize;
    }
    public int getPlainTextBlockSize() {
        return plainTextBlockSize;
    }
    public int getSignatureSize() {
        return signatureSize;
    }
    public int getMaxChunkSize() {
        return maxChunkSize;
    }
    public int getPaddingOverhead() {
        return paddingOverhead;
    }
    public int getMaxCipherTextSize() {
        return maxCipherTextSize;
    }
    public int getMaxCipherTextBlocks() {
        return maxCipherTextBlocks;
    }
    public int getMaxPlainTextSize() {
        return maxPlainTextSize;
    }
    public int getMaxBodySize() {
        return maxBodySize;
    }

    public boolean isAsymmetric() {
        return asymmetric;
    }

    public boolean isEncrypted() {
        return encrypted;
    }

    public boolean isSigned() {
        return signed;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof Chunk)) {
            return false;
        }
        Chunk chunk = (Chunk) o;
        return getSecurityHeaderSize() == chunk.getSecurityHeaderSize()
            && getCipherTextBlockSize() == chunk.getCipherTextBlockSize()
            && getPlainTextBlockSize() == chunk.getPlainTextBlockSize()
            && getSignatureSize() == chunk.getSignatureSize()
            && getMaxChunkSize() == chunk.getMaxChunkSize()
            && getPaddingOverhead() == chunk.getPaddingOverhead()
            && getMaxCipherTextSize() == chunk.getMaxCipherTextSize()
            && getMaxCipherTextBlocks() == chunk.getMaxCipherTextBlocks()
            && getMaxPlainTextSize() == chunk.getMaxPlainTextSize()
            && getMaxBodySize() == chunk.getMaxBodySize();
    }

    @Override
    public int hashCode() {
        return Objects.hash(getSecurityHeaderSize(), getCipherTextBlockSize(),
            getPlainTextBlockSize(),
            getSignatureSize(), getMaxChunkSize(), getPaddingOverhead(), getMaxCipherTextSize(),
            getMaxCipherTextBlocks(), getMaxPlainTextSize(), getMaxBodySize());
    }

    @Override
    public String toString() {
        return "Chunk" +
            "{ securityHeaderSize=" + securityHeaderSize +
            ", cipherTextBlockSize=" + cipherTextBlockSize +
            ", plainTextBlockSize=" + plainTextBlockSize +
            ", signatureSize=" + signatureSize +
            ", maxChunkSize=" + maxChunkSize +
            ", paddingOverhead=" + paddingOverhead +
            ", maxCipherTextSize=" + maxCipherTextSize +
            ", maxCipherTextBlocks=" + maxCipherTextBlocks +
            ", maxPlainTextSize=" + maxPlainTextSize +
            ", maxBodySize=" + maxBodySize +
            '}';
    }
}
