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
package org.apache.plc4x.java.spi.buffers.api;

import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;

import java.math.BigDecimal;
import java.math.BigInteger;

/**
 * Interface for writing data to a buffer with bit-level precision.
 * Provides methods for writing various data types with specified bit lengths.
 */
public interface WriteBuffer extends Buffer {

    /**
     * Writes a single bit to the buffer.
     *
     * @param value the bit value to write (true for 1, false for 0)
     */
    void writeBit(boolean value, WithOption... options) throws BufferException;

    /**
     * Writes a specified number of bits to the buffer.
     *
     * @param numBits the number of bits to write
     * @param value   the byte array containing the bits to write
     */
    void writeBits(int numBits, byte[] value, WithOption... options) throws BufferException;

    /**
     * Writes an unsigned byte value with the specified bit length.
     *
     * @param numBits the number of bits to write
     * @param value   the unsigned byte value as a byte value
     * @param options additional options for writing
     */
    void writeUnsignedByte(int numBits, byte value, WithOption... options) throws BufferException;

    /**
     * Writes an unsigned short value with the specified bit length.
     *
     * @param numBits the number of bits to write
     * @param value   the unsigned short value as an short value
     * @param options additional options for writing
     */
    void writeUnsignedShort(int numBits, short value, WithOption... options) throws BufferException;

    /**
     * Writes an unsigned int value with the specified bit length.
     *
     * @param numBits the number of bits to write
     * @param value   the unsigned int value as a int value
     * @param options additional options for writing
     */
    void writeUnsignedInt(int numBits, int value, WithOption... options) throws BufferException;

    /**
     * Writes an unsigned long value with the specified bit length.
     *
     * @param numBits the number of bits to write
     * @param value   the unsigned long value as a long value
     * @param options additional options for writing
     */
    void writeUnsignedLong(int numBits, long value, WithOption... options) throws BufferException;

    /**
     * Writes an unsigned big integer with the specified bit length.
     *
     * @param numBits the number of bits to write
     * @param value   the unsigned BigInteger value
     * @param options additional options for writing
     */
    void writeUnsignedBigInteger(int numBits, BigInteger value, WithOption... options) throws BufferException;

    /**
     * Writes a signed byte value with the specified bit length.
     *
     * @param numBits the number of bits to write
     * @param value   the signed byte value
     * @param options additional options for writing
     */
    void writeSignedByte(int numBits, byte value, WithOption... options) throws BufferException;

    /**
     * Writes a signed short value with the specified bit length.
     *
     * @param numBits the number of bits to write
     * @param value   the signed short value
     * @param options additional options for writing
     */
    void writeSignedShort(int numBits, short value, WithOption... options) throws BufferException;

    /**
     * Writes a signed int value with the specified bit length.
     *
     * @param numBits the number of bits to write
     * @param value   the signed int value
     * @param options additional options for writing
     */
    void writeSignedInt(int numBits, int value, WithOption... options) throws BufferException;

    /**
     * Writes a signed long value with the specified bit length.
     *
     * @param numBits the number of bits to write
     * @param value   the signed long value
     * @param options additional options for writing
     */
    void writeSignedLong(int numBits, long value, WithOption... options) throws BufferException;

    /**
     * Writes a signed big integer with the specified bit length.
     *
     * @param numBits the number of bits to write
     * @param value   the signed big integer value
     * @param options additional options for writing
     */
    void writeSignedBigInteger(int numBits, BigInteger value, WithOption... options) throws BufferException;

    /**
     * Writes a float value with the specified bit length.
     *
     * @param numBits the number of bits to write
     * @param value   the float value
     * @param options additional options for writing
     */
    void writeFloat(int numBits, float value, WithOption... options) throws BufferException;

    /**
     * Writes a double value with the specified bit length.
     *
     * @param numBits the number of bits to write
     * @param value   the double value
     * @param options additional options for writing
     */
    void writeDouble(int numBits, double value, WithOption... options) throws BufferException;

    /**
     * Writes a big decimal value with the specified bit length.
     *
     * @param numBits the number of bits to write
     * @param value   the big decimal value
     * @param options additional options for writing
     */
    void writeBigDecimal(int numBits, BigDecimal value, WithOption... options) throws BufferException;

    /**
     * Writes a string value with the specified bit length.
     *
     * @param numBits the number of bits to write
     * @param value   the string value
     * @param options additional options for writing
     */
    void writeString(int numBits, String value, WithOption... options) throws BufferException;

    /**
     * this method can be used to influence serializing (e.g. intercept whole types and render them in a simplified form)
     *
     * @param message the value to be serialized
     * @throws BufferException if something goes wrong
     */
    default void writeMessage(Message message) throws BufferException {
        if (message == null) {
            return;
        }
        message.serialize(this);
    }

    /**
     * Creates a sub-buffer with the specified bit length.
     *
     * @param numBits the number of bits for the sub-buffer
     * @param options additional options for creating the sub-buffer
     * @return a new WriteBuffer instance representing the sub-buffer
     */
    WriteBuffer createSubBuffer(int numBits, WithOption... options) throws BufferException;

    /**
     * Gets the current position in the buffer in bits.
     *
     * @return the current position in bits
     */
    int getPositionInBits();

    /**
     * Gets the number of remaining bits in the buffer.
     *
     * @return the number of remaining bits
     */
    int getRemainingBits();

    /**
     * Gets the byte array representation of the buffer's content.
     *
     * @return the byte array containing the buffer's data
     */
    byte[] getBytes();

    /**
     * This is only implemented to return true for byte-based byte buffers.
     *
     * @return true, if this is a byte-based buffer.
     */
    default boolean isByteBased() {
        return false;
    }

}
