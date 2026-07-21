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
package org.apache.plc4x.java.slmp;

import org.apache.plc4x.java.api.exceptions.PlcIncompatibleDatatypeException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.values.PlcDINT;
import org.apache.plc4x.java.spi.values.PlcINT;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcREAL;
import org.apache.plc4x.java.spi.values.PlcUDINT;
import org.apache.plc4x.java.spi.values.PlcUINT;
import org.apache.plc4x.java.spi.values.PlcWORD;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;

/**
 * Maps a supported PLC4X data type to its SLMP word footprint and the conversion
 * between little-endian SLMP words and a {@link PlcValue} in both directions. The SLMP
 * wire layer leaves {@code responseData}/request data as raw bytes; this enum owns typed
 * decoding and encoding of word-unit values (word units only).
 */
public enum SlmpDataType {
    WORD(1),
    INT(1),
    UINT(1),
    DINT(2),
    UDINT(2),
    REAL(2);

    private static final Logger LOGGER = LoggerFactory.getLogger(SlmpDataType.class);

    private final int wordsPerElement;

    SlmpDataType(int wordsPerElement) {
        this.wordsPerElement = wordsPerElement;
    }

    public int getWordsPerElement() {
        return wordsPerElement;
    }

    /**
     * Decode {@code quantity} elements of this type from little-endian SLMP response bytes.
     * Returns {@code null} when {@code responseData} is shorter than required (caller maps to INVALID_DATA).
     */
    public PlcValue decode(byte[] responseData, int quantity) {
        int requiredBytes = quantity * wordsPerElement * 2;
        if (responseData == null || responseData.length < requiredBytes) {
            return null;
        }
        ReadBufferByteBased buffer = new ReadBufferByteBased(responseData,
            WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"),
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithOption.WithFloatEncoding("IEEE754"));
        try {
            if (quantity == 1) {
                return readOne(buffer);
            }
            List<PlcValue> values = new ArrayList<>(quantity);
            for (int i = 0; i < quantity; i++) {
                values.add(readOne(buffer));
            }
            return new PlcList(values);
        } catch (BufferException e) {
            LOGGER.warn("Failed to decode SLMP {} value", this, e);
            return null;
        }
    }

    /**
     * Encode {@code quantity} elements of this type to little-endian SLMP request bytes
     * ({@code 2 * wordsPerElement} per element). Returns {@code null} on a type/arity mismatch
     * (caller maps to INVALID_DATA), symmetric with {@link #decode}.
     */
    public byte[] encode(PlcValue value, int quantity) {
        int totalBytes = quantity * wordsPerElement * 2;
        WriteBufferByteBased buffer = new WriteBufferByteBased(new byte[totalBytes],
            WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"),
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithOption.WithFloatEncoding("IEEE754"));
        try {
            if (quantity == 1) {
                if (value == null || value.isList()) {
                    return null;
                }
                writeOne(buffer, value);
            } else {
                if (value == null || !value.isList() || value.getList().size() != quantity) {
                    return null;
                }
                for (PlcValue element : value.getList()) {
                    writeOne(buffer, element);
                }
            }
            return buffer.getBytes();
        } catch (PlcIncompatibleDatatypeException e) {
            // An un-coercible value is an expected caller error (mapped to INVALID_DATA), not an
            // internal fault -- keep it out of the WARN log and skip the stack trace.
            LOGGER.debug("Incompatible value for SLMP {} encode: {}", this, e.getMessage());
            return null;
        } catch (BufferException | PlcRuntimeException e) {
            LOGGER.warn("Failed to encode SLMP {} value", this, e);
            return null;
        }
    }

    private void writeOne(WriteBufferByteBased buffer, PlcValue value) throws BufferException {
        switch (this) {
            case WORD:
            case UINT:
                // writeUnsignedInt(16, int) accepts the full 0..65535 range; writeUnsignedShort(16, short)
                // takes a SIGNED short and would reject 65535 (0xFFFF) as negative, corrupting unsigned WORD/UINT.
                buffer.writeUnsignedInt(16, value.getInt());
                break;
            case INT:
                buffer.writeSignedShort(16, value.getShort());
                break;
            case DINT:
                buffer.writeSignedInt(32, value.getInt());
                break;
            case UDINT:
                buffer.writeUnsignedLong(32, value.getLong());
                break;
            case REAL:
                buffer.writeFloat(32, value.getFloat());
                break;
            default:
                throw new BufferException("Unsupported SLMP data type " + this);
        }
    }

    private PlcValue readOne(ReadBufferByteBased buffer) throws BufferException {
        switch (this) {
            case WORD:
                // readUnsignedInt(16) returns int 0..65535; readUnsignedShort(16) returns a SIGNED short
                // (negative for values > 0x7FFF), which would corrupt unsigned WORD/UINT.
                return new PlcWORD(buffer.readUnsignedInt(16));
            case INT:
                return new PlcINT(buffer.readSignedShort(16));
            case UINT:
                return new PlcUINT(buffer.readUnsignedInt(16));
            case DINT:
                return new PlcDINT(buffer.readSignedInt(32));
            case UDINT:
                return new PlcUDINT(buffer.readUnsignedLong(32));
            case REAL:
                return new PlcREAL(buffer.readFloat(32));
            default:
                throw new BufferException("Unsupported SLMP data type " + this);
        }
    }
}
