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

import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
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

public class ReadBufferByteBased extends AbstractBufferByteBased implements ReadBuffer, ReadBufferRaw {

    public ReadBufferByteBased(byte[] buffer, WithOption... options) {
        this(buffer, 0, buffer.length * 8, options);
    }

    private ReadBufferByteBased(byte[] buffer, int startBit, int sizeInBits, WithOption... options) {
        super(buffer, startBit, sizeInBits,
            WithByteBasedOption.extractByteOrderManager(options).orElseGet(ByteOrderManager::getDefault),
            WithByteBasedOption.extractEncodingManager(options).orElseGet(EncodingManager::getDefault),
            options);
    }

    @Override
    public boolean readBit(WithOption... options) throws BufferException {
        ensureAvailable(1);
        int absoluteBitIndex = startBit + positionInBits;
        int byteIndex = absoluteBitIndex / 8;
        int bitIndex = absoluteBitIndex % 8;
        boolean bit = (buffer[byteIndex] & (0x80 >> bitIndex)) != 0;
        positionInBits++;
        return bit;
    }

    @Override
    public byte[] readBits(int numBits, WithOption... options) throws BufferException {
        if (numBits < 0) {
            throw new BufferException("Number of bits must be positive");
        } else if (numBits == 0) {
            return new byte[0];
        }
        ensureAvailable(numBits);

        // Calculate the number of bytes needed to store the bits
        int numBytes = (numBits + 7) / 8;
        byte[] result = new byte[numBytes];

        // Simple case, in which we're reading full bytes and we are byte-aligned.
        // Using this option is significantly quicker.
        if (isAligned() && (numBits % 8 == 0)) {
            // Fast path for byte-aligned reads
            int byteIndex = (startBit + positionInBits) / 8;
            System.arraycopy(buffer, byteIndex, result, 0, numBytes);
            positionInBits += numBits;
            return result;
        } else {
            // Potentially align when doing partial byte reads.
            int resultBitIndex = (8 - (numBits % 8)) % 8;
            for (int i = 0; i < numBits; i++) {
                if (readBit()) {
                    int resultByteIndex = resultBitIndex / 8;
                    int resultBitPosition = 7 - (resultBitIndex % 8);
                    result[resultByteIndex] |= (byte) (1 << resultBitPosition);
                }
                resultBitIndex++;
            }
            return result;
        }
    }

    /**
     * Reads {@code numBits} (a whole number of bytes) big-endian from the current byte-aligned
     * position into a long and advances the position. Callers gate this behind the fast-path
     * eligibility checks; no intermediate byte[] or per-field ByteOrder object is allocated.
     */
    private long readAlignedBytesBE(int numBits) throws BufferException {
        ensureAvailable(numBits);
        int byteIndex = (startBit + positionInBits) / 8;
        int numBytes = numBits / 8;
        long v = 0;
        for (int i = 0; i < numBytes; i++) {
            v = (v << 8) | (buffer[byteIndex + i] & 0xFF);
        }
        positionInBits += numBits;
        return v;
    }

    @Override
    public byte readUnsignedByte(int numBits, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 7) {
            throw new BufferException("Unsigned byte can only read between 1 and 7 bits");
        }

