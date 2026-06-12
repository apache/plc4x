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
package org.apache.plc4x.java.spi.buffers.bytebased;

import org.apache.plc4x.java.spi.buffers.api.AbstractBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.byteorder.ByteOrder;
import org.apache.plc4x.java.spi.buffers.bytebased.byteorder.ByteOrderBigEndian;
import org.apache.plc4x.java.spi.buffers.bytebased.byteorder.ByteOrderManager;
import org.apache.plc4x.java.spi.buffers.bytebased.encoding.Encoding;
import org.apache.plc4x.java.spi.buffers.bytebased.encoding.EncodingManager;

import java.util.Objects;
import java.util.Optional;

public abstract class AbstractBufferByteBased extends AbstractBuffer {

    protected final byte[] buffer;
    protected final int startBit;
    protected final int sizeInBits;
    protected int positionInBits;

    private final ByteOrderManager byteOrderManager;
    private final EncodingManager encodingManager;

    public AbstractBufferByteBased(byte[] buffer, int startBit, int sizeInBits, ByteOrderManager byteOrderManager, EncodingManager encodingManager, WithOption[] options) {
        super(options);
        this.buffer = Objects.requireNonNull(buffer);
        this.startBit = startBit;
        this.sizeInBits = sizeInBits;
        this.positionInBits = 0;
        this.byteOrderManager = byteOrderManager;
        this.encodingManager = encodingManager;

        if (startBit < 0 || startBit > sizeInBits) {
            throw new IllegalArgumentException("Start bit must be between 0 and sizeInBits");
        }
        if (sizeInBits > buffer.length * 8) {
            throw new IllegalArgumentException("Size in bits must not exceed buffer size");
        }
    }

    protected ByteOrder getByteOrder(WithOption... options) {
        Optional<String> byteOrderName = WithByteBasedOption.extractByteOrder(options, getContext());
        if (byteOrderName.isPresent()) {
            Optional<ByteOrder> byteOrder = byteOrderManager.getByteOrder(byteOrderName.get());
            if (byteOrder.isPresent()) {
                return byteOrder.get();
            }
        }
        return new ByteOrderBigEndian();
    }

    protected Optional<Encoding> getUnsignedIntegerEncoding(WithOption... options) {
        // First, try to get the encoding for this explicit type of field.
        Optional<String> encodingName = WithOption.extractUnsignedIntegerEncoding(options, getContext());
        // If that fails, try to get the global encoder.
        if (encodingName.isEmpty()) {
            encodingName = WithOption.extractEncoding(options, getContext());
            if(encodingName.isPresent()) {
                return encodingManager.getEncoding(encodingName.get());
            }
        } else {
            return encodingManager.getEncoding(encodingName.get());
        }
        return Optional.empty();
    }

    protected Optional<Encoding> getSignedIntegerEncoding(WithOption... options) {
        // First, try to get the encoding for this explicit type of field.
        Optional<String> encodingName = WithOption.extractSignedIntegerEncoding(options, getContext());
        // If that fails, try to get the global encoder.
        if (encodingName.isEmpty()) {
            encodingName = WithOption.extractEncoding(options, getContext());
            if(encodingName.isPresent()) {
                return encodingManager.getEncoding(encodingName.get());
            }
        }
        if (encodingName.isPresent()) {
            return encodingManager.getEncoding(encodingName.get());
        }
        return Optional.empty();
    }

    protected Optional<Encoding> getFloatEncoding(WithOption... options) {
        // First, try to get the encoding for this explicit type of field.
        Optional<String> encodingName = WithOption.extractFloatEncoding(options, getContext());
        // If that fails, try to get the global encoder.
        if (encodingName.isEmpty()) {
            encodingName = WithOption.extractEncoding(options, getContext());
            if(encodingName.isPresent()) {
                return encodingManager.getEncoding(encodingName.get());
            }
        }
        if (encodingName.isPresent()) {
            return encodingManager.getEncoding(encodingName.get());
        }
        return Optional.empty();
    }

    protected Optional<Encoding> getStringEncoding(WithOption... options) {
        // First, try to get the encoding for this explicit type of field.
        Optional<String> encodingName = WithOption.extractStringEncoding(options, getContext());
        // If that fails, try to get the global encoder.
        if (encodingName.isEmpty()) {
            encodingName = WithOption.extractEncoding(options, getContext());
            if(encodingName.isPresent()) {
                return encodingManager.getEncoding(encodingName.get());
            }
        }
        if (encodingName.isPresent()) {
            return encodingManager.getEncoding(encodingName.get());
        }
        return Optional.empty();
    }

    public byte[] getBytes() {
        return buffer;
    }

    public int getPositionInBits() {
        return positionInBits;
    }

    public int getRemainingBits() {
        return sizeInBits - positionInBits;
    }

    public void setPositionInBits(int positionInBits) {
        if (positionInBits < 0 || positionInBits > sizeInBits) {
            throw new IllegalArgumentException("Position must be between 0 and sizeInBits");
        }
        this.positionInBits = positionInBits;
    }

    protected void ensureAvailable(int bitsNeeded) throws BufferException {
        if (positionInBits + bitsNeeded > sizeInBits) {
            throw new BufferException("Not enough bits left to read");
        }
    }

    protected boolean isAligned() {
        return (positionInBits % 8) == 0;
    }

}
