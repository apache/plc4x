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

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.byteorder.ByteOrder;
import org.apache.plc4x.java.spi.buffers.bytebased.byteorder.ByteOrderBigEndian;
import org.apache.plc4x.java.spi.buffers.bytebased.byteorder.ByteOrderManager;
import org.apache.plc4x.java.spi.buffers.bytebased.encoding.Encoding;
import org.apache.plc4x.java.spi.buffers.bytebased.encoding.EncodingDefault;
import org.apache.plc4x.java.spi.buffers.bytebased.encoding.EncodingManager;
import org.apache.plc4x.java.spi.buffers.bytebased.encoding.EncodingRaw;
import org.apache.plc4x.java.spi.buffers.bytebased.encoding.EncodingTwosComplement;
import org.apache.plc4x.java.spi.buffers.bytebased.encoding.EncodingUnsignedBinary;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.Optional;

public class WriteBufferByteBased extends AbstractBufferByteBased implements WriteBuffer {

    public WriteBufferByteBased(byte[] buffer, WithOption... defaultOptions) {
        this(buffer, 0, buffer.length * 8, defaultOptions);
    }

    public WriteBufferByteBased(byte[] buffer, int startBit, int sizeInBits, WithOption... defaultOptions) {
        super(buffer, startBit, sizeInBits,
            WithByteBasedOption.extractByteOrderManager(defaultOptions).orElseGet(ByteOrderManager::getDefault),
            WithByteBasedOption.extractEncodingManager(defaultOptions).orElseGet(EncodingManager::getDefault),
            defaultOptions);
    }

    @Override
    public void writeBit(boolean value, WithOption... options) throws BufferException {
        ensureAvailable(1);
        // Index by the ABSOLUTE bit position (startBit + positionInBits), like readBit and the
        // whole-byte arraycopy path in writeBits — a buffer constructed with a non-zero startBit
        // must not write from the start of the backing array.
        int absoluteBitIndex = startBit + positionInBits;
        int byteIndex = absoluteBitIndex / 8;
        int bitIndex = absoluteBitIndex % 8;
        if (value) {
            buffer[byteIndex] |= (byte) (0x80 >> bitIndex);
        }
        positionInBits++;
    }

    @Override
    public void writeBits(int numBits, byte[] value, WithOption... options) throws BufferException {
        if (numBits < 0) {
            throw new BufferException("Number of bits must be positive");
        } else if (numBits == 0) {
            return;
        }
        ensureAvailable(numBits);

        // Simple case, in which we're writing full bytes and we are byte-aligned.
        if (isAligned() && (numBits % 8 == 0)) {
            int byteIndex = (startBit + positionInBits) / 8;
            System.arraycopy(value, 0, buffer, byteIndex, numBits / 8);
            positionInBits += numBits;
        } else {
            int valueBitIndex = (8 - (numBits % 8)) % 8;
            for (int i = 0; i < numBits; i++) {
                int valueByteIndex = valueBitIndex / 8;
                int valueBitPosition = 7 - (valueBitIndex % 8);
                boolean bit = (value[valueByteIndex] & (1 << valueBitPosition)) != 0;
                writeBit(bit);
                valueBitIndex++;
            }
        }
    }

    /**
     * Writes the low {@code numBits} (a whole number of bytes) of {@code value} big-endian into the
     * backing array at the current byte-aligned position and advances the position. Callers gate this
     * behind the fast-path eligibility checks; no intermediate byte[] or ByteOrder object is allocated.
     * Unlike {@code readAlignedBytesBE} this method performs NO capacity check of its own — every
     * caller invokes {@code ensureAvailable(numBits)} before the fast-path branch, and that call must
     * not be removed as "redundant" with the slow path's.
     */
    private void writeAlignedBytesBE(int numBits, long value) {
        int byteIndex = (startBit + positionInBits) / 8;
        int numBytes = numBits / 8;
        for (int i = 0; i < numBytes; i++) {
            buffer[byteIndex + i] = (byte) ((value >>> ((numBytes - 1 - i) * 8)) & 0xFF);
        }
        positionInBits += numBits;
    }