        Optional<Encoding> encodingOptional = getUnsignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for unsigned integer values"));
        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = readBits(numBits);
            bytes = getByteOrder(options).process(bytes);
            return encodingDefault.decodeByte(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            return encodingRaw.decodeByte(numBits, this);
        }
        throw new BufferException("Unsupported encoding type: " +  encoding.getClass().getName());
    }

    @Override
    public short readUnsignedShort(int numBits, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 15) {
            throw new BufferException("Unsigned short can only read between 1 and 15 bits");
        }

        Optional<Encoding> encodingOptional = getUnsignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for unsigned integer values"));
        ByteOrder byteOrder = getByteOrder(options);
        if (isByteAlignedWholeBytes(numBits, options)
            && encoding instanceof EncodingUnsignedBinary
            && byteOrder instanceof ByteOrderBigEndian) {
            return (short) readAlignedBytesBE(numBits);
        }

        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = readBits(numBits);
            bytes = byteOrder.process(bytes);
            return encodingDefault.decodeShort(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            return encodingRaw.decodeShort(numBits, this);
        }
        throw new BufferException("Unsupported encoding type: " +  encoding.getClass().getName());
    }

    @Override
    public int readUnsignedInt(int numBits, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 31) {
            throw new BufferException("Unsigned int can only read between 1 and 31 bits");
        }

        Optional<Encoding> encodingOptional = getUnsignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for unsigned integer values"));
        ByteOrder byteOrder = getByteOrder(options);
        // Fast path: byte-aligned whole-byte plain-binary big-endian read straight from the backing
        // array; encoding/byte order are resolved once above and reused by the slow path below.
        if (isByteAlignedWholeBytes(numBits, options)
            && encoding instanceof EncodingUnsignedBinary
            && byteOrder instanceof ByteOrderBigEndian) {
            return (int) readAlignedBytesBE(numBits);
        }

        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = readBits(numBits);
            bytes = byteOrder.process(bytes);
            return encodingDefault.decodeInt(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            return encodingRaw.decodeInt(numBits, this);
        }
        throw new BufferException("Unsupported encoding type: " +  encoding.getClass().getName());
    }

    @Override
    public long readUnsignedLong(int numBits, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 63) {
            throw new BufferException("Unsigned long can only read between 1 and 63 bits");
        }

        Optional<Encoding> encodingOptional = getUnsignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for unsigned integer values"));
        ByteOrder byteOrder = getByteOrder(options);
        if (isByteAlignedWholeBytes(numBits, options)
            && encoding instanceof EncodingUnsignedBinary
            && byteOrder instanceof ByteOrderBigEndian) {
            return readAlignedBytesBE(numBits);
        }

        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = readBits(numBits);
            bytes = byteOrder.process(bytes);
            return encodingDefault.decodeLong(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            return encodingRaw.decodeLong(numBits, this);
        }
        throw new BufferException("Unsupported encoding type: " +  encoding.getClass().getName());
    }

    @Override
    public BigInteger readUnsignedBigInteger(int numBits, WithOption... options) throws BufferException {
        if (numBits < 1) {
            throw new BufferException("Unsigned BigInteger can only read 1 or more bits");
        }

        Optional<Encoding> encodingOptional = getUnsignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for unsigned integer values"));
        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = readBits(numBits);
            bytes = getByteOrder(options).process(bytes);
            return encodingDefault.decodeBigInteger(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            return encodingRaw.decodeBigInteger(numBits, this);
        }
        throw new BufferException("Unsupported encoding type: " +  encoding.getClass().getName());
    }

    @Override
    public byte readSignedByte(int numBits, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 8) {
            throw new BufferException("Signed byte can only read between 1 and 8 bits");
        }

        Optional<Encoding> encodingOptional = getSignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for signed integer values"));
        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = readBits(numBits);
            bytes = getByteOrder(options).process(bytes);
            return encodingDefault.decodeByte(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            return encodingRaw.decodeByte(numBits, this);
        }
        throw new BufferException("Unsupported encoding type: " +  encoding.getClass().getName());
    }

    @Override
    public short readSignedShort(int numBits, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 16) {
            throw new BufferException("Signed short can only read between 1 and 16 bits");
        }

        Optional<Encoding> encodingOptional = getSignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for signed integer values"));
        ByteOrder byteOrder = getByteOrder(options);
        if (isByteAlignedWholeBytes(numBits, options)
            && encoding instanceof EncodingTwosComplement
            && byteOrder instanceof ByteOrderBigEndian) {
            return (short) signExtend(readAlignedBytesBE(numBits), numBits);
        }

        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = readBits(numBits);
            bytes = byteOrder.process(bytes);
            return encodingDefault.decodeShort(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            return encodingRaw.decodeShort(numBits, this);
        }
        throw new BufferException("Unsupported encoding type: " +  encoding.getClass().getName());
    }

    @Override
    public int readSignedInt(int numBits, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 32) {
            throw new BufferException("Signed int can only read between 1 and 32 bits");
        }

        Optional<Encoding> encodingOptional = getSignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for signed integer values"));
        ByteOrder byteOrder = getByteOrder(options);
        if (isByteAlignedWholeBytes(numBits, options)
            && encoding instanceof EncodingTwosComplement
            && byteOrder instanceof ByteOrderBigEndian) {
            return (int) signExtend(readAlignedBytesBE(numBits), numBits);
        }

        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = readBits(numBits);
            bytes = byteOrder.process(bytes);
            return encodingDefault.decodeInt(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            return encodingRaw.decodeInt(numBits, this);
        }
        throw new BufferException("Unsupported encoding type: " +  encoding.getClass().getName());
    }

    @Override
    public long readSignedLong(int numBits, WithOption... options) throws BufferException {
        if (numBits < 1 || numBits > 64) {
            throw new BufferException("Signed long can only read between 1 and 64 bits");
        }

        Optional<Encoding> encodingOptional = getSignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for signed integer values"));
        ByteOrder byteOrder = getByteOrder(options);
        if (isByteAlignedWholeBytes(numBits, options)
            && encoding instanceof EncodingTwosComplement
            && byteOrder instanceof ByteOrderBigEndian) {
            return signExtend(readAlignedBytesBE(numBits), numBits);
        }

        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = readBits(numBits);
            bytes = byteOrder.process(bytes);
            return encodingDefault.decodeLong(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            return encodingRaw.decodeLong(numBits, this);
        }
        throw new BufferException("Unsupported encoding type: " +  encoding.getClass().getName());
    }

    @Override
    public BigInteger readSignedBigInteger(int numBits, WithOption... options) throws BufferException {
        if (numBits < 1) {
            throw new BufferException("Signed BigInteger can only read 1 or more bits");
        }

        // For signed BigInteger, we can directly use the BigInteger constructor
        // that takes a byte array, which interprets it as a two's-complement representation
        Optional<Encoding> encodingOptional = getSignedIntegerEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for signed integer values"));
        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = readBits(numBits);
            bytes = getByteOrder(options).process(bytes);
            return encodingDefault.decodeBigInteger(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            return encodingRaw.decodeBigInteger(numBits, this);
        }
        throw new BufferException("Unsupported encoding type: " +  encoding.getClass().getName());
    }

    @Override
    public float readFloat(int numBits, WithOption... options) throws BufferException {
        Optional<Encoding> encodingOptional = getFloatEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for float values"));
        // IEEE 754 only supports 32-bit floats; other encodings (e.g. KNXFloat,
        // which is a 16-bit format) decide their own size constraints inside
        // their decodeFloat implementation.
        if (encoding instanceof EncodingDefault encodingDefault) {
            if ("IEEE754".equals(encoding.getName()) && numBits != 32) {
                throw new BufferException("IEEE 754 floats can only be read using 32 bits");
            }
            ensureAvailable(numBits);
            byte[] bytes = readBits(numBits);
            bytes = getByteOrder(options).process(bytes);
            return encodingDefault.decodeFloat(numBits, bytes);
        } else if (encoding instanceof EncodingRaw encodingRaw) {
            return encodingRaw.decodeFloat(numBits, this);
        }
        throw new BufferException("Unsupported encoding type: " +  encoding.getClass().getName());
    }

    @Override
    public double readDouble(int numBits, WithOption... options) throws BufferException {
        if (numBits != 64) {
            throw new BufferException("Double can only be read using 64 bits");
        }

        Optional<Encoding> encodingOptional = getFloatEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for signed integer values"));
        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = readBits(numBits);
            bytes = getByteOrder(options).process(bytes);
            return encodingDefault.decodeDouble(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            return encodingRaw.decodeDouble(numBits, this);
        }
        throw new BufferException("Unsupported encoding type: " +  encoding.getClass().getName());
    }

    @Override
    public BigDecimal readBigDecimal(int numBits, WithOption... options) throws BufferException {
        if (numBits < 1) {
            throw new BufferException("BigDecimal can only be read using 1 or more bits");
        }

        Optional<Encoding> encodingOptional = getFloatEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for signed integer values"));
        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = readBits(numBits);
            bytes = getByteOrder(options).process(bytes);
            return encodingDefault.decodeBigDecimal(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            return encodingRaw.decodeBigDecimal(numBits, this);
        }
        throw new BufferException("Unsupported encoding type: " +  encoding.getClass().getName());
    }

    @Override
    public String readString(int numBits, WithOption... options) throws BufferException {
        if (numBits < 0) {
            throw new BufferException("Number of bits must be non-negative");
        }
        if (numBits == 0) {
            return "";
        }

        Optional<Encoding> encodingOptional = getStringEncoding(options);
        Encoding encoding = encodingOptional.orElseThrow(() -> new BufferException("No encoding defined for string values"));
        if(encoding instanceof EncodingDefault encodingDefault) {
            ensureAvailable(numBits);
            byte[] bytes = readBits(numBits);
            // Do not apply byte order to UTF-8 strings as it would reverse the bytes
            // and make the string unreadable
            return encodingDefault.decodeString(numBits, bytes);
        } else if(encoding instanceof EncodingRaw encodingRaw) {
            return encodingRaw.decodeString(numBits, this);
        }
        throw new BufferException("Unsupported encoding type: " +  encoding.getClass().getName());
    }

    @Override
    public ReadBufferByteBased createSubBuffer(int numBits, WithOption... options) throws BufferException {
        ensureAvailable(numBits);
        ReadBufferByteBased subBuffer = new ReadBufferByteBased(
            buffer,
            startBit + positionInBits,
            numBits,
            WithOption.AddOptions(getContext(), options)
        );
        positionInBits += numBits;
        return subBuffer;
    }

}