    @Override
    public void writeUnsignedByte(int numBits, byte value, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 7) {
            throw new BufferException("Unsigned byte can only read between 1 and 7 bits");
        }
        // Check if the value is within the valid range for the given number of bits
        int maxValue = numBits == 7 ? 0xFF : (1 << numBits) - 1;
        if (value < 0 || value > maxValue) {
            throw new BufferException("Value " + value + " is out of range for " + numBits + " bits. Valid range is 0 to " + maxValue);
        }
        ensureAvailable(numBits);

        Optional<Encoding> encodingOptional = getUnsignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for unsigned integer values"));
        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = encodingDefault.encodeByte(numBits, value);
            bytes = getByteOrder(options).process(bytes);
            writeBits(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            byte[] bytes = encodingRaw.encodeByte(numBits, value);
            writeBits(bytes.length * 8, bytes);
        } else {
            throw new BufferException("Unsupported encoding type: " + encoding.getClass().getName());
        }
    }

    @Override
    public void writeUnsignedShort(int numBits, short value, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 15) {
            throw new BufferException("Unsigned short can only read between 1 and 15 bits");
        }
        // Check if the value is within the valid range for the given number of bits
        long maxValue = numBits == 15 ? 0xFFFFL : (1L << numBits) - 1;
        if (value < 0 || value > maxValue) {
            throw new BufferException("Value " + value + " is out of range for " + numBits + " bits. Valid range is 0 to " + maxValue);
        }
        ensureAvailable(numBits);

        Optional<Encoding> encodingOptional = getUnsignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for unsigned integer values"));
        ByteOrder byteOrder = getByteOrder(options);
        if (isByteAlignedWholeBytes(numBits, options)
            && encoding.getClass() == EncodingUnsignedBinary.class
            && byteOrder.getClass() == ByteOrderBigEndian.class) {
            writeAlignedBytesBE(numBits, value);
            return;
        }

        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = encodingDefault.encodeShort(numBits, value);
            bytes = byteOrder.process(bytes);
            writeBits(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            byte[] bytes = encodingRaw.encodeShort(numBits, value);
            writeBits(bytes.length * 8, bytes);
        } else {
            throw new BufferException("Unsupported encoding type: " + encoding.getClass().getName());
        }
    }

    @Override
    public void writeUnsignedInt(int numBits, int value, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 31) {
            throw new BufferException("Unsigned int can only read between 1 and 31 bits");
        }
        // Check if the value is within the valid range for the given number of bits
        long maxValue = numBits == 31 ? 0xFFFFFFFFL : (1L << numBits) - 1;
        if (value < 0 || value > maxValue) {
            throw new BufferException("Value " + value + " is out of range for " + numBits + " bits. Valid range is 0 to " + maxValue);
        }
        ensureAvailable(numBits);

        Optional<Encoding> encodingOptional = getUnsignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for unsigned integer values"));
        ByteOrder byteOrder = getByteOrder(options);
        // Fast path: byte-aligned whole-byte plain-binary big-endian write straight into the backing
        // array; encoding/byte order are resolved once above and reused by the slow path below.
        if (isByteAlignedWholeBytes(numBits, options)
            && encoding.getClass() == EncodingUnsignedBinary.class
            && byteOrder.getClass() == ByteOrderBigEndian.class) {
            writeAlignedBytesBE(numBits, value);
            return;
        }

        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = encodingDefault.encodeInt(numBits, value);
            bytes = byteOrder.process(bytes);
            writeBits(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            byte[] bytes = encodingRaw.encodeInt(numBits, value);
            writeBits(bytes.length * 8, bytes);
        } else {
            throw new BufferException("Unsupported encoding type: " + encoding.getClass().getName());
        }
    }

    @Override
    public void writeUnsignedLong(int numBits, long value, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 63) {
            throw new BufferException("Unsigned long can only read between 1 and 63 bits");
        }
        // Check if the value is within the valid range for the given number of bits
        long maxValue = numBits == 63 ? 0x7FFFFFFFFFFFFFFFL : (1L << numBits) - 1;
        if (value < 0 || value > maxValue) {
            throw new BufferException("Value " + value + " is out of range for " + numBits + " bits. Valid range is 0 to " + maxValue);
        }
        ensureAvailable(numBits);

        Optional<Encoding> encodingOptional = getUnsignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for unsigned integer values"));
        ByteOrder byteOrder = getByteOrder(options);
        if (isByteAlignedWholeBytes(numBits, options)
            && encoding.getClass() == EncodingUnsignedBinary.class
            && byteOrder.getClass() == ByteOrderBigEndian.class) {
            writeAlignedBytesBE(numBits, value);
            return;
        }

        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = encodingDefault.encodeLong(numBits, value);
            bytes = byteOrder.process(bytes);
            writeBits(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            byte[] bytes = encodingRaw.encodeLong(numBits, value);
            writeBits(bytes.length * 8, bytes);
        } else {
            throw new BufferException("Unsupported encoding type: " + encoding.getClass().getName());
        }
    }

    @Override
    public void writeUnsignedBigInteger(int numBits, BigInteger value, WithOption... options) throws BufferException {
        if (numBits < 1) {
            throw new BufferException("Unsigned BigInteger can only read between 1 or more bits");
        }
        // Check if the value is within the valid range for the given number of bits
        BigInteger maxValue = BigInteger.ONE.shiftLeft(numBits).subtract(BigInteger.ONE);
        if (value == null || value.compareTo(BigInteger.ZERO) < 0 || value.compareTo(maxValue) > 0) {
            throw new BufferException("Value " + value + " is out of range for " + numBits + " bits. Valid range is 0 to " + maxValue);
        }

        Optional<Encoding> encodingOptional = getUnsignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for unsigned integer values"));
        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = encodingDefault.encodeBigInteger(numBits, value);
            bytes = getByteOrder(options).process(bytes);
            writeBits(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            byte[] bytes = encodingRaw.encodeBigInteger(numBits, value);
            writeBits(bytes.length * 8, bytes);
        } else {
            throw new BufferException("Unsupported encoding type: " + encoding.getClass().getName());
        }
    }

    @Override
    public void writeSignedByte(int numBits, byte value, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 8) {
            throw new BufferException("Signed byte can only read between 1 and 8 bits");
        }
        // Check if the value is within the valid range for the given number of bits
        int minValue = numBits == 8 ? Byte.MIN_VALUE : -(1 << (numBits - 1));
        int maxValue = numBits == 8 ? Byte.MAX_VALUE : (1 << (numBits - 1)) - 1;
        if (value < minValue || value > maxValue) {
            throw new BufferException("Value " + value + " is out of range for " + numBits + " bits. Valid range is " + minValue + " to " + maxValue);
        }
        ensureAvailable(numBits);

        Optional<Encoding> encodingOptional = getSignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for signed integer values"));
        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = encodingDefault.encodeByte(numBits, value);
            bytes = getByteOrder(options).process(bytes);
            writeBits(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            byte[] bytes = encodingRaw.encodeByte(numBits, value);
            writeBits(bytes.length * 8, bytes);
        } else {
            throw new BufferException("Unsupported encoding type: " + encoding.getClass().getName());
        }
    }

    @Override
    public void writeSignedShort(int numBits, short value, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 16) {
            throw new BufferException("Signed short can only read between 1 and 16 bits");
        }
        // Check if the value is within the valid range for the given number of bits
        int minValue = numBits == 16 ? Short.MIN_VALUE : -(1 << (numBits - 1));
        int maxValue = numBits == 16 ? Short.MAX_VALUE : (1 << (numBits - 1)) - 1;
        if (value < minValue || value > maxValue) {
            throw new BufferException("Value " + value + " is out of range for " + numBits + " bits. Valid range is " + minValue + " to " + maxValue);
        }
        ensureAvailable(numBits);

        Optional<Encoding> encodingOptional = getSignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for signed integer values"));
        ByteOrder byteOrder = getByteOrder(options);
        if (isByteAlignedWholeBytes(numBits, options)
            && encoding.getClass() == EncodingTwosComplement.class
            && byteOrder.getClass() == ByteOrderBigEndian.class) {
            writeAlignedBytesBE(numBits, value);
            return;
        }

        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = encodingDefault.encodeShort(numBits, value);
            bytes = byteOrder.process(bytes);
            writeBits(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            byte[] bytes = encodingRaw.encodeShort(numBits, value);
            writeBits(bytes.length * 8, bytes);
        } else {
            throw new BufferException("Unsupported encoding type: " + encoding.getClass().getName());
        }
    }

    @Override
    public void writeSignedInt(int numBits, int value, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 32) {
            throw new BufferException("Signed int can only read between 1 and 32 bits");
        }
        // Check if the value is within the valid range for the given number of bits
        long minValue = numBits == 32 ? Integer.MIN_VALUE : -(1L << (numBits - 1));
        long maxValue = numBits == 32 ? Integer.MAX_VALUE : (1L << (numBits - 1)) - 1;
        if (value < minValue || value > maxValue) {
            throw new BufferException("Value " + value + " is out of range for " + numBits + " bits. Valid range is " + minValue + " to " + maxValue);
        }
        ensureAvailable(numBits);

        Optional<Encoding> encodingOptional = getSignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for signed integer values"));
        ByteOrder byteOrder = getByteOrder(options);
        if (isByteAlignedWholeBytes(numBits, options)
            && encoding.getClass() == EncodingTwosComplement.class
            && byteOrder.getClass() == ByteOrderBigEndian.class) {
            writeAlignedBytesBE(numBits, value);
            return;
        }

        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = encodingDefault.encodeInt(numBits, value);
            bytes = byteOrder.process(bytes);
            writeBits(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            byte[] bytes = encodingRaw.encodeInt(numBits, value);
            writeBits(bytes.length * 8, bytes);
        } else {
            throw new BufferException("Unsupported encoding type: " + encoding.getClass().getName());
        }
    }

    @Override
    public void writeSignedLong(int numBits, long value, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 64) {
            throw new BufferException("Signed long can only read between 1 and 64 bits");
        }
        // Check if the value is within the valid range for the given number of bits
        long minValue = numBits == 64 ? Long.MIN_VALUE : -(1L << (numBits - 1));
        long maxValue = numBits == 64 ? Long.MAX_VALUE : (1L << (numBits - 1)) - 1;
        if (value < minValue || value > maxValue) {
            throw new BufferException("Value " + value + " is out of range for " + numBits + " bits. Valid range is " + minValue + " to " + maxValue);
        }
        ensureAvailable(numBits);

        Optional<Encoding> encodingOptional = getSignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for signed integer values"));
        ByteOrder byteOrder = getByteOrder(options);
        if (isByteAlignedWholeBytes(numBits, options)
            && encoding.getClass() == EncodingTwosComplement.class
            && byteOrder.getClass() == ByteOrderBigEndian.class) {
            writeAlignedBytesBE(numBits, value);
            return;
        }

        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = encodingDefault.encodeLong(numBits, value);
            bytes = byteOrder.process(bytes);
            writeBits(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            byte[] bytes = encodingRaw.encodeLong(numBits, value);
            writeBits(bytes.length * 8, bytes);
        } else {
            throw new BufferException("Unsupported encoding type: " + encoding.getClass().getName());
        }
    }

    public void writeSignedBigInteger(int numBits, BigInteger value, WithOption... options) throws BufferException {
        if (numBits < 1) {
            throw new BufferException("Signed BigInteger can only read between 1 or more bits");
        }
        // Check if the value is within the valid range for the given number of bits
        BigInteger minValue = BigInteger.ONE.shiftLeft(numBits - 1).negate();
        BigInteger maxValue = BigInteger.ONE.shiftLeft(numBits - 1).subtract(BigInteger.ONE);
        if (value == null || value.compareTo(minValue) < 0 || value.compareTo(maxValue) > 0) {
            throw new BufferException("Value " + value + " is out of range for " + numBits + " bits. Valid range is " + minValue + " to " + maxValue);
        }
        ensureAvailable(numBits);

        Optional<Encoding> encodingOptional = getSignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for signed integer values"));
        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = encodingDefault.encodeBigInteger(numBits, value);
            bytes = getByteOrder(options).process(bytes);
            writeBits(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            byte[] bytes = encodingRaw.encodeBigInteger(numBits, value);
            writeBits(bytes.length * 8, bytes);
        } else {
            throw new BufferException("Unsupported encoding type: " + encoding.getClass().getName());
        }
    }

    @Override
    public void writeFloat(int numBits, float value, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 32) {
            throw new BufferException("Float can only be written using between 1 and 32 bits");
        }
        ensureAvailable(numBits);

        Optional<Encoding> encodingOptional = getFloatEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for floating point values"));
        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = encodingDefault.encodeFloat(numBits, value);
            bytes = getByteOrder(options).process(bytes);
            writeBits(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            byte[] bytes = encodingRaw.encodeFloat(numBits, value);
            writeBits(bytes.length * 8, bytes);
        } else {
            throw new BufferException("Unsupported encoding type: " + encoding.getClass().getName());
        }
    }

    @Override
    public void writeDouble(int numBits, double value, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 64) {
            throw new BufferException("Double can only be written using between 1 and 64 bits");
        }
        ensureAvailable(numBits);

        Optional<Encoding> encodingOptional = getFloatEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for floating point values"));
        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = encodingDefault.encodeDouble(numBits, value);
            bytes = getByteOrder(options).process(bytes);
            writeBits(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            byte[] bytes = encodingRaw.encodeDouble(numBits, value);
            writeBits(bytes.length * 8, bytes);
        } else {
            throw new BufferException("Unsupported encoding type: " + encoding.getClass().getName());
        }
    }

    @Override
    public void writeBigDecimal(int numBits, BigDecimal value, WithOption... options) throws BufferException {
        if (numBits < 1) {
            throw new BufferException("BigDecimal can only be written using 1 or more bits");
        }
        if (value == null) {
            throw new BufferException("Value null is out of range for " + numBits + " bits.");
        }
        ensureAvailable(numBits);

        Optional<Encoding> encodingOptional = getFloatEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for floating point values"));
        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = encodingDefault.encodeBigDecimal(numBits, value);
            bytes = getByteOrder(options).process(bytes);
            writeBits(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            byte[] bytes = encodingRaw.encodeBigDecimal(numBits, value);
            writeBits(bytes.length * 8, bytes);
        } else {
            throw new BufferException("Unsupported encoding type: " + encoding.getClass().getName());
        }
    }

    @Override
    public void writeString(int numBits, String value, WithOption... options) throws BufferException {
        if (numBits == 0) {
            return;
        }
        if (numBits < 1) {
            throw new BufferException("Number of bits must be written using 1 or more bits");
        }
        if (value == null) {
            throw new BufferException("Value null is out of range for " + numBits + " bits.");
        }
        ensureAvailable(numBits);

        // Convert string to bytes
        Optional<Encoding> encoding = getStringEncoding(options);
        // TODO: Possibly add support for the Raw encoding here ...
        byte[] bytes = encoding.orElseThrow(() -> new BufferException("No encoding defined for string values")).encodeString(numBits, value);

        // If padding is configured, update the 0x00 bytes to the provided character value
        Optional<Character> paddingChar = WithByteBasedOption.extractPaddingChar(options, getContext());
        if (paddingChar.isPresent()) {
            Character c = paddingChar.get();
            for (int i = 0; i < bytes.length; i++) {
                if (bytes[i] == 0x00) {
                    bytes[i] = (byte) c.charValue();
                }
            }
        }

        // If the string's byte representation is longer than the specified number of bits,
        // truncate it to fit
        int maxBytes = (numBits + 7) / 8;
        if (bytes.length > maxBytes) {
            byte[] truncatedBytes = new byte[maxBytes];
            System.arraycopy(bytes, 0, truncatedBytes, 0, maxBytes);
            bytes = truncatedBytes;
        }

        // If the string's byte representation is shorter than the specified number of bits,
        // pad it with zeros
        if (bytes.length * 8 < numBits) {
            int paddedLength = (numBits + 7) / 8;
            byte[] paddedBytes = new byte[paddedLength];
            System.arraycopy(bytes, 0, paddedBytes, 0, bytes.length);
            bytes = paddedBytes;
        }

        // Do not apply byte order to UTF-8 strings as it would reverse the bytes
        // and make the string unreadable
        writeBits(numBits, bytes);
    }

    @Override
    public WriteBuffer createSubBuffer(int numBits, WithOption... options) throws BufferException {
        ensureAvailable(numBits);
        int lengthInBytes = (numBits + 7) / 8;
        ByteOrder byteOrder = getByteOrder(options);
        WithOption byteOrderOption = WithByteBasedOption.WithByteOrder(byteOrder.getName());
        byte[] subBufferBytes = new byte[lengthInBytes];
        return new WriteBufferByteBased(subBufferBytes, byteOrderOption);
    }

    @Override
    public boolean isByteBased() {
        return true;
    }

}
